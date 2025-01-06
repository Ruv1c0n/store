CREATE TABLE Orders (
    OrderID INT NOT NULL,
    ProductID INT NOT NULL,
    OrderDate TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    Status VARCHAR(50) DEFAULT 'в обработке',
    CustomerID INT,
    Quantity INT NOT NULL,
    PricePerUnit NUMERIC(10, 2) NOT NULL,
    TotalPrice NUMERIC(10, 2) GENERATED ALWAYS AS (Quantity * PricePerUnit) STORED,
    PRIMARY KEY (OrderID, ProductID),
    FOREIGN KEY (ProductID) REFERENCES Catalog(ProductID) ON DELETE CASCADE
);