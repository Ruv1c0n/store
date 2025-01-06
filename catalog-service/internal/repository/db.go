package db

import (
	"context"
	"github.com/jackc/pgx/v4"
	"store/catalog-service/internal/proto"
)

//go:generate mockgen -source=db.go -destination=mock/mock.go

type CatalogDB interface {
	AddProduct(productName string, stockQuantity int, pricePerUnit float64) (int, error)
	GetProductByID(productID int32) (string, int, float64, error) // Используем int32
	GetAllProducts() ([]*proto.Product, error)
	UpdateProduct(productID int, productName string, stockQuantity int, pricePerUnit float64) error
	DeleteProduct(productID int) error
}

// catalogDB реализует интерфейс CatalogDB
type catalogDB struct {
	conn *pgx.Conn
}

// NewCatalogDB создает новый экземпляр catalogDB
func NewCatalogDB(conn *pgx.Conn) CatalogDB {
	return &catalogDB{conn: conn}
}

func (db *catalogDB) AddProduct(productName string, stockQuantity int, pricePerUnit float64) (int, error) {
	var productID int
	err := db.conn.QueryRow(context.Background(),
		"INSERT INTO Catalog (ProductName, StockQuantity, PricePerUnit) VALUES ($1, $2, $3) RETURNING ProductID",
		productName, stockQuantity, pricePerUnit).Scan(&productID)
	if err != nil {
		return 0, err
	}
	return productID, nil
}

func (db *catalogDB) GetAllProducts() ([]*proto.Product, error) {
	rows, err := db.conn.Query(context.Background(), "SELECT ProductID, ProductName, StockQuantity, PricePerUnit FROM Catalog")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*proto.Product
	for rows.Next() {
		var product proto.Product
		err := rows.Scan(&product.ProductId, &product.ProductName, &product.StockQuantity, &product.PricePerUnit)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, nil
}

func (db *catalogDB) GetProductByID(productID int32) (string, int, float64, error) {
	var productName string
	var stockQuantity int
	var pricePerUnit float64
	err := db.conn.QueryRow(context.Background(),
		"SELECT ProductName, StockQuantity, PricePerUnit FROM Catalog WHERE ProductID=$1", productID).
		Scan(&productName, &stockQuantity, &pricePerUnit)
	if err != nil {
		return "", 0, 0, err
	}
	return productName, stockQuantity, pricePerUnit, nil
}

func (db *catalogDB) UpdateProduct(productID int, productName string, stockQuantity int, pricePerUnit float64) error {
	_, err := db.conn.Exec(context.Background(),
		"UPDATE Catalog SET ProductName=$1, StockQuantity=$2, PricePerUnit=$3 WHERE ProductID=$4",
		productName, stockQuantity, pricePerUnit, productID)
	return err
}

func (db *catalogDB) DeleteProduct(productID int) error {
	_, err := db.conn.Exec(context.Background(), "DELETE FROM Catalog WHERE ProductID=$1", productID)
	return err
}
