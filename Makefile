# Makefile untuk mengelola Docker commands - Optimized Version

.PHONY: help build up down restart logs clean dev prod dev-build dev-up dev-down dev-logs

# Color codes untuk output yang lebih menarik
GREEN := \033[0;32m
YELLOW := \033[1;33m
RED := \033[0;31m
BLUE := \033[0;34m
PURPLE := \033[0;35m
CYAN := \033[0;36m
NC := \033[0m # No Color

# Project variables
PROJECT_NAME := rijig_backend
DEV_COMPOSE_FILE := docker-compose.dev.yml

# Default target
help:
	@echo "$(GREEN)🚀 $(PROJECT_NAME) - Available Commands:$(NC)"
	@echo ""
	@echo "$(YELLOW)📦 Development Commands (Hot Reload):$(NC)"
	@echo "  $(CYAN)dev$(NC)          - Complete development setup (build + up)"
	@echo "  $(CYAN)dev-build$(NC)    - Build development images"
	@echo "  $(CYAN)dev-up$(NC)       - Start development environment"
	@echo "  $(CYAN)dev-down$(NC)     - Stop development environment"
	@echo "  $(CYAN)dev-restart$(NC)  - Restart development services"
	@echo "  $(CYAN)dev-logs$(NC)     - Show development logs (all services)"
	@echo "  $(CYAN)dev-clean$(NC)    - Clean development environment"
	@echo ""
	@echo "$(YELLOW)🛠️  Development Utilities:$(NC)"
	@echo "  $(CYAN)dev-app-logs$(NC) - Show only app logs"
	@echo "  $(CYAN)dev-db-logs$(NC)  - Show only database logs"
	@echo "  $(CYAN)dev-shell$(NC)    - Access app container shell"
	@echo "  $(CYAN)dev-status$(NC)   - Check development services status"
	@echo "  $(CYAN)psql$(NC)         - Connect to development PostgreSQL"
	@echo "  $(CYAN)redis-cli$(NC)    - Connect to development Redis"
	@echo ""
	@echo "$(YELLOW)🧹 Maintenance:$(NC)"
	@echo "  $(RED)clean-all$(NC)     - Clean everything (containers, volumes, images)"
	@echo "  $(RED)system-prune$(NC)  - Clean Docker system"
	@echo "  $(CYAN)stats$(NC)         - Show container resource usage"

# ======================
# DEVELOPMENT COMMANDS
# ======================

# Quick development setup (recommended)
dev: dev-build dev-up
	@echo "$(GREEN)✨ Development environment ready!$(NC)"
	@echo "$(BLUE)🌐 Services:$(NC)"
	@echo "  • API Server: $(CYAN)http://localhost:7000$(NC)"
	@echo "  • PostgreSQL: $(CYAN)localhost:5433$(NC)"
	@echo "  • Redis: $(CYAN)localhost:6378$(NC)"
	@echo "  • pgAdmin: $(CYAN)http://localhost:8080$(NC) (admin@rijig.com / admin123)"
	@echo "  • Redis Commander: $(CYAN)http://localhost:8081$(NC)"
	@echo ""
	@echo "$(GREEN)🔥 Hot reload is active! Edit your Go files and see changes automatically$(NC)"

dev-build:
	@echo "$(YELLOW)🔨 Building development images...$(NC)"
	@docker compose -f $(DEV_COMPOSE_FILE) build --no-cache
	@echo "$(GREEN)✅ Development images built successfully!$(NC)"

dev-up:
	@echo "$(YELLOW)🚀 Starting development services...$(NC)"
	@docker compose -f $(DEV_COMPOSE_FILE) up -d
	@echo "$(GREEN)✅ Development services started!$(NC)"
	@make dev-status

dev-down:
	@echo "$(RED)🛑 Stopping development services...$(NC)"
	@docker compose -f $(DEV_COMPOSE_FILE) down
	@echo "$(GREEN)✅ Development services stopped!$(NC)"

dev-restart:
	@echo "$(YELLOW)🔄 Restarting development services...$(NC)"
	@docker compose -f $(DEV_COMPOSE_FILE) restart
	@echo "$(GREEN)✅ Development services restarted!$(NC)"

dev-logs:
	@echo "$(CYAN)📋 Showing development logs (Ctrl+C to exit)...$(NC)"
	@docker compose -f $(DEV_COMPOSE_FILE) logs -f --tail=100

dev-clean:
	@echo "$(RED)🧹 Cleaning development environment...$(NC)"
	@docker compose -f $(DEV_COMPOSE_FILE) down -v --remove-orphans
	@echo "$(GREEN)✅ Development environment cleaned!$(NC)"

# ======================
# DEVELOPMENT UTILITIES
# ======================

dev-app-logs:
	@echo "$(CYAN)📋 Showing app logs (Ctrl+C to exit)...$(NC)"
	@docker compose -f $(DEV_COMPOSE_FILE) logs -f --tail=50 app

dev-db-logs:
	@echo "$(CYAN)📋 Showing database logs (Ctrl+C to exit)...$(NC)"
	@docker compose -f $(DEV_COMPOSE_FILE) logs -f --tail=50 postgres

dev-shell:
	@echo "$(CYAN)🐚 Accessing app container shell...$(NC)"
	@docker compose -f $(DEV_COMPOSE_FILE) exec app sh

dev-status:
	@echo "$(BLUE)📊 Development services status:$(NC)"
	@docker compose -f $(DEV_COMPOSE_FILE) ps

psql:
	@echo "$(CYAN)🐘 Connecting to development PostgreSQL...$(NC)"
	@docker compose -f $(DEV_COMPOSE_FILE) exec postgres psql -U postgres -d apirijig_v2

redis-cli:
	@echo "$(CYAN)⚡ Connecting to development Redis...$(NC)"
	@docker compose -f $(DEV_COMPOSE_FILE) exec redis redis-cli

# ======================
# MAINTENANCE COMMANDS
# ======================

clean-all:
	@echo "$(RED)🧹 Performing complete cleanup...$(NC)"
	@docker compose -f $(DEV_COMPOSE_FILE) down -v --remove-orphans 2>/dev/null || true
	@echo "$(YELLOW)🗑️  Removing unused containers, networks, and images...$(NC)"
	@docker system prune -a -f --volumes
	@echo "$(GREEN)✅ Complete cleanup finished!$(NC)"

system-prune:
	@echo "$(YELLOW)🗑️  Cleaning Docker system...$(NC)"
	@docker system prune -f
	@echo "$(GREEN)✅ Docker system cleaned!$(NC)"

stats:
	@echo "$(BLUE)📈 Container resource usage:$(NC)"
	@docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}\t{{.BlockIO}}"

# ======================
# QUICK COMMANDS
# ======================

# App only restart (faster for development)
app-restart:
	@echo "$(YELLOW)🔄 Restarting app container only...$(NC)"
	@docker compose -f $(DEV_COMPOSE_FILE) restart app
	@echo "$(GREEN)✅ App container restarted!$(NC)"

# Check if containers are healthy
health-check:
	@echo "$(BLUE)🏥 Checking container health...$(NC)"
	@docker compose -f $(DEV_COMPOSE_FILE) ps --format "table {{.Name}}\t{{.Status}}"