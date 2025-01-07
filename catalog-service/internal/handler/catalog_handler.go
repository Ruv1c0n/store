package handler

import (
	"context"
	"log"
	db "store/catalog-service/internal/repository"
	proto2 "store/proto"
)

type CatalogHandler struct {
	proto2.UnimplementedProductServiceServer
	db db.CatalogDB // Добавляем поле db
}

func NewCatalogHandler(db db.CatalogDB) *CatalogHandler {
	return &CatalogHandler{db: db}
}

func (h *CatalogHandler) UpdateProduct(ctx context.Context, req *proto2.UpdateProductRequest) (*proto2.UpdateProductResponse, error) {
	log.Printf("Получен запрос UpdateProduct для product_id: %d", req.ProductId)

	// Получаем текущие данные о товаре
	productName, stockQuantity, pricePerUnit, err := h.db.GetProductByID(req.ProductId)
	if err != nil {
		log.Printf("Ошибка при получении товара: %v", err)
		return nil, err
	}

	// Обновляем только те поля, которые переданы в запросе
	if req.ProductName != "" {
		productName = req.ProductName
	}
	if req.StockQuantity != 0 {
		stockQuantity = int(req.StockQuantity)
	}
	if req.PricePerUnit != 0 {
		pricePerUnit = req.PricePerUnit
	}

	// Обновляем товар в базе данных
	err = h.db.UpdateProduct(int(req.ProductId), productName, stockQuantity, pricePerUnit)
	if err != nil {
		log.Printf("Ошибка при обновлении товара: %v", err)
		return nil, err
	}

	// Возвращаем успешный ответ
	return &proto2.UpdateProductResponse{
		Success: true,
	}, nil
}

func (h *CatalogHandler) AddProduct(ctx context.Context, req *proto2.AddProductRequest) (*proto2.AddProductResponse, error) {
	log.Printf("Получен запрос AddProduct: %v", req)

	// Добавляем продукт в базу данных
	productID, err := h.db.AddProduct(req.ProductName, int(req.StockQuantity), req.PricePerUnit)
	if err != nil {
		log.Printf("Ошибка при добавлении продукта: %v", err)
		return nil, err
	}

	// Возвращаем ответ
	return &proto2.AddProductResponse{
		ProductId: int32(productID),
	}, nil
}

func (h *CatalogHandler) GetProductByID(ctx context.Context, req *proto2.GetProductByIDRequest) (*proto2.GetProductByIDResponse, error) {
	log.Printf("Получен запрос GetProductByID для product_id: %d", req.ProductId)

	// Используем реальную базу данных
	productName, stockQuantity, pricePerUnit, err := h.db.GetProductByID(req.ProductId)
	if err != nil {
		log.Printf("Ошибка при получении продукта: %v", err)
		return nil, err
	}

	// Возвращаем ответ
	return &proto2.GetProductByIDResponse{
		Product: &proto2.Product{
			ProductId:     req.ProductId,
			ProductName:   productName,
			StockQuantity: int32(stockQuantity),
			PricePerUnit:  pricePerUnit,
		},
	}, nil
}

func (h *CatalogHandler) GetAllProducts(ctx context.Context, req *proto2.GetAllProductsRequest) (*proto2.GetAllProductsResponse, error) {
	log.Println("Получен запрос GetAllProducts")

	// Получаем все продукты из базы данных
	products, err := h.db.GetAllProducts()
	if err != nil {
		log.Printf("Ошибка при получении продуктов: %v", err)
		return nil, err
	}

	// Возвращаем ответ
	return &proto2.GetAllProductsResponse{
		Products: products,
	}, nil
}

func (h *CatalogHandler) DeleteProduct(ctx context.Context, req *proto2.DeleteProductRequest) (*proto2.DeleteProductResponse, error) {
	log.Printf("Получен запрос DeleteProduct для product_id: %d", req.ProductId)

	// Удаляем продукт из базы данных
	err := h.db.DeleteProduct(int(req.ProductId))
	if err != nil {
		log.Printf("Ошибка при удалении продукта: %v", err)
		return nil, err
	}

	// Возвращаем ответ
	return &proto2.DeleteProductResponse{
		Success: true,
	}, nil
}
