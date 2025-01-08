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
