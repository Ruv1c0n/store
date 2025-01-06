package db

import (
    "context"
    "fmt"
    "github.com/jackc/pgx/v4"
    "log"
)

var conn *pgx.Conn

func Connect() {
    var err error
    conn, err = pgx.Connect(context.Background(), "postgres://username:password@localhost:5432/yourdbname")
    if err != nil {
        log.Fatalf("Unable to connect to database: %v\n", err)
    }
    fmt.Println("Connected to PostgreSQL!")
}

func Disconnect() {
    if err := conn.Close(context.Background()); err != nil {
        log.Fatalf("Unable to close connection: %v\n", err)
    }
    fmt.Println("Disconnected from PostgreSQL!")
}

// Добавление нового заказа
func AddOrder(orderID int, productID int, customerID int, quantity int, pricePerUnit float64) error {
    _, err := conn.Exec(context.Background(),
        "INSERT INTO Orders (OrderID, ProductID, CustomerID, Quantity, PricePerUnit) VALUES ($1, $2, $3, $4, $5)",
        orderID, productID, customerID, quantity, pricePerUnit)
    return err
}

// Получение всех заказов для конкретного продукта
func GetOrdersByProductID(productID int) ([]Order, error) {
    rows, err := conn.Query(context.Background(),
        "SELECT OrderID, OrderDate, Status, CustomerID, Quantity, PricePerUnit, TotalPrice FROM Orders WHERE ProductID=$1", productID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var orders []Order
    for rows.Next() {
        var order Order
        err := rows.Scan(&order.OrderID, &order.OrderDate, &order.Status, &order.CustomerID, &order.Quantity, &order.PricePerUnit, &order.TotalPrice)
        if err != nil {
            return nil, err
        }
        orders = append(orders, order)
    }
    return orders, nil
}

// Обновление информации о заказе
func UpdateOrder(orderID int, productID int, customerID int, quantity int, pricePerUnit float64, status string) error {
    _, err := conn.Exec(context.Background(),
        "UPDATE Orders SET ProductID=$1, CustomerID=$2, Quantity=$3, PricePerUnit=$4, Status=$5 WHERE OrderID=$6",
        productID, customerID, quantity, pricePerUnit, status, orderID)
    return err
}

// Удаление заказа
func DeleteOrder(orderID int, productID int) error {
    _, err := conn.Exec(context.Background(),
        "DELETE FROM Orders WHERE OrderID=$1 AND ProductID=$2", orderID, productID)
    return err
}

// Определите структуру Order для хранения данных заказа
type Order struct {
    OrderID      int
    OrderDate    string
    Status       string
    CustomerID   int
    Quantity     int
    PricePerUnit float64
    TotalPrice   float64
}