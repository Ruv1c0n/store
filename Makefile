# Variables
GO_CMD = go
GOOSE_CMD = goose
PROTOC_CMD = protoc
BUILD_DIR = bin
CATALOG_SERVICE_DIR = ./catalog-service
ORDER_SERVICE_DIR = ./order-service
CATALOG_BINARY = $(BUILD_DIR)/catalog-service
ORDER_BINARY = $(BUILD_DIR)/order-service
CATALOG_PORT = 50051
ORDER_PORT = 50052

# Path to proto files
CATALOG_PROTO_FILES = $(CATALOG_SERVICE_DIR)/internal/proto/catalog.proto
ORDER_PROTO_FILES = $(ORDER_SERVICE_DIR)/internal/proto/order.proto

# Repository files
CATALOG_REPO_FILES = $(CATALOG_SERVICE_DIR)/internal/repository/db.go
ORDER_REPO_FILES = $(ORDER_SERVICE_DIR)/internal/repository/db.go

# Default goal
.DEFAULT_GOAL := help

# Create directory for binaries
$(BUILD_DIR):
	mkdir $(BUILD_DIR)

# Generate Go files from proto for CatalogService
init-proto-catalog: $(CATALOG_PROTO_FILES) ## Initialize proto files for CatalogService
	@echo "Generating Go files from proto for CatalogService..."
	$(PROTOC_CMD) --go_out=$(CATALOG_SERVICE_DIR) --go-grpc_out=$(CATALOG_SERVICE_DIR) $(CATALOG_PROTO_FILES)

# Generate Go files from proto for OrderService
init-proto-order: $(ORDER_PROTO_FILES) ## Initialize proto files for OrderService
	@echo "Generating Go files from proto for OrderService..."
	$(PROTOC_CMD) --go_out=$(ORDER_SERVICE_DIR) --go-grpc_out=$(ORDER_SERVICE_DIR) $(ORDER_PROTO_FILES)

# Initialize proto files for both services
init-proto: init-proto-catalog init-proto-order ## Initialize proto files for all services

# Build CatalogService
build-catalog: $(BUILD_DIR) init-proto $(CATALOG_REPO_FILES) ## Build CatalogService
	$(GO_CMD) build -o $(CATALOG_BINARY) $(CATALOG_SERVICE_DIR)/cmd

# Build OrderService
build-order: $(BUILD_DIR) init-proto $(ORDER_REPO_FILES) ## Build OrderService
	$(GO_CMD) build -o $(ORDER_BINARY) $(ORDER_SERVICE_DIR)/cmd

# Build all services
build: build-catalog build-order ## Build all services

# Run CatalogService
run-catalog: build-catalog ## Run CatalogService
	$(CATALOG_BINARY)

# Run OrderService
run-order: build-order ## Run OrderService
	$(ORDER_BINARY)

# Run all services
run: ## Run all services
	@echo "Starting CatalogService on port $(CATALOG_PORT)..."
	@start cmd /c $(CATALOG_BINARY)
	@echo "Starting OrderService on port $(ORDER_PORT)..."
	@start cmd /c $(ORDER_BINARY)

# Migrations for CatalogService
migrate-catalog-up: ## Apply migrations for CatalogService
	$(GOOSE_CMD) -dir $(CATALOG_SERVICE_DIR)/migrations postgres "user=your_user password=your_password dbname=catalog_db sslmode=disable" up

migrate-catalog-down: ## Rollback migrations for CatalogService
	$(GOOSE_CMD) -dir $(CATALOG_SERVICE_DIR)/migrations postgres "user=your_user password=your_password dbname=catalog_db sslmode=disable" down

# Migrations for OrderService
migrate-order-up: ## Apply migrations for OrderService
	$(GOOSE_CMD) -dir $(ORDER_SERVICE_DIR)/migrations postgres "user=your_user password=your_password dbname=order_db sslmode=disable" up

migrate-order-down: ## Rollback migrations for OrderService
	$(GOOSE_CMD) -dir $(ORDER_SERVICE_DIR)/migrations postgres "user=your_user password=your_password dbname=order_db sslmode=disable" down

# Apply migrations for all services
migrate-up: migrate-catalog-up migrate-order-up ## Apply migrations for all services

# Rollback migrations for all services
migrate-down: migrate-catalog-down migrate-order-down ## Rollback migrations for all services

# Clean binaries
clean: ## Remove compiled binaries
	rmdir /s /q $(BUILD_DIR)

# Help
help: ## Show help
	@echo "Available commands:"
	@echo "  init-proto          Initialize proto files"
	@echo "  build-catalog       Build CatalogService"
	@echo "  build-order         Build OrderService"
	@echo "  build               Build all services"
	@echo "  run-catalog         Run CatalogService"
	@echo "  run-order           Run OrderService"
	@echo "  run                 Run all services"
	@echo "  migrate-catalog-up  Apply migrations for CatalogService"
	@echo "  migrate-catalog-down Rollback migrations for CatalogService"
	@echo "  migrate-order-up    Apply migrations for OrderService"
	@echo "  migrate-order-down  Rollback migrations for OrderService"
	@echo "  migrate-up          Apply migrations for all services"
	@echo "  migrate-down        Rollback migrations for all services"
	@echo "  clean               Remove compiled binaries"
