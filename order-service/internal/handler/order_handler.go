package handler

import (
	"context"
	"store/order-service/internal/proto"
)

type OrderHandler struct {
	proto.UnimplementedOrderServiceServer
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *proto.CreateOrderRequest) (*proto.CreateOrderResponse, error) {
	// Пример: возвращаем заглушку
	return &proto.CreateOrderResponse{
		OrderId: 1,
	}, nil
}
