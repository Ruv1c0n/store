package client

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"store/proto"
)

// CatalogClient интерфейс для взаимодействия с catalog-service
type CatalogClient interface {
	UpdateProductStock(productID int32, newStockQuantity int32) error
	Close()
}

// CatalogClientImpl реализует интерфейс CatalogClient
type CatalogClientImpl struct {
	conn   *grpc.ClientConn
	client proto.ProductServiceClient
}

// NewCatalogClient создает новый экземпляр CatalogClient
func NewCatalogClient(address string) (CatalogClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure()) // Устанавливаем соединение
	if err != nil {
		return nil, err
	}
	client := proto.NewProductServiceClient(conn) // Создаем клиент
	return &CatalogClientImpl{conn: conn, client: client}, nil
}

// Close закрывает соединение с catalog-service
func (c *CatalogClientImpl) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

// UpdateProductStock обновляет количество товара в каталоге через gRPC
func (c *CatalogClientImpl) UpdateProductStock(productID int32, newStockQuantity int32) error {
	req := &proto.UpdateProductRequest{
		ProductId:     productID,
		StockQuantity: newStockQuantity,
	}
	_, err := c.client.UpdateProduct(context.Background(), req) // Вызываем метод catalog-service
	if err != nil {
		log.Printf("Failed to update product stock: %v", err)
		return err
	}
	return nil
}
