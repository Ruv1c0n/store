-- Создание последовательности для генерации order_id
CREATE SEQUENCE orders_orderid_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

-- Создание таблицы Orders
CREATE TABLE Orders (
    OrderID         INT             NOT NULL    DEFAULT nextval('orders_orderid_seq'),
    ProductID       INT             NOT NULL,
    CustomerID      INT             NOT NULL,
    Quantity        INT             NOT NULL,
    PricePerUnit    NUMERIC(10, 2)  NOT NULL,
    OrderDate       TIMESTAMP       NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    Status          VARCHAR(50)                 DEFAULT 'в обработке', 
    PRIMARY KEY (OrderID, ProductID)
);