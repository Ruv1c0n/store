package handler

import (
	"context"
	"database/sql"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	mock "store/order-service/internal/repository/mock"
	"store/proto"
)

func TestCreateOrder_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockOrderDB(ctrl)
	handler := NewOrderHandler(mockDB)

	customerID := int32(1)
	productID := int32(2)
	quantity := int32(2)
	pricePerUnit := 19.99
	orderID := int32(2)
	stockQuantity := 10

	mockDB.EXPECT().
		GetNextOrderID(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, id *int32) error {
			*id = orderID
			return nil
		})

	mockDB.EXPECT().
		GetProductByID(productID).
		Return("ProductName", stockQuantity, pricePerUnit, nil)

	mockDB.EXPECT().
		CreateOrder(gomock.Any(), orderID, productID, customerID, quantity, pricePerUnit).
		Return(nil)

	req := &proto.CreateOrderRequest{
		CustomerId: customerID,
		Items: []*proto.OrderItem{
			{
				ProductId: productID,
				Quantity:  quantity,
			},
		},
	}
	resp, err := handler.CreateOrder(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, orderID, resp.OrderId)
}

func TestCreateOrder_InsufficientStock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock for OrderDB
	mockDB := mock.NewMockOrderDB(ctrl)

	// Create the OrderHandler with the mock
	handler := NewOrderHandler(mockDB)

	// Define test data
	customerID := int32(1)
	productID := int32(2)
	quantity := int32(20) // Requested quantity exceeds stock
	stockQuantity := 10

	// Mock GetNextOrderID to return a new order ID
	mockDB.EXPECT().
		GetNextOrderID(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, id *int32) error {
			*id = int32(1)
			return nil
		})

	// Mock GetProductByID to return product details
	mockDB.EXPECT().
		GetProductByID(productID).
		Return("ProductName", stockQuantity, 19.99, nil)

	// Call the CreateOrder method
	req := &proto.CreateOrderRequest{
		CustomerId: customerID,
		Items: []*proto.OrderItem{
			{
				ProductId: productID,
				Quantity:  quantity,
			},
		},
	}
	resp, err := handler.CreateOrder(context.Background(), req)

	// Assert the results
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Not enough stock for the product")
}

func TestGetOrderByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock for OrderDB
	mockDB := mock.NewMockOrderDB(ctrl)

	// Create the OrderHandler with the mock
	handler := NewOrderHandler(mockDB)

	// Define test data
	orderID := int32(2)
	order := &proto.Order{
		OrderId: orderID,
		Items: []*proto.OrderItem{
			{
				ProductId: 2,
				Quantity:  3,
			},
		},
		OrderDate:  "2025-01-07 21:13:18.41648",
		Status:     "в обработке",
		CustomerId: 2,
	}

	// Mock GetOrderByID to return the order
	mockDB.EXPECT().
		GetOrderByID(orderID).
		Return(order, nil)

	// Call the GetOrderByID method
	req := &proto.GetOrderByIDRequest{
		OrderId: orderID,
	}
	resp, err := handler.GetOrderByID(context.Background(), req)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, order, resp.Order)
}

func TestGetOrderByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок для OrderDB
	mockDB := mock.NewMockOrderDB(ctrl)

	// Создаем OrderHandler с моком
	handler := NewOrderHandler(mockDB)

	// Определяем тестовые данные
	orderID := int32(1)

	// Мокируем вызов GetOrderByID, чтобы он возвращал ошибку "order not found"
	mockDB.EXPECT().
		GetOrderByID(orderID).
		Return(nil, sql.ErrNoRows) // Используем sql.ErrNoRows для имитации отсутствия заказа

	// Вызываем метод GetOrderByID
	req := &proto.GetOrderByIDRequest{
		OrderId: orderID,
	}
	resp, err := handler.GetOrderByID(context.Background(), req)

	// Проверяем, что возвращена ошибка с кодом NotFound и сообщением "Заказ не найден"
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Заказ не найден")

	// Проверяем, что ошибка имеет код NotFound
	status, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, status.Code())
}

