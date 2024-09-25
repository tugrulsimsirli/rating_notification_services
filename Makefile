# Makefile for Docker Compose management

# Docker Compose commands
DOCKER_COMPOSE = docker compose
MIN_COVERAGE = 100

# Targets
.PHONY: up down restart clean build

# Start containers
up:
	$(DOCKER_COMPOSE) up --build

# Stop and remove containers, networks, volumes, and images
down:
	$(DOCKER_COMPOSE) down

# Restart the services
restart: down up

# Stop and remove containers, networks, volumes (clean)
clean:
	$(DOCKER_COMPOSE) down -v

# Build or rebuild services
build: test
	$(DOCKER_COMPOSE) down -v
	$(DOCKER_COMPOSE) up --build

test:
	@echo "Running tests for notification_service..."
	@cd notification_service && go test ./internal/app/services -coverprofile=coverage.out
	@echo "Running tests for rating_service..."
	@cd rating_service && go test ./internal/app/services -coverprofile=coverage.out
