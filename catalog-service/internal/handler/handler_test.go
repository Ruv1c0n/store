package handler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock" // Используем go.uber.org/mock/gomock
	"store/catalog-service/internal/proto"
	"store/catalog-service/internal/repository/mock" // Импортируем моки
)

func TestAddProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок для CatalogDB
	mockDB := mock.NewMockCatalogDB(ctrl)

	// Создаем экземпляр CatalogHandler с моком
	h := NewCatalogHandler(mockDB)

	// Мокируем вызов AddProduct
	mockDB.EXPECT().
		AddProduct("Test Product", 10, 19.99).
		Return(1, nil)

	// Вызов метода AddProduct
	req := &proto.AddProductRequest{
		ProductName:   "Test Product",
		StockQuantity: 10,
		PricePerUnit:  19.99,
	}
	resp, err := h.AddProduct(context.Background(), req)

	// Проверяем результат
	assert.NoError(t, err)
	assert.Equal(t, int32(1), resp.ProductId)
}

func TestUpdateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок для CatalogDB
	mockDB := mock.NewMockCatalogDB(ctrl)

	// Создаем экземпляр CatalogHandler с моком
	h := NewCatalogHandler(mockDB)

	// Мокируем вызов UpdateProduct
	mockDB.EXPECT().
		UpdateProduct(1, "Updated Product", 20, 29.99).
		Return(nil)

	// Вызов метода UpdateProduct
	req := &proto.UpdateProductRequest{
		ProductId:     1,
		ProductName:   "Updated Product",
		StockQuantity: 20,
		PricePerUnit:  29.99,
	}
	resp, err := h.UpdateProduct(context.Background(), req)

	// Проверяем результат
	assert.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestGetProductByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок для CatalogDB
	mockDB := mock.NewMockCatalogDB(ctrl)

	// Создаем экземпляр CatalogHandler с моком
	h := NewCatalogHandler(mockDB)

	// Мокируем вызов GetProductByID
	mockDB.EXPECT().
		GetProductByID(int32(1)). // Используем int32
		Return("Test Product", 10, 19.99, nil)

	// Вызов метода GetProductByID
	req := &proto.GetProductByIDRequest{
		ProductId: 1,
	}
	resp, err := h.GetProductByID(context.Background(), req)

	// Проверяем результат
	assert.NoError(t, err)
	assert.Equal(t, &proto.Product{
		ProductId:     1,
		ProductName:   "Test Product",
		StockQuantity: 10,
		PricePerUnit:  19.99,
	}, resp.Product)
}

func TestGetAllProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок для CatalogDB
	mockDB := mock.NewMockCatalogDB(ctrl)

	// Создаем экземпляр CatalogHandler с моком
	h := NewCatalogHandler(mockDB)

	// Мокируем вызов GetAllProducts
	expectedProducts := []*proto.Product{
		{
			ProductId:     1,
			ProductName:   "Product 1",
			StockQuantity: 10,
			PricePerUnit:  19.99,
		},
		{
			ProductId:     2,
			ProductName:   "Product 2",
			StockQuantity: 20,
			PricePerUnit:  29.99,
		},
	}
	mockDB.EXPECT().
		GetAllProducts().
		Return(expectedProducts, nil)

	// Вызов метода GetAllProducts
	req := &proto.GetAllProductsRequest{}
	resp, err := h.GetAllProducts(context.Background(), req)

	// Проверяем результат
	assert.NoError(t, err)
	assert.Equal(t, expectedProducts, resp.Products)
}

func TestDeleteProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок для CatalogDB
	mockDB := mock.NewMockCatalogDB(ctrl)

	// Создаем экземпляр CatalogHandler с моком
	h := NewCatalogHandler(mockDB)

	// Мокируем вызов DeleteProduct
	mockDB.EXPECT().
		DeleteProduct(1).
		Return(nil)

	// Вызов метода DeleteProduct
	req := &proto.DeleteProductRequest{
		ProductId: 1,
	}
	resp, err := h.DeleteProduct(context.Background(), req)

	// Проверяем результат
	assert.NoError(t, err)
	assert.True(t, resp.Success)
}
