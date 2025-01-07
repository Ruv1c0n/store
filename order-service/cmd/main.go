package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"store/order-service/internal/client" // Импортируем пакет client
	"store/order-service/internal/handler"
	db "store/order-service/internal/repository"
	"store/proto"
)

func createDatabaseIfNotExists(dbURL, dbName string) error {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow("SELECT EXISTS (SELECT FROM pg_database WHERE datname = $1)", dbName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	if !exists {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		log.Printf("Database '%s' created successfully.\n", dbName)
	} else {
		log.Printf("Database '%s' already exists.\n", dbName)
	}

	return nil
}

func runMigrations(databaseURL string) error {
	m, err := migrate.New(
		"file://order-service/migrations", // Путь к миграциям
		databaseURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Println("Migrations applied successfully!")
	return nil
}

func main() {
	// Параметры подключения
	dbURL := "postgres://postgres:0000@localhost:5432?sslmode=disable"
	dbName := "catalog" // Используем базу данных "catalog"

	// Создаём базу данных, если её нет
	if err := createDatabaseIfNotExists(dbURL, dbName); err != nil {
		log.Fatalf("Failed to create database: %v\n", err)
	}

	// Подключаемся к базе данных
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:0000@localhost:5432/catalog?sslmode=disable")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())
	fmt.Println("Connected to PostgreSQL!")

	// Применяем миграции
	if err := runMigrations("postgres://postgres:0000@localhost:5432/catalog?sslmode=disable&x-migrations-table=order_migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v\n", err)
	}
	fmt.Println("Migrations applied successfully!")

	// Создаем клиент для CatalogService
	catalogClient, err := client.NewCatalogClient("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to create catalog client: %v", err)
	}
	defer catalogClient.Close()

	// Создаем экземпляр OrderDB
	orderDB := db.NewOrderDB(conn, catalogClient)

	// Создаем новый gRPC сервер
	grpcServer := grpc.NewServer()

	// Регистрируем обработчик
	orderHandler := handler.NewOrderHandler(orderDB)
	proto.RegisterOrderServiceServer(grpcServer, orderHandler)

	// Включаем Reflection
	reflection.Register(grpcServer)

	// Запускаем сервер на порту 50052
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}

	log.Println("gRPC сервер запущен на порту 50052...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка при работе сервера: %v", err)
	}
}
