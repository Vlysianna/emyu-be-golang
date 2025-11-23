.PHONY: help build run dev test clean install db-create db-migrate db-fresh db-seed format lint setup watch

help:
	@echo "Emyu E-Commerce API - Available Commands:"
	@echo ""
	@echo "  make install          - Install dependencies"
	@echo "  make build            - Build binary"
	@echo "  make run              - Run the application"
	@echo "  make dev              - Run in development mode"
	@echo "  make test             - Run tests"
	@echo "  make clean            - Clean build files"
	@echo "  make db-create        - Create database"
	@echo "  make db-migrate       - Run migrations"
	@echo "  make db-fresh         - Migrate fresh (drop & recreate database)"
	@echo "  make db-seed          - Seed database with sample data"
	@echo "  make db-fresh-seed    - Fresh database + seed data (recommended for dev)"
	@echo "  make setup-fresh      - Complete setup: install, db-fresh, db-seed"
	@echo "  make format           - Format code"
	@echo "  make lint             - Run linter"
	@echo "  make watch            - Watch mode (auto-reload on file changes)"
	@echo ""
	@echo "Quick Start:"
	@echo "  make setup-fresh && make dev"

install:
	go mod download
	go mod tidy

build:
	go build -o bin/api cmd/api/main.go

run: build
	./bin/api

dev:
	go run cmd/api/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/
	go clean

db-create:
	/Applications/XAMPP/bin/mysql -u root -e "CREATE DATABASE emyu CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

db-migrate:
	/Applications/XAMPP/bin/mysql -u root emyu < database/schema.sql
	@echo "✅ Migrations completed!"

db-fresh:
	@echo "🔄 Running migration fresh..."
	/Applications/XAMPP/bin/mysql -u root -e "DROP DATABASE IF EXISTS emyu;"
	/Applications/XAMPP/bin/mysql -u root -e "CREATE DATABASE emyu CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
	/Applications/XAMPP/bin/mysql -u root emyu < database/schema.sql
	@echo "✅ Database fresh migration completed!"

db-seed: build
	@echo "🌱 Seeding database..."
	./bin/api --seed
	@echo "✅ Database seeding completed!"

db-fresh-seed: db-fresh db-seed
	@echo "✅ Database fresh migration and seeding completed!"

format:
	go fmt ./...

lint:
	go vet ./...

# Development commands
setup: install db-create db-migrate
	@echo "✅ Setup complete! Run 'make dev' to start the server"

setup-fresh: install db-fresh db-seed
	@echo "✅ Setup with fresh data complete! Run 'make dev' to start the server"

watch:
	find . -name "*.go" -type f | entr -r go run cmd/api/main.go
