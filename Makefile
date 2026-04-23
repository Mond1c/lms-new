SERVICES := gateway identity submission grading vcs
COMPOSE := docker compose -f deploy/docker-compose.yml

.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make <target>\n\nTargets:\n"} \
		/^[a-zA-Z0-9_.-]+:.*##/ {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' \
		$(MAKEFILE_LIST)

.PHONY: proto
proto:
	cd proto && buf generate

.PHONY: proto-lint
proto-lint:
	cd proto && buf lint

.PHONY: proto-breaking
proto-breaking:
	cd proto && buf breaking --against '../.git#branch=main,subdir=proto'

# ---- Go ----

.PHONY: tidy
tidy:
	@for m in pkg gen/go $(SERVICES:%=services/%) test/e2e; do \
		echo "==> tidy $$m"; (cd $$m && go mod tidy) || exit 1; \
	done

.PHONY: build
build:
	@for s in $(SERVICES); do \
		echo "==> build $$s"; \
		(cd services/$$s && go build -o bin/svc ./cmd) || exit 1; \
	done

.PHONY: test
test:
	@for m in pkg $(SERVICES:%=services/%); do \
		echo "==> test $$m"; (cd $$m && go test -race -count=1 ./...) || exit 1; \
	done

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: up
up:
	$(COMPOSE) up -d --build

.PHONY: down
down:
	$(COMPOSE) down

.PHONY: nuke
nuke:
	$(COMPOSE) down -v

.PHONY: logs
logs:
	$(COMPOSE) logs -f --tail=100

.PHONY: ps
ps:
	$(COMPOSE) ps

.PHONY: smoke
smoke:
	@for port in 8080 8081 8082 8083 8084; do \
		echo "port $$port:"; \
		curl -fsS "http://localhost:$$port/healthz" && echo; \
		curl -fsS "http://localhost:$$port/readyz" && echo; \
	done

.PHONY: clean
clean:
	find services -name bin -type d -exec rm -rf {} +
