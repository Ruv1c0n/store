package main

import (
	"log"
	"net"

	"store/order-service/internal/repository" 
	"store/order-service/internal/handler"
	"store/order-service/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterOrderServiceServer(grpcServer, handler.NewOrderHandler())

	// Включение рефлексии
	reflection.Register(grpcServer)
	db.Connect() // Подключение к базе данных
    defer db.Disconnect() // Отключение при завершении работы

	log.Println("OrderService is running on port 50052...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}