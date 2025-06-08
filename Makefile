# Makefile untuk mengelola Docker commands

.PHONY: help build up down restart logs clean dev prod dev-build dev-up dev-down dev-logs

# Color codes untuk output yang lebih menarik
GREEN := \033[0;32m
YELLOW := \033[1;33m
RED := \033[0;31m
NC := \033[0m # No Color

# Default target
help:
	@echo "$(GREEN)Available commands:$(NC)"
	@echo "$(YELLOW)Production:$(NC)"
	@echo "  build      - Build all Docker images"
	@echo "  up         - Start all services"
	@echo "  down       - Stop all services"
	@echo "  restart    - Restart all services"
	@echo "  logs       - Show logs for all services"
	@echo "  clean      - Remove all containers and volumes"
	@echo "  prod       - Start production environment"
	@echo "  psql-prod  - Execute psql in production postgres"
	@echo "  redis-prod - Execute redis-cli in production redis"
	@echo ""
	@echo "$(YELLOW)Development (dengan Air hot reload):$(NC)"
	@echo "  dev-build  - Build development images"
	@echo "  dev-up     - Start development environment dengan hot reload"
	@echo "  dev-down   - Stop development environment"
	@echo "  dev-logs   - Show development logs"
	@echo "  dev-clean  - Clean development environment"
	@echo "  dev-restart- Restart development environment"
	@echo "  psql       - Execute psql in development postgres"
	@echo "  redis-cli  - Execute redis-cli in development redis"
	@echo ""
	@echo "$(YELLOW)Utilities:$(NC)"
	@echo "  app-logs   - Show only app logs"
	@echo "  db-logs    - Show only database logs"
	@echo "  status     - Check service status"
	@echo "  shell      - Execute bash in app container"

# Production Commands
build:
	@echo "$(GREEN)Building production images...$(NC)"
	docker compose build --no-cache

up:
	@echo "$(GREEN)Starting production services...$(NC)"
	docker compose up -d

down:
	@echo "$(RED)Stopping production services...$(NC)"
	docker compose down

restart:
	@echo "$(YELLOW)Restarting production services...$(NC)"
	docker compose restart

logs:
	@echo "$(GREEN)Showing production logs...$(NC)"
	docker compose logs -f

clean:
	@echo "$(RED)Cleaning production environment...$(NC)"
	docker compose down -v --remove-orphans
	docker system prune -f
	docker volume prune -f

prod:
	@echo "$(GREEN)Starting production environment...$(NC)"
	docker compose up -d

# Production utilities
psql-prod:
	@echo "$(GREEN)Connecting to production PostgreSQL...$(NC)"
	docker compose exec postgres psql -U postgres -d apirijig_v2

redis-prod:
	@echo "$(GREEN)Connecting to production Redis...$(NC)"
	docker compose exec redis redis-cli

# Development Commands (dengan Air hot reload)
dev-build:
	@echo "$(GREEN)Building development images dengan Air...$(NC)"
	docker compose -f docker-compose.dev.yml build --no-cache

dev-up:
	@echo "$(GREEN)Starting development environment dengan Air hot reload...$(NC)"
	docker compose -f docker-compose.dev.yml up -d
	@echo "$(GREEN)Development services started!$(NC)"
	@echo "$(YELLOW)API Server: http://localhost:7000$(NC)"
	@echo "$(YELLOW)PostgreSQL: localhost:5433$(NC)"
	@echo "$(YELLOW)Redis: localhost:6378$(NC)"
	@echo "$(YELLOW)pgAdmin: http://localhost:8080 (admin@rijig.com / admin123)$(NC)"
	@echo "$(YELLOW)Redis Commander: http://localhost:8081$(NC)"
	@echo ""
	@echo "$(GREEN)✨ Hot reload is active! Edit your Go files and see changes automatically ✨$(NC)"

dev-down:
	@echo "$(RED)Stopping development services...$(NC)"
	docker compose -f docker-compose.dev.yml down

dev-logs:
	@echo "$(GREEN)Showing development logs...$(NC)"
	docker compose -f docker-compose.dev.yml logs -f

dev-clean:
	@echo "$(RED)Cleaning development environment...$(NC)"
	docker compose -f docker-compose.dev.yml down -v --remove-orphans
	docker system prune -f

dev-restart:
	@echo "$(YELLOW)Restarting development services...$(NC)"
	docker compose -f docker-compose.dev.yml restart

# Development utilities (FIXED - menggunakan -f docker-compose.dev.yml)
dev-app-logs:
	@echo "$(GREEN)Showing development app logs...$(NC)"
	docker compose -f docker-compose.dev.yml logs -f app

dev-db-logs:
	@echo "$(GREEN)Showing development database logs...$(NC)"
	docker compose -f docker-compose.dev.yml logs -f postgres

dev-shell:
	@echo "$(GREEN)Accessing development app container...$(NC)"
	docker compose -f docker-compose.dev.yml exec app sh

dev-status:
	@echo "$(GREEN)Development service status:$(NC)"
	docker compose -f docker-compose.dev.yml ps

# FIXED: Development database access (default untuk development)
psql:
	@echo "$(GREEN)Connecting to development PostgreSQL...$(NC)"
	docker compose -f docker-compose.dev.yml exec postgres psql -U postgres -d apirijig_v2

redis-cli:
	@echo "$(GREEN)Connecting to development Redis...$(NC)"
	docker compose -f docker-compose.dev.yml exec redis redis-cli

# Shared utilities (default ke production)
app-logs:
	docker compose logs -f app

db-logs:
	docker compose logs -f postgres

status:
	docker compose ps

shell:
	docker compose exec app sh

# Rebuild and restart app only
app-rebuild:
	docker compose build app
	docker compose up -d app

# View real-time resource usage
stats:
	docker stats

# Quick development setup (recommended)
dev:
	@echo "$(GREEN)Setting up complete development environment...$(NC)"
	make dev-build
	make dev-up