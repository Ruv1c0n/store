package db

import (
	"context"
	"fmt"
	"log"
	"store/order-service/internal/client"
	"store/proto"
	"time"

	"github.com/jackc/pgx/v4"
)

// OrderDB интерфейс для работы с заказами
type OrderDB interface {
	GetNextOrderID(ctx context.Context, orderID *int32) error
	CreateOrder(ctx context.Context, orderID int32, productID int32, customerID int32, quantity int32, pricePerUnit float64) error
	GetOrderByID(orderID int32) (*proto.Order, error)
	GetAllOrders() ([]*proto.Order, error)
	UpdateOrder(orderID int32, status string) error
	DeleteOrder(orderID int32) error
	// GetProductByID(productID int32) (string, int, float64, error)
}

//go:generate mockgen -source=db.go -destination=mock/mock.go

// orderDB реализует интерфейс OrderDB
type orderDB struct {
	conn          *pgx.Conn
	catalogClient client.CatalogClient // Используем интерфейс
}

// NewOrderDB создает новый экземпляр orderDB
func NewOrderDB(conn *pgx.Conn, catalogClient client.CatalogClient) OrderDB {
	return &orderDB{
		conn:          conn,
		catalogClient: catalogClient,
	}
}

// CreateOrder добавляет новый заказ в базу данных
func (db *orderDB) GetNextOrderID(ctx context.Context, orderID *int32) error {
	return db.conn.QueryRow(ctx, "SELECT nextval('orders_orderid_seq')").Scan(orderID)
}

func (db *orderDB) CreateOrder(ctx context.Context, orderID int32, productID int32, customerID int32, quantity int32, pricePerUnit float64) error {
	_, err := db.conn.Exec(ctx, `
        INSERT INTO Orders (OrderID, ProductID, CustomerID, Quantity, PricePerUnit)
        VALUES ($1, $2, $3, $4, $5)
    `, orderID, productID, customerID, quantity, pricePerUnit)
	return err
}

// GetOrderByID возвращает заказ по его ID
func (db *orderDB) GetOrderByID(orderID int32) (*proto.Order, error) {
	// Основная информация о заказе
	var order proto.Order
	order.OrderId = orderID

	// Получаем список продуктов в заказе
	rows, err := db.conn.Query(context.Background(), `
        SELECT productid, quantity 
        FROM orders 
        WHERE orderid = $1`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Чтение данных о продуктах
	for rows.Next() {
		var item proto.OrderItem
		err := rows.Scan(&item.ProductId, &item.Quantity)
		if err != nil {
			return nil, err
		}

		// Добавляем продукт в заказ
		order.Items = append(order.Items, &item)
	}

	// Получаем общую информацию о заказе (дата, статус, клиент)
	var orderDate time.Time
	err = db.conn.QueryRow(context.Background(), `
        SELECT orderdate, status, customerid 
        FROM orders 
        WHERE orderid = $1 
        LIMIT 1`, orderID).Scan(&orderDate, &order.Status, &order.CustomerId)
	if err != nil {
		return nil, err
	}

	// Преобразуем время в строку
	order.OrderDate = orderDate.Format(time.RFC3339)

	return &order, nil
}

func (db *orderDB) GetAllOrders() ([]*proto.Order, error) {
	rows, err := db.conn.Query(context.Background(), `
        SELECT orderid, productid, quantity, orderdate, status, customerid
        FROM orders`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*proto.Order
	orderMap := make(map[int32]*proto.Order) // Для группировки товаров по заказам

	for rows.Next() {
		var orderID int32
		var productID int32
		var quantity int32
		var orderDate time.Time
		var status string
		var customerID int32

		err := rows.Scan(&orderID, &productID, &quantity, &orderDate, &status, &customerID)
		if err != nil {
			return nil, err
		}

		// Преобразование времени в строку
		orderDateStr := orderDate.Format(time.RFC3339)

		// Если заказ с таким ID уже есть в мапе, добавляем товар в его список
		if order, exists := orderMap[orderID]; exists {
			order.Items = append(order.Items, &proto.OrderItem{
				ProductId: productID,
				Quantity:  quantity,
			})
		} else {
			// Создаем новый заказ и добавляем его в мапу
			order := &proto.Order{
				OrderId:    orderID,
				OrderDate:  orderDateStr,
				Status:     status,
				CustomerId: customerID,
				Items: []*proto.OrderItem{
					{
						ProductId: productID,
						Quantity:  quantity,
					},
				},
			}
			orderMap[orderID] = order
			orders = append(orders, order)
		}
	}

	return orders, nil
}

func (db *orderDB) UpdateOrder(orderID int32, status string) error {
	// Начинаем транзакцию
	tx, err := db.conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(context.Background())

	// Обновляем статус заказа
	_, err = tx.Exec(context.Background(), `
        UPDATE Orders 
        SET status = $1 
        WHERE orderid = $2`, status, orderID)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	// Завершаем транзакцию
	if err := tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// DeleteOrder удаляет заказ из базы данных и восстанавливает количество товаров в каталоге
func (db *orderDB) DeleteOrder(orderID int32) error {
	// Начинаем транзакцию
	tx, err := db.conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(context.Background())

	// Получаем информацию о заказе
	order, err := db.GetOrderByID(orderID)
	if err != nil {
		return fmt.Errorf("failed to get order details: %w", err)
	}

	// Создаем gRPC-клиент для catalog-service
	catalogClient, err := client.NewCatalogClient("localhost:50051") // Укажите адрес catalog-service
	if err != nil {
		log.Printf("Ошибка при создании gRPC-клиента для catalog-service: %v", err)
		return err
	}
	defer catalogClient.Close()

	// Восстанавливаем количество товаров в каталоге
	for _, item := range order.Items {
		// Получаем текущее количество товара на складе из каталога
		_, stockQuantity, _, err := catalogClient.GetProductByID(item.ProductId)
		if err != nil {
			return fmt.Errorf("failed to get product stock quantity: %w", err)
		}

		// Восстанавливаем количество товара на складе
		newStockQuantity := stockQuantity + int(item.Quantity)
		err = db.catalogClient.UpdateProductStock(item.ProductId, int32(newStockQuantity))
		if err != nil {
			return fmt.Errorf("failed to update catalog stock via gRPC: %w", err)
		}
	}

	// Удаляем заказ из базы данных
	_, err = tx.Exec(context.Background(), `DELETE FROM orders WHERE orderid = $1`, orderID)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	// Завершаем транзакцию
	if err := tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// // GetProductByID возвращает информацию о товаре по его ID
// func (db *orderDB) GetProductByID(productID int32) (string, int, float64, error) {
// 	var productName string
// 	var stockQuantity int
// 	var pricePerUnit float64
// 	err := db.conn.QueryRow(context.Background(), `
//         SELECT ProductName, StockQuantity, PricePerUnit FROM Catalog 
//         WHERE ProductID = $1`, productID).Scan(&productName, &stockQuantity, &pricePerUnit)
// 	if err != nil {
// 		return "", 0, 0, err
// 	}
// 	return productName, stockQuantity, pricePerUnit, nil
// }
