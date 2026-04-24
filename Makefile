# LMS Makefile. `make help` to list targets.
#
# Conventions:
#   SERVICE=<n> scopes per-service targets (e.g. make build SERVICE=identity).
#   Most targets without SERVICE act on all services.

SERVICES := gateway identity submission grading vcs
COMPOSE  := docker compose -f deploy/docker-compose.yml

# Colors for output readability.
BLUE   := \033[36m
YELLOW := \033[33m
GREEN  := \033[32m
RED    := \033[31m
RESET  := \033[0m

.DEFAULT_GOAL := help

# ============================================================
# Help
# ============================================================

.PHONY: help
help: ## Print this help
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make $(BLUE)<target>$(RESET) [SERVICE=<n>]\n\nTargets:\n"} \
		/^## ---/ {printf "\n$(YELLOW)%s$(RESET)\n", substr($$0, 8)} \
		/^[a-zA-Z0-9_.-]+:.*##/ {printf "  $(BLUE)%-22s$(RESET) %s\n", $$1, $$2}' \
		$(MAKEFILE_LIST)

## --- Proto ---

.PHONY: proto
proto: ## Regenerate Go code from .proto files
	cd proto && buf generate

.PHONY: proto-lint
proto-lint: ## Lint .proto files
	cd proto && buf lint

.PHONY: proto-breaking
proto-breaking: ## Check breaking changes against main
	cd proto && buf breaking --against '../.git#branch=main,subdir=proto'

.PHONY: proto-watch
proto-watch: ## Regenerate proto on file changes (needs entr)
	@command -v entr >/dev/null || { echo "install entr: pacman -S entr"; exit 1; }
	@find proto -name '*.proto' | entr -c $(MAKE) proto

## --- Go: build & test ---

.PHONY: tidy
tidy: ## Run go mod tidy in every module
	@for m in pkg gen/go $(SERVICES:%=services/%) test/e2e; do \
		echo "$(BLUE)==>$(RESET) tidy $$m"; \
		(cd $$m && go mod tidy) || exit 1; \
	done

.PHONY: build
build: ## Build all services (or one with SERVICE=<n>)
ifdef SERVICE
	@echo "$(BLUE)==>$(RESET) build $(SERVICE)"
	@cd services/$(SERVICE) && go build -o bin/svc ./cmd
else
	@for s in $(SERVICES); do \
		echo "$(BLUE)==>$(RESET) build $$s"; \
		(cd services/$$s && go build -o bin/svc ./cmd) || exit 1; \
	done
endif

.PHONY: test
test: ## Run unit tests (or one module with SERVICE=<n>)
ifdef SERVICE
	@cd services/$(SERVICE) && go test -race -count=1 ./...
else
	@for m in pkg $(SERVICES:%=services/%); do \
		echo "$(BLUE)==>$(RESET) test $$m"; \
		(cd $$m && go test -race -count=1 ./...) || exit 1; \
	done
endif

.PHONY: test-cover
test-cover: ## Run tests with coverage, open HTML report
	@cd services/$(or $(SERVICE),identity) && \
		go test -race -coverprofile=coverage.out ./... && \
		go tool cover -html=coverage.out -o coverage.html && \
		echo "coverage.html generated"

.PHONY: test-e2e
test-e2e: ## Run e2e tests against live stack
	@echo "$(BLUE)==>$(RESET) ensuring stack is up"
	@$(COMPOSE) up -d --wait
	@echo "$(BLUE)==>$(RESET) running e2e"
	@cd test/e2e && go test -v -count=1 ./... || \
		($(COMPOSE) logs --tail=50 && exit 1)

.PHONY: lint
lint: ## Run golangci-lint on all modules
	golangci-lint run ./...

.PHONY: lint-fix
lint-fix: ## Run golangci-lint with --fix
	golangci-lint run --fix ./...

.PHONY: fmt
fmt: ## Format all Go files
	@gofmt -w -s $$(find . -name '*.go' -not -path './gen/*')
	@echo "$(GREEN)formatted$(RESET)"

.PHONY: generate
generate: proto sqlc ## Run all code generation

.PHONY: sqlc
sqlc: ## Regenerate sqlc code for all services that have sqlc.yaml
	@for s in $(SERVICES); do \
		if [ -f services/$$s/sqlc.yaml ]; then \
			echo "$(BLUE)==>$(RESET) sqlc $$s"; \
			(cd services/$$s && sqlc generate) || exit 1; \
		fi; \
	done

## --- Docker compose ---

.PHONY: up
up: ## Start full stack (only rebuilds changed services)
	$(COMPOSE) up -d --build --wait

.PHONY: up-infra
up-infra: ## Start only infrastructure (postgres, redpanda, minio, jaeger)
	$(COMPOSE) up -d postgres redpanda minio jaeger --wait

.PHONY: up-svc
up-svc: ## Start/restart ONE service. Requires SERVICE=<n>
	@test -n "$(SERVICE)" || { echo "$(RED)SERVICE is required$(RESET)"; exit 1; }
	$(COMPOSE) up -d --build --no-deps $(SERVICE)

.PHONY: down
down: ## Stop containers (keeps volumes)
	$(COMPOSE) down

.PHONY: nuke
nuke: ## Stop AND delete volumes (full reset)
	$(COMPOSE) down -v

.PHONY: restart
restart: ## Restart one service. Requires SERVICE=<n>
	@test -n "$(SERVICE)" || { echo "$(RED)SERVICE is required$(RESET)"; exit 1; }
	$(COMPOSE) restart $(SERVICE)

.PHONY: rebuild
rebuild: ## Rebuild and restart one service. Requires SERVICE=<n>
	@test -n "$(SERVICE)" || { echo "$(RED)SERVICE is required$(RESET)"; exit 1; }
	$(COMPOSE) up -d --build --no-deps $(SERVICE)

.PHONY: logs
logs: ## Tail logs (all or one via SERVICE=<n>)
	$(COMPOSE) logs -f --tail=100 $(SERVICE)

.PHONY: ps
ps: ## Show container status
	$(COMPOSE) ps

.PHONY: stats
stats: ## Show container resource usage
	@docker stats --no-stream $$($(COMPOSE) ps -q)

## --- Dev shortcuts: shell access ---

.PHONY: psql
psql: ## Open psql shell. DB=<n> for specific db (default: lms_identity)
	$(COMPOSE) exec postgres psql -U lms -d $(or $(DB),lms_identity)

.PHONY: psql-list
psql-list: ## List all databases
	$(COMPOSE) exec postgres psql -U lms -l

.PHONY: pg-dump
pg-dump: ## Dump one DB to stdout. DB=<n> required
	@test -n "$(DB)" || { echo "$(RED)DB is required$(RESET)"; exit 1; }
	$(COMPOSE) exec postgres pg_dump -U lms $(DB)

.PHONY: kafka-topics
kafka-topics: ## List Kafka topics
	$(COMPOSE) exec redpanda rpk topic list

.PHONY: kafka-consume
kafka-consume: ## Consume from a topic. TOPIC=<n> required
	@test -n "$(TOPIC)" || { echo "$(RED)TOPIC is required$(RESET)"; exit 1; }
	$(COMPOSE) exec redpanda rpk topic consume $(TOPIC)

.PHONY: kafka-describe
kafka-describe: ## Describe a topic. TOPIC=<n> required
	@test -n "$(TOPIC)" || { echo "$(RED)TOPIC is required$(RESET)"; exit 1; }
	$(COMPOSE) exec redpanda rpk topic describe $(TOPIC)

.PHONY: sh
sh: ## Open shell in a service container. SERVICE=<n> required
	@test -n "$(SERVICE)" || { echo "$(RED)SERVICE is required$(RESET)"; exit 1; }
	$(COMPOSE) exec $(SERVICE) sh

## --- Health & debugging ---

.PHONY: smoke
smoke: ## Curl every service's healthz and readyz
	@for port in 8080 8081 8082 8083 8084; do \
		printf "port %s: " "$$port"; \
		curl -fsS "http://localhost:$$port/healthz" >/dev/null && \
			printf "$(GREEN)healthz$(RESET) " || printf "$(RED)healthz$(RESET) "; \
		curl -fsS "http://localhost:$$port/readyz" >/dev/null && \
			printf "$(GREEN)readyz$(RESET)\n" || printf "$(RED)readyz$(RESET)\n"; \
	done

.PHONY: grpc-list
grpc-list: ## List gRPC services exposed by a service. SERVICE=<n> required
	@test -n "$(SERVICE)" || { echo "$(RED)SERVICE is required$(RESET)"; exit 1; }
	@case "$(SERVICE)" in \
		gateway)    port=8080 ;; \
		identity)   port=8081 ;; \
		submission) port=8082 ;; \
		grading)    port=8083 ;; \
		vcs)        port=8084 ;; \
	esac; \
	grpcurl -plaintext localhost:$$port list

.PHONY: ui
ui: ## Open all dev UIs in browser
	@xdg-open http://localhost:16686 >/dev/null 2>&1 &
	@xdg-open http://localhost:9001 >/dev/null 2>&1 &
	@echo "Jaeger: http://localhost:16686"
	@echo "MinIO:  http://localhost:9001 (minio / minio12345)"

## --- Cleanup ---

.PHONY: clean
clean: ## Remove build artefacts
	find services -name bin -type d -exec rm -rf {} + 2>/dev/null || true
	find . -name coverage.out -delete 2>/dev/null || true
	find . -name coverage.html -delete 2>/dev/null || true

.PHONY: clean-all
clean-all: clean nuke ## clean + nuke volumes
	@echo "$(GREEN)full reset complete$(RESET)"