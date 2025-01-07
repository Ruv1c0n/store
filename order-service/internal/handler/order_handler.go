package handler

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"store/order-service/internal/client"
	db "store/order-service/internal/repository"
	proto2 "store/proto"
)

type OrderHandler struct {
	proto2.UnimplementedOrderServiceServer
	db db.OrderDB // Поле для работы с базой данных
}

func NewOrderHandler(db db.OrderDB) *OrderHandler {
	return &OrderHandler{db: db}
}

// CreateOrder обрабатывает создание нового заказа
func (h *OrderHandler) CreateOrder(ctx context.Context, req *proto2.CreateOrderRequest) (*proto2.CreateOrderResponse, error) {
	log.Printf("Получен запрос CreateOrder для customer_id: %d", req.CustomerId)

	// Генерируем новый OrderID
	var orderID int32
	err := h.db.GetNextOrderID(ctx, &orderID)
	if err != nil {
		log.Printf("Ошибка при генерации OrderID: %v", err)
		return nil, err
	}

	// Создаем gRPC-клиент для catalog-service
	catalogClient, err := client.NewCatalogClient("localhost:50051") // Укажите адрес catalog-service
	if err != nil {
		log.Printf("Ошибка при создании gRPC-клиента для catalog-service: %v", err)
		return nil, err
	}
	defer catalogClient.Close()

	// Обрабатываем каждый товар в заказе
	for _, item := range req.Items {
		// Получаем информацию о товаре, включая цену
		_, stockQuantity, pricePerUnit, err := h.db.GetProductByID(item.ProductId)
		if err != nil {
			log.Printf("Ошибка при получении товара: %v", err)
			return nil, err
		}

		// Проверяем наличие товара в достаточном количестве
		if stockQuantity < int(item.Quantity) {
			log.Printf("Недостаточно товара в наличии для product_id: %d", item.ProductId)
			return nil, fmt.Errorf("Not enough stock for the product")
		}

		// Создаем запись в таблице Orders
		err = h.db.CreateOrder(ctx, orderID, item.ProductId, req.CustomerId, item.Quantity, pricePerUnit)
		if err != nil {
			log.Printf("Ошибка при создании заказа: %v", err)
			return nil, err
		}

		// Обновляем количество товара в каталоге через gRPC
		newStockQuantity := stockQuantity - int(item.Quantity)
		err = catalogClient.UpdateProductStock(item.ProductId, int32(newStockQuantity))
		if err != nil {
			log.Printf("Ошибка при обновлении количества товара: %v", err)
			return nil, err
		}

		log.Printf("Добавлен товар в заказ: OrderID=%d, ProductID=%d, Quantity=%d", orderID, item.ProductId, item.Quantity)
	}

	log.Printf("Создан заказ с OrderID: %d", orderID)

	// Возвращаем ответ
	return &proto2.CreateOrderResponse{
		OrderId: orderID,
	}, nil
}

// GetOrderByID обрабатывает запрос на получение заказа по ID
func (h *OrderHandler) GetOrderByID(ctx context.Context, req *proto2.GetOrderByIDRequest) (*proto2.GetOrderByIDResponse, error) {
	log.Printf("Получен запрос GetOrderByID для order_id: %d", req.OrderId)

	// Получаем заказ из базы данных
	order, err := h.db.GetOrderByID(req.OrderId)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Если заказ не найден, возвращаем ошибку "Not Found"
			log.Printf("Заказ с order_id=%d не найден", req.OrderId)
			return nil, status.Errorf(codes.NotFound, "Заказ не найден")
		}
		// В случае других ошибок возвращаем Internal Server Error
		log.Printf("Ошибка при получении заказа: %v", err)
		return nil, status.Errorf(codes.Internal, "Внутренняя ошибка сервера")
	}

	// Возвращаем заказ
	return &proto2.GetOrderByIDResponse{
		Order: order,
	}, nil
}

// GetAllOrders обрабатывает запрос на получение всех заказов
func (h *OrderHandler) GetAllOrders(ctx context.Context, req *proto2.GetAllOrdersRequest) (*proto2.GetAllOrdersResponse, error) {
	log.Println("Получен запрос GetAllOrders")

	// Получаем все заказы из базы данных
	orders, err := h.db.GetAllOrders()
	if err != nil {
		log.Printf("Ошибка при получении заказов: %v", err)
		return nil, err
	}

	// Возвращаем ответ
	return &proto2.GetAllOrdersResponse{
		Orders: orders,
	}, nil
}

// UpdateOrder обрабатывает запрос на обновление заказа
func (h *OrderHandler) UpdateOrder(ctx context.Context, req *proto2.UpdateOrderRequest) (*proto2.UpdateOrderResponse, error) {
	log.Printf("Получен запрос UpdateOrder для order_id: %d", req.OrderId)

	// Обновляем информацию о заказе
	err := h.db.UpdateOrder(req.OrderId, req.Quantity, req.Status)
	if err != nil {
		log.Printf("Ошибка при обновлении заказа: %v", err)
		return nil, err
	}

	// Возвращаем успешный ответ
	return &proto2.UpdateOrderResponse{
		Success: true,
	}, nil
}

// DeleteOrder обрабатывает запрос на удаление заказа
func (h *OrderHandler) DeleteOrder(ctx context.Context, req *proto2.DeleteOrderRequest) (*proto2.DeleteOrderResponse, error) {
	log.Printf("Получен запрос DeleteOrder для order_id: %d", req.OrderId)

	// Удаляем заказ из базы данных
	err := h.db.DeleteOrder(req.OrderId)
	if err != nil {
		log.Printf("Ошибка при удалении заказа: %v", err)
		return nil, err
	}

	// Возвращаем успешный ответ
	return &proto2.DeleteOrderResponse{
		Success: true,
	}, nil
}
