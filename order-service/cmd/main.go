package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"net"
	"store/order-service/internal/client"
	"store/order-service/internal/handler"
	db "store/order-service/internal/repository"
	"store/proto"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Config структура для хранения конфигурации
type Config struct {
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
	SSLMode  string
}

// loadConfig загружает конфигурацию из текстового файла
func loadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue // пропускаем пустые строки
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // пропускаем некорректные строки
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "DB_USERNAME":
			config.Username = value
		case "DB_PASSWORD":
			config.Password = value
		case "DB_HOST":
			config.Host = value
		case "DB_PORT":
			config.Port = value
		case "DB_NAME":
			config.DBName = value
		case "DB_SSLMODE":
			config.SSLMode = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return config, nil
}

// generateDBURL генерирует строку подключения к базе данных
func generateDBURL(cfg Config, dbname bool) string {
	if dbname {
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)
	} else {
		return fmt.Sprintf("postgres://%s:%s@%s:%s?sslmode=%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.SSLMode)
	}
}

// createDatabaseIfNotExists создаёт базу данных, если её ещё нет
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
	// Загружаем конфигурацию
	config, err := loadConfig("config.txt") // Укажите путь к вашему текстовому файлу
	if err != nil {
		log.Fatalf("Failed to load config: %v\n", err)
	}

	// Генерируем строку подключения
	dbURL := generateDBURL(*config, false)

	// Создаём базу данных, если её нет
	if err := createDatabaseIfNotExists(dbURL, config.DBName); err != nil {
		log.Fatalf("Failed to create database: %v\n", err)
	}

	// Подключаемся к базе данных
	dbURLWithDB := generateDBURL(*config, true)
	conn, err := pgx.Connect(context.Background(), dbURLWithDB)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())
	fmt.Println("Connected to PostgreSQL!")

	// Применяем миграции
	if err := runMigrations(dbURLWithDB+"&x-migrations-table=order_migrations"); err != nil {
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
