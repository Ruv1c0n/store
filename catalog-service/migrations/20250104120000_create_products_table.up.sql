CREATE TABLE Catalog (
    ProductID SERIAL PRIMARY KEY,
    ProductName VARCHAR(255) NOT NULL,
    StockQuantity INT NOT NULL,
    PricePerUnit NUMERIC(10, 2) NOT NULL
);