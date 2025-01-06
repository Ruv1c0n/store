Микросервис онлайн магазина

```
store-main
├─ .gitignore
├─ .idea
│  ├─ .gitignore
│  ├─ modules.xml
│  ├─ store-main.iml
│  └─ workspace.xml
├─ .vscode
│  └─ settings.json
├─ catalog-service
│  ├─ cmd
│  │  └─ main.go
│  ├─ internal
│  │  ├─ handler
│  │  │  ├─ catalog_handler.go
│  │  │  └─ handler_test.go
│  │  ├─ proto
│  │  │  ├─ catalog.pb.go
│  │  │  ├─ catalog.proto
│  │  │  └─ catalog_grpc.pb.go
│  │  ├─ repository
│  │  │  ├─ db.go
│  │  │  └─ mock
│  │  │     └─ mock.go
│  │  └─ service
│  │     └─ service.go
│  └─ migrations
│     └─ 20250104120000_create_products_table.up.sql
├─ Db
├─ go.mod
├─ go.sum
├─ Makefile
├─ order-service
│  ├─ cmd
│  │  └─ main.go
│  ├─ internal
│  │  ├─ handler
│  │  │  └─ order_handler.go
│  │  ├─ proto
│  │  │  ├─ order.pb.go
│  │  │  ├─ order.proto
│  │  │  └─ order_grpc.pb.go
│  │  ├─ repository
│  │  │  └─ db.go
│  │  └─ service
│  │     └─ service.go
│  └─ migrations
│     └─ 20250104121000_create_orders_table.sql
└─ README.md

```