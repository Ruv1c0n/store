#!/bin/bash

# Проверка наличия jq
if ! command -v jq &> /dev/null; then
    echo "jq не установлен. Пожалуйста, установите jq и попробуйте снова."
    exit 1
fi

# Чтение конфигурации из JSON файла
CONFIG_FILE="database.config"

# Проверка существования файла конфигурации
if [ ! -f "$CONFIG_FILE" ]; then
    echo "Файл конфигурации не найден!"
    exit 1
fi

# Используем jq для извлечения значений
export DB_NAME=$(jq -r '.dbName' "$CONFIG_FILE")
export SERVER_NAME=$(jq -r '.serverName' "$CONFIG_FILE")
export SERVER_PASSWORD=$(jq -r '.serverPassword' "$CONFIG_FILE")
export SERVER_HOSTNAME=$(jq -r '.serverHostName' "$CONFIG_FILE")
export SERVER_PORT=$(jq -r '.serverPort' "$CONFIG_FILE")

# Формируем строку подключения
export CATALOG_DB_URL="postgres://${SERVER_NAME}:${SERVER_PASSWORD}@${SERVER_HOSTNAME}:${SERVER_PORT}/${DB_NAME}?sslmode=disable&x-migrations-table=catalog_migrations"
export ORDER_DB_URL="postgres://${SERVER_NAME}:${SERVER_PASSWORD}@${SERVER_HOSTNAME}:${SERVER_PORT}/${DB_NAME}?sslmode=disable&x-migrations-table=order_migrations"

echo "Конфигурация загружена успешно."