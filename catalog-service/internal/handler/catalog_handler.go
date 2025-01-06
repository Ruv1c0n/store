package handler

import (
	"context"
	"log"
	"store/catalog-service/internal/proto"
	db "store/catalog-service/internal/repository"
)

type CatalogHandler struct {
	proto.UnimplementedProductServiceServer
	db db.CatalogDB // Добавляем поле db
}

func NewCatalogHandler(db db.CatalogDB) *CatalogHandler {
	return &CatalogHandler{db: db}
}

func (h *CatalogHandler) UpdateProduct(ctx context.Context, req *proto.UpdateProductRequest) (*proto.UpdateProductResponse, error) {
	log.Printf("Получен запрос UpdateProduct для product_id: %d", req.ProductId)

	// Обновляем продукт в базе данных
	err := h.db.UpdateProduct(int(req.ProductId), req.ProductName, int(req.StockQuantity), req.PricePerUnit)
	if err != nil {
		log.Printf("Ошибка при обновлении продукта: %v", err)
		return nil, err
	}

	// Возвращаем ответ
	return &proto.UpdateProductResponse{
		Success: true,
	}, nil
}

func (h *CatalogHandler) AddProduct(ctx context.Context, req *proto.AddProductRequest) (*proto.AddProductResponse, error) {
	log.Printf("Получен запрос AddProduct: %v", req)

	// Добавляем продукт в базу данных
	productID, err := h.db.AddProduct(req.ProductName, int(req.StockQuantity), req.PricePerUnit)
	if err != nil {
		log.Printf("Ошибка при добавлении продукта: %v", err)
		return nil, err
	}

	// Возвращаем ответ
	return &proto.AddProductResponse{
		ProductId: int32(productID),
	}, nil
}

func (h *CatalogHandler) GetProductByID(ctx context.Context, req *proto.GetProductByIDRequest) (*proto.GetProductByIDResponse, error) {
	log.Printf("Получен запрос GetProductByID для product_id: %d", req.ProductId)

	// Используем реальную базу данных
	productName, stockQuantity, pricePerUnit, err := h.db.GetProductByID(req.ProductId)
	if err != nil {
		log.Printf("Ошибка при получении продукта: %v", err)
		return nil, err
	}

	// Возвращаем ответ
	return &proto.GetProductByIDResponse{
		Product: &proto.Product{
			ProductId:     req.ProductId,
			ProductName:   productName,
			StockQuantity: int32(stockQuantity),
			PricePerUnit:  pricePerUnit,
		},
	}, nil
}

func (h *CatalogHandler) GetAllProducts(ctx context.Context, req *proto.GetAllProductsRequest) (*proto.GetAllProductsResponse, error) {
	log.Println("Получен запрос GetAllProducts")

	// Получаем все продукты из базы данных
	products, err := h.db.GetAllProducts()
	if err != nil {
		log.Printf("Ошибка при получении продуктов: %v", err)
		return nil, err
	}

	// Возвращаем ответ
	return &proto.GetAllProductsResponse{
		Products: products,
	}, nil
}

func (h *CatalogHandler) DeleteProduct(ctx context.Context, req *proto.DeleteProductRequest) (*proto.DeleteProductResponse, error) {
	log.Printf("Получен запрос DeleteProduct для product_id: %d", req.ProductId)

	// Удаляем продукт из базы данных
	err := h.db.DeleteProduct(int(req.ProductId))
	if err != nil {
		log.Printf("Ошибка при удалении продукта: %v", err)
		return nil, err
	}

	// Возвращаем ответ
	return &proto.DeleteProductResponse{
		Success: true,
	}, nil
}
