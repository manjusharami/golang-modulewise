 # ==========================================
# Student CRUD - Makefile
# ==========================================

APP=student-crud

DB_CONTAINER=student-postgres
DB_USER=postgres
DB_PASS=postgres
DB_NAME=studentdb
DB_PORT=5432

.PHONY: all postgres wait table deps build run clean

all: postgres wait table deps build run

postgres:
	@echo "Starting PostgreSQL container..."
	-docker rm -f $(DB_CONTAINER) >/dev/null 2>&1
	docker run -d \
		--name $(DB_CONTAINER) \
		-e POSTGRES_USER=$(DB_USER) \
		-e POSTGRES_PASSWORD=$(DB_PASS) \
		-e POSTGRES_DB=$(DB_NAME) \
		-p $(DB_PORT):5432 \
		postgres:16

wait:
	@echo "Waiting for PostgreSQL to start..."
	@sleep 8

table:
	@echo "Creating students table..."
	@echo "CREATE TABLE IF NOT EXISTS students (\
id SERIAL PRIMARY KEY,\
name VARCHAR(100) NOT NULL,\
age INT,\
email VARCHAR(100) UNIQUE\
);" | docker exec -i $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME)

deps:
	@echo "Installing Go dependencies..."
	go mod tidy

build:
	@echo "Building application..."
	go build -o $(APP)

run:
	@echo "Starting application..."
	./$(APP)

clean:
	@echo "Cleaning up..."
	-docker rm -f $(DB_CONTAINER)
	-rm -f $(APP)
	go clean
