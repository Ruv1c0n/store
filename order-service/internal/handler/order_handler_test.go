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

	// Create a mock for OrderDB
	mockDB := mock.NewMockOrderDB(ctrl)

	// Create the OrderHandler with the mock
	handler := NewOrderHandler(mockDB)

	// Define test data
	customerID := int32(1)
	productID := int32(2)
	quantity := int32(2)
	pricePerUnit := 19.99
	orderID := int32(2)
	stockQuantity := 10

	// Mock GetNextOrderID to return a new order ID
	mockDB.EXPECT().
		GetNextOrderID(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, id *int32) error {
			*id = orderID
			return nil
		})

	// Mock GetProductByID to return product details
	mockDB.EXPECT().
		GetProductByID(productID).
		Return("ProductName", stockQuantity, pricePerUnit, nil)

	// Mock CreateOrder to simulate creating an order
	mockDB.EXPECT().
		CreateOrder(gomock.Any(), orderID, productID, customerID, quantity, pricePerUnit).
		Return(nil)

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
