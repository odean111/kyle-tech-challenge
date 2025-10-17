.PHONY: help up down build logs clean migrate test

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

up: ## Start all services
	docker-compose up -d

down: ## Stop all services
	docker-compose down

build: ## Build all services
	docker-compose build

logs: ## Show logs for all services
	docker-compose logs -f

logs-backend: ## Show backend logs
	docker-compose logs -f backend

logs-frontend: ## Show frontend logs
	docker-compose logs -f frontend

logs-postgres: ## Show postgres logs
	docker-compose logs -f postgres

clean: ## Clean up containers and volumes
	docker-compose down -v
	docker system prune -f

migrate: ## Run database migrations
	docker-compose exec backend sh -c "cd migrations && sqitch deploy"

migrate-status: ## Check migration status
	docker-compose exec backend sh -c "cd migrations && sqitch status"

migrate-revert: ## Revert last migration
	docker-compose exec backend sh -c "cd migrations && sqitch revert --to @HEAD^"

test-backend: ## Test backend connection
	curl http://localhost:8080/api/v1/

dev: up ## Start development environment
	@echo "Starting development environment..."
	@echo "Backend: http://localhost:8080"
	@echo "Frontend: http://localhost:5173"
	@echo "API Status: http://localhost:8080/api/v1/"