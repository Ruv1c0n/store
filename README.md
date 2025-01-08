# Микросервис онлайн магазина
## Архитектура проекта
```
store
├─ catalog-service
│  ├─ cmd
│  │  └─ main.go
│  ├─ internal
│  │  ├─ handler
│  │  │  ├─ catalog_handler.go
│  │  │  └─ handler_test.go
│  │  └─ repository
│  │     ├─ mock
│  │     │  └─ mock.go
│  │     └─ db.go
│  └─ migrations
│     ├─ 20250104120000_create_products_table.down.sql
│     └─ 20250104120000_create_products_table.up.sql
├─ order-service
│  ├─ cmd
│  │  └─ main.go
│  ├─ internal
│  │  ├─ client
│  │  │  └─ catalog_client.go
│  │  ├─ handler
│  │  │  └─ order_handler.go
│  │  │  └─ order_handler_test.go
│  │  └─ repository
│  │     ├─ mock
│  │     │  └─ mock.go
│  │     └─ db.go
│  └─ migrations
│     ├─ 20250104120000_create_orders_table.down.sql
│     └─ 20250104121000_create_orders_table.up.sql
├─ proto
│  └─ catalog.proto
│  └─ order.proto
├─ .gitignore
├─  config.txt
├─  go.mod
├─  go.sum
├─  Makefile
└─ README.md
```
## Как собирать проект
#### Вызов всех make команд
```
make | make help
```

#### Инициализация проекта
```
make init-proto
```

#### Cборка проекта
##### Отдельных сервисов
```
make build-catalog | make build-order
```
##### Всего проекта
```
make build
```

#### Запуск сервисов
```
make run-catalog | make run-order
```

#### Миграции
```
make migrate-/catalog|order| /-/up|down/
```

#### Очистка бинарников(для пересборки проекта)
```
make clean
```

## Как тестировать проект
#### Целиком
```
go test -v -coverprofile=./...
```

#### Отдельные части
Тесты прописаны для handler'ов каждого из сервисов
```
cd ./путь/до/тестового/файла
go test -v -cover
```

## Тестирование работы сервисов
#### Для CATALOG
- Создаание товаров в каталоге
```
grpcurl -plaintext -d '{\"product_name\": \"Кофеварка\", \"stock_quantity\": 15, \"price_per_unit\": 9500}' localhost:50051 catalog.ProductService/AddProduct
grpcurl -plaintext -d '{\"product_name\": \"Чайник\", \"stock_quantity\": 100, \"price_per_unit\": 5700}' localhost:50051 catalog.ProductService/AddProduct
grpcurl -plaintext -d '{\"product_name\": \"Кастрюля\", \"stock_quantity\": 31, \"price_per_unit\": 3700}' localhost:50051 catalog.ProductService/AddProduct
grpcurl -plaintext -d '{\"product_name\": \"Кружка\", \"stock_quantity\": 100, \"price_per_unit\": 600}' localhost:50051 catalog.ProductService/AddProduct
grpcurl -plaintext -d '{\"product_name\": \"Набор столовых приборов\", \"stock_quantity\": 17, \"price_per_unit\": 7000}' localhost:50051 catalog.ProductService/AddProduct
```
- Вывод всех продуктов
```
grpcurl -plaintext localhost:50051 catalog.ProductService/GetAllProducts
```
- Обновление чайника(стоимость)
```
grpcurl -plaintext -d '{\"product_id\": 2, \"product_name\": \"Чайник\", \"stock_quantity\": 100, \"price_per_unit\": 4700}' localhost:50051 catalog.ProductService/UpdateProduct
```
- Вывод по ИД чайник(показываем изменения)
```
grpcurl -plaintext -d '{\"product_id\": 2}' localhost:50051 catalog.ProductService/GetProductByID
```
- Удаляем кофеварку
```
grpcurl -plaintext -d '{\"product_id\": 1}' localhost:50051 catalog.ProductService/DeleteProduct
```
-----------------------------------------

#### Для ORDER
- Создание Заказа(Два чайника)
```
grpcurl -plaintext -d '{\"customer_id\": 1, \"items\": [{\"product_id\": 2, \"quantity\": 2}]}' localhost:50052 order.OrderService/CreateOrder
```
- Создание Заказа(Три кружки)
```
grpcurl -plaintext -d '{\"customer_id\": 1, \"items\": [{\"product_id\": 4, \"quantity\": 6}]}' localhost:50052 order.OrderService/CreateOrder
```
- Вывод двух заказов
```
grpcurl -plaintext localhost:50052 order.OrderService/GetAllOrders
```
- Изменился статус заказа 1
```
grpcurl -plaintext -d '{\"order_id\": 1, \"status\": \"Выполнен\"}' localhost:50052 order.OrderService/UpdateOrder
```
- Вывод заказа по ИД
```
grpcurl -plaintext -d '{\"order_id\": 1}' localhost:50052 order.OrderService/GetOrderByID
```
- Удалили заказ на кружки
```
grpcurl -plaintext -d '{\"order_id\": 2}' localhost:50052 order.OrderService/DeleteOrder
```
- Вывод одного заказа
```
grpcurl -plaintext localhost:50052 order.OrderService/GetAllOrders
```
