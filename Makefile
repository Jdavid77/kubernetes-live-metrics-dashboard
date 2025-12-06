.PHONY: help build run clean test docker-build docker-up docker-down frontend-install backend-deps

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

# Backend targets
backend-deps: ## Install backend dependencies
	cd backend && go mod download

backend-build: ## Build backend binary
	cd backend && go build -o api ./cmd/api

backend-run: ## Run backend server
	cd backend && go run cmd/api/main.go

backend-test: ## Run backend tests
	cd backend && go test -v ./...

# Frontend targets
frontend-install: ## Install frontend dependencies
	cd frontend && npm install

frontend-dev: ## Run frontend development server
	cd frontend && npm start

frontend-build: ## Build frontend for production
	cd frontend && npm run build

frontend-test: ## Run frontend tests
	cd frontend && npm test

# Docker targets
docker-build: ## Build Docker images
	docker-compose build

docker-up: ## Start Docker containers
	docker-compose up -d

docker-down: ## Stop Docker containers
	docker-compose down

docker-logs: ## View Docker logs
	docker-compose logs -f

docker-restart: ## Restart Docker containers
	docker-compose restart

# Combined targets
build: backend-build frontend-build ## Build both backend and frontend

install: backend-deps frontend-install ## Install all dependencies

run: ## Run both backend and frontend in development mode
	@echo "Starting backend and frontend..."
	@make -j2 backend-run frontend-dev

clean: ## Clean build artifacts
	rm -rf backend/api
	rm -rf frontend/build
	rm -rf frontend/node_modules

# Kubernetes targets
k8s-apply-rbac: ## Apply Kubernetes RBAC configuration
	kubectl apply -f rbac.yaml

k8s-delete-rbac: ## Delete Kubernetes RBAC configuration
	kubectl delete -f rbac.yaml