func TestUpdateOrder_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock.NewMockOrderDB(ctrl)
	handler := NewOrderHandler(mockDB)

	orderID := int32(1)
	status := "new_status"

	mockDB.EXPECT().
		UpdateOrder(orderID, status).
		Return(nil)

	req := &proto.UpdateOrderRequest{
		OrderId: orderID,
		Status:  status,
	}
	resp, err := handler.UpdateOrder(context.Background(), req)

	assert.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestDeleteOrder_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock for OrderDB
	mockDB := mock.NewMockOrderDB(ctrl)

	// Create the OrderHandler with the mock
	handler := NewOrderHandler(mockDB)

	// Define test data
	orderID := int32(2)

	// Mock DeleteOrder to simulate successful deletion
	mockDB.EXPECT().
		DeleteOrder(orderID).
		Return(nil)

	// Call the DeleteOrder method
	req := &proto.DeleteOrderRequest{
		OrderId: orderID,
	}
	resp, err := handler.DeleteOrder(context.Background(), req)

	// Assert the results
	assert.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestDeleteOrder_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock for OrderDB
	mockDB := mock.NewMockOrderDB(ctrl)

	// Create the OrderHandler with the mock
	handler := NewOrderHandler(mockDB)

	// Define test data
	orderID := int32(2)

	// Mock DeleteOrder to return an error
	mockDB.EXPECT().
		DeleteOrder(orderID).
		Return(errors.New("failed to delete order"))

	// Call the DeleteOrder method
	req := &proto.DeleteOrderRequest{
		OrderId: orderID,
	}
	resp, err := handler.DeleteOrder(context.Background(), req)

	// Assert the results
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestCreateOrder_DBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock for OrderDB
	mockDB := mock.NewMockOrderDB(ctrl)

	// Create the OrderHandler with the mock
	handler := NewOrderHandler(mockDB)

	// Define test data
	customerID := int32(1)
	productID := int32(2)
	quantity := int32(2)

	// Mock GetNextOrderID to simulate a database error
	mockDB.EXPECT().
		GetNextOrderID(gomock.Any(), gomock.Any()).
		Return(errors.New("database error"))

	// Call the CreateOrder method
	req := &proto.CreateOrderRequest{
		CustomerId: customerID,
		Items: []*proto.OrderItem{
			{
				ProductId: productID,
				Quantity:  quantity,
			},
		},
	}
	resp, err := handler.CreateOrder(context.Background(), req)

	// Assert the results
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "database error")
}

func TestCreateOrder_StockValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock for OrderDB
	mockDB := mock.NewMockOrderDB(ctrl)

	// Create the OrderHandler with the mock
	handler := NewOrderHandler(mockDB)

	// Define test data
	customerID := int32(1)
	productID := int32(2)
	quantity := int32(5)
	stockQuantity := 3 // Requested quantity exceeds stock

	// Mock GetNextOrderID to return a new order ID
	mockDB.EXPECT().GetNextOrderID(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, id *int32) error {
		*id = int32(1)
		return nil
	})

	// Mock GetProductByID to return product details
	mockDB.EXPECT().GetProductByID(productID).Return("ProductName", stockQuantity, 19.99, nil)

	// Call the CreateOrder method
	req := &proto.CreateOrderRequest{
		CustomerId: customerID,
		Items: []*proto.OrderItem{
			{
				ProductId: productID,
				Quantity:  quantity,
			},
		},
	}
	resp, err := handler.CreateOrder(context.Background(), req)

	// Assert the results
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Not enough stock for the product")
}

func TestCreateOrder_ProductNotFoundError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создание мока для OrderDB
	mockDB := mock.NewMockOrderDB(ctrl)

	// Создание обработчика с моками
	handler := NewOrderHandler(mockDB)

	// Определение тестовых данных
	customerID := int32(1)
	productID := int32(2)
	quantity := int32(2)

	// Мокаем GetNextOrderID для возврата нового ID заказа
	mockDB.EXPECT().GetNextOrderID(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, id *int32) error {
		*id = 1
		return nil
	})

	// Мокаем ошибку при получении товара
	mockDB.EXPECT().GetProductByID(productID).Return("", 0, 0.0, errors.New("Product not found"))

	// Запрос для создания заказа
	req := &proto.CreateOrderRequest{
		CustomerId: customerID,
		Items: []*proto.OrderItem{
			{
				ProductId: productID,
				Quantity:  quantity,
			},
		},
	}

	// Выполнение запроса
	resp, err := handler.CreateOrder(context.Background(), req)

	// Проверки
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Product not found")
}

func TestCreateOrder_CreateOrderError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создание мока для OrderDB
	mockDB := mock.NewMockOrderDB(ctrl)

	// Создание обработчика с моками
	handler := NewOrderHandler(mockDB)

	// Определение тестовых данных
	customerID := int32(1)
	productID := int32(2)
	quantity := int32(2)
	orderID := int32(2)
	stockQuantity := 10

	// Мокаем GetNextOrderID для возврата нового ID заказа
	mockDB.EXPECT().GetNextOrderID(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, id *int32) error {
		*id = orderID
		return nil
	})

	// Мокаем GetProductByID для возврата данных о товаре
	mockDB.EXPECT().GetProductByID(productID).Return("ProductName", stockQuantity, 19.99, nil)

	// Мокаем ошибку при создании заказа
	mockDB.EXPECT().CreateOrder(gomock.Any(), orderID, productID, customerID, quantity, 19.99).Return(errors.New("Database error"))

	// Запрос для создания заказа
	req := &proto.CreateOrderRequest{
		CustomerId: customerID,
		Items: []*proto.OrderItem{
			{
				ProductId: productID,
				Quantity:  quantity,
			},
		},
	}

	// Выполнение запроса
	resp, err := handler.CreateOrder(context.Background(), req)

	// Проверки
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Database error")
}
