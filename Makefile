# Variables
DOCKER_COMPOSE=docker-compose
GO_CMD=go
PROJECT_NAME=log-tracer

# Docker Commands
.PHONY: start stop restart build clean logs shell

start:
	@echo "Starting the ${PROJECT_NAME} services..."
	${DOCKER_COMPOSE} up -d

stop:
	@echo "Stopping the ${PROJECT_NAME} services..."
	${DOCKER_COMPOSE} down

restart: stop start

build:
	@echo "Building the Docker containers..."
	${DOCKER_COMPOSE} build

clean:
	@echo "Cleaning up Docker containers, networks, and volumes..."
	${DOCKER_COMPOSE} down --volumes --remove-orphans
	${DOCKER_COMPOSE} rm -fsv

logs:
	@echo "Displaying logs for all services..."
	${DOCKER_COMPOSE} logs -f

shell:
	@echo "Opening a shell in the Kafka container..."
	${DOCKER_COMPOSE} exec kafka /bin/bash

# Go Commands
.PHONY: fmt vet test mod

fmt:
	@echo "Running go fmt..."
	${GO_CMD} fmt ./...

vet:
	@echo "Running go vet..."
	${GO_CMD} vet ./...

test:
	@echo "Running tests..."
	${GO_CMD} test ./...

mod:
	@echo "Running go mod tidy..."
	${GO_CMD} mod tidy

# Combined Commands
.PHONY: dev prod

dev: mod fmt vet build start
	@echo "Development environment started."

prod: mod fmt vet build start
	@echo "Production environment started."

# Utility Commands
.PHONY: kafka-producer kafka-consumer dashboard

kafka-producer:
	@echo "Building the Kafka producer..."
	${GO_CMD} build -o bin/producer ./cmd/producer

kafka-consumer:
	@echo "Building the Kafka consumer..."
	${GO_CMD} build -o bin/consumer ./cmd/consumer

dashboard:
	@echo "Building the dashboard..."
	${GO_CMD} build -o bin/dashboard ./cmd/dashboard

# Full Build
.PHONY: all

all: kafka-producer kafka-consumer dashboard
	@echo "All components built successfully."

