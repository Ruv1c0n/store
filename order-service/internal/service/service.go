package service

import (
	"log"

	"google.golang.org/grpc"
)

// Предположим, что у вас есть структура Product
type Product struct {
	Id    int64
	Name  string
	Price float64
}

// Функция для проверки доступности продукта
func CheckProductAvailability(productID int64) {
	// Устанавливаем соединение с другим сервисом, если это необходимо
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to service: %v", err)
	}
	defer conn.Close()

	// Здесь вы можете использовать свой клиент для получения информации о продукте
	// Например, если у вас есть другой gRPC клиент, замените следующую строку
	// client := your.NewYourServiceClient(conn)

	// Вместо вызова GetProduct, вы можете использовать другую логику
	// Например, просто создадим продукт вручную для демонстрации
	product := Product{
		Id:    productID,
		Name:  "Sample Product",
		Price: 19.99,
	}

	// Обрабатываем ответ
	log.Printf("Product received: ID=%d, Name=%s, Price=%.2f", product.Id, product.Name, product.Price)
}