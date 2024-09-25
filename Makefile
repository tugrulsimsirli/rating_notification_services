# Makefile for Docker Compose management

# Docker Compose commands
DOCKER_COMPOSE = docker compose

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
build:
	$(DOCKER_COMPOSE) down -v
	$(DOCKER_COMPOSE) up --build
