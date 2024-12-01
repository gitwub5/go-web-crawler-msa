.PHONY: help build up down logs ps test clean
.DEFAULT_GOAL := help

DOCKER_COMPOSE_FILE := docker-compose.yml

# Docker Compose targets
build: ## Build docker images for all services
	docker compose -f $(DOCKER_COMPOSE_FILE) build

up: ## Start all services defined in docker-compose.yml
	docker compose -f $(DOCKER_COMPOSE_FILE) up -d

down: ## Stop all services defined in docker-compose.yml
	docker compose -f $(DOCKER_COMPOSE_FILE) down

logs: ## Tail logs for all services
	docker compose -f $(DOCKER_COMPOSE_FILE) logs -f

ps: ## Show container status
	docker compose -f $(DOCKER_COMPOSE_FILE) ps

test: ## Run tests for the project
	go test -race -shuffle=on ./...

clean: ## Remove unused Docker objects
	docker system prune -f

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'