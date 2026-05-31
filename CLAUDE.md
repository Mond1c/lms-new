# CLAUDE.md — Agent guide for the LMS project

Operating manual for working in this repo. Read it before making changes.

## What this is

**LMS** — a learning-management platform for **programming courses where
students submit work through their own Git repositories** and get **automatically
graded**. Think *GitHub Classroom + an autograder*, self-hostable across multiple
VCS providers (Gitea, GitHub, GitLab).

The end-to-end flow it is built around:

1. An **instructor** creates a **Course** (bound to a VCS org) and **Assignments**
   (template repo, deadlines, max score).
2. A **student** enrolls, links a **VCS identity** via OAuth, and gets a
   **student repo** provisioned from the assignment template.
3. The student pushes code / opens a PR → the VCS provider fires a **webhook** →
   the **vcs** service normalizes it → an event goes onto **Kafka**.
4. The **submission** service turns that event into a **Submission**; CI runs the
   grader; the **grading** service ingests the **result** (per-test-case scores,
   feedback) and emits a result event.
5. The **gateway** (BFF) serves the web frontend: login, dashboard, submission
   detail.

Status: **early development.** The proto contracts and infra wiring are
substantially defined; **identity** is the only service with real business logic
— but it now implements its full RPC surface (User, Course, Enrollment,
Assignment, VCSIdentity, StudentRepo). The other services are app/handler
skeletons returning `Unimplemented`. There is no frontend yet.

## ⚠️ Module layout — read first

This is a **Go workspace** (`go.work`), not a single module. Module path prefix is
**`github.com/Mond1c/lms`**, Go **1.26**. Each of these is its own module:

```
gen/go            services/gateway    services/identity   services/vcs
pkg               services/grading    services/submission test/e2e
```

Practical consequences:

- There is **no root `go.mod`** — `go <cmd> ./...` from the repo root won't work.
  Run Go tooling **inside a module dir** (`cd services/identity && go test ./...`)
  or, preferably, use the **Makefile** (below), which loops over modules for you.
- Cross-module deps use `replace` directives in each service's `go.mod`
  (e.g. `replace github.com/Mond1c/lms/pkg => ../../pkg`).

## How to build / test / run — use the Makefile

`make help` lists everything. The ones you'll use most:

```bash
make build              # build all services (SERVICE=identity for one)
make test               # go test -race -count=1 across pkg + all services
make test SERVICE=identity
make lint               # golangci-lint
make fmt                # gofmt -s -w (skips gen/)
make tidy               # go mod tidy in every module

make proto              # regenerate Go from proto/  (cd proto && buf generate)
make proto-lint         # buf lint
make sqlc               # regenerate sqlc code for services that have sqlc.yaml
make generate           # proto + sqlc

make up                 # full docker-compose stack
make up-infra           # just postgres + redpanda + minio + jaeger
make smoke              # curl every service's /healthz + /readyz
make test-e2e           # bring up stack + run test/e2e
make logs SERVICE=...   make psql DB=lms_identity   make down   make nuke
```

## Repository layout

```
proto/lms/v1/      # protobuf — THE source of truth for every API (package lms.v1)
  common, identity, vcs, submission, grading, gateway, events .proto
  buf.yaml, buf.gen.yaml
gen/go/            # buf output (protobuf + ConnectRPC stubs). NEVER edit by hand.
  lms/v1/*.pb.go, lms/v1/lmsv1connect/*.connect.go
services/<svc>/    # one module per service, consistent internal layout:
  cmd/main.go              # entrypoint: load config, app.New, app.Run
  internal/config/         # env-based config (caarlos0/env)
  internal/app/            # wiring: obs, db pool, mux, ConnectRPC handler, health
  internal/handler/        # ConnectRPC handlers (proto in/out)  -> service
  internal/repo/           # persistence: sqlc queries + mapping  (sqlcgen/, queries/)
  internal/domain/         # domain types / value objects (Email, Password, Role…)
  migrations/              # golang-migrate SQL (NNN_name.up.sql / .down.sql)
pkg/               # shared libs: jwt, grpcauth, kafka, outbox, obs, pg
deploy/            # docker-compose.yml, Dockerfile (multi-service), init-db.sql
test/e2e/          # cross-service smoke/e2e tests
```

Note: identity's business layer lives in `internal/service/` (the old top-level
`service/` package was moved there in Phase 1). Keep new services under
`internal/`.

## Architecture & stack

- **Transport:** [ConnectRPC](https://connectrpc.com) over HTTP/2 (`h2c`). Each
  service mounts its generated handler plus `/healthz` and `/readyz`. The
  **gateway** is the edge/BFF that fronts the browser and calls the other services
  as Connect clients.
- **Contracts:** Protobuf via **buf**; generated into `gen/go` with
  `go_package_prefix = github.com/Mond1c/lms/gen/go`, `paths=source_relative`.
- **Persistence:** PostgreSQL via **pgx/v5**, type-safe queries via **sqlc**
  (write SQL in `internal/repo/queries/*.sql`, run `make sqlc`). Migrations via
  **golang-migrate**. **One database per service** (`lms_identity`,
  `lms_submission`, `lms_grading`, `lms_vcs`); the **gateway is stateless**.
  Services never touch each other's DB — only RPC.
- **Eventing:** **Kafka** (Redpanda in dev) via **franz-go**, with a
  **transactional outbox** (`pkg/outbox`) so events publish atomically with DB
  writes. Event schemas live in `proto/lms/v1/events.proto`.
- **Observability:** `pkg/obs` sets up `slog` + OpenTelemetry tracing (OTLP →
  Jaeger). All services take `OTLP_ENDPOINT`, `LOG_LEVEL`, `LOG_FORMAT`.
- **Auth:** JWT HS256 (`pkg/jwt`), validated by a Connect interceptor
  (`pkg/grpcauth`). Passwords: bcrypt. IDs: ULID.
- **Object storage:** MinIO (S3) for submission artifacts.
- **Config:** environment variables, parsed with `caarlos0/env` (see each
  `internal/config/config.go`). DB URL is `required`; everything else has
  defaults.
- **Tests:** `testify`; repo-layer tests spin up real Postgres via
  `ory/dockertest` (see `internal/repo/testmain_test.go`).

### Ports

| Service | Port | | Infra | Port |
|---|---|---|---|---|
| gateway | 8080 | | postgres | 5432 |
| identity | 8081 | | redpanda (kafka) | 9092 |
| submission | 8082 | | minio | 9000 / console 9001 |
| grading | 8083 | | jaeger UI | 16686 |
| vcs | 8084 | | otlp | 4317 / 4318 |

## Conventions (match existing code)

- **Proto-first.** Change behavior by editing `.proto`, then `make proto`. Never
  hand-edit `gen/`. Keep `lms.v1` package; shared types go in `common.proto`,
  VCS-shared types (`ProviderRef`, `Repo`, enums) in `vcs.proto`.
- **Layering:** `handler` (proto ⇄ domain, maps errors to Connect codes) →
  `service` (business rules, defines the repo interface it needs) → `repo` (sqlc).
  Dependencies point inward; `service` imports neither Connect nor SQL.
- **Interfaces at the consumer.** e.g. `UserRepo` is declared in the `service`
  package that uses it, not the `repo` package.
- **Errors:** sentinel errors per layer (`repo.ErrNotFound`,
  `domain.ErrInvalidEmail`); wrap with `fmt.Errorf("...: %w", err)`; map to
  Connect codes only in the handler (`toConnectErr`).
- **SQL:** all queries in `queries/*.sql`, regenerated via sqlc — don't write SQL
  in Go. Migrations are append-only; never edit a committed migration.
- **Formatting:** `make fmt` (gofmt -s). Lint clean with `make lint`.
- **Commits:** `[scope] type: summary` — e.g. `[identity] feat: add users repo`.
  Scope = a service name, `proto`, `pkg`, or `build`. The history even tags
  intentional defects, e.g. `... bug: NewEmail returns original value`.
- **Git:** branch off `main`, keep it green.

## Build/test status (verified)

`identity` builds, is lint-clean, and its full RPC surface is implemented with
passing `service` unit tests and `repo` integration tests (testcontainers spins
up Postgres for the repo layer). The historical `NewEmail returns original value`
bug from commit `85b6cdf` has **since been fixed** — `NewEmail` now
trims+lowercases and returns the normalized value.

## Known issues / rough edges

- **`services/vcs/internal/domain/provider.go`** is a `// Temp` stub: type
  `NormilizedEvent` (note the misspelling vs. proto `NormalizedEvent`) and
  `ProviderKind = int32` alias. The real VCS domain isn't built yet.
- **Skeleton services.** vcs / submission / grading / gateway have app + config +
  a handler embedding `Unimplemented…Handler` — no business logic yet. identity is
  the only service with logic, and it now implements its whole RPC surface.
- identity Course persists a VCS binding but has no setter RPC yet (no input on
  `CreateCourse`); `GetUser` doesn't embed VCS identities (use `ListVCSIdentities`).

## Where to look first

- Architecture + diagrams: **[`docs/architecture.md`](docs/architecture.md)**.
- Roadmap, status, open questions: **[`docs/plan.md`](docs/plan.md)**.
- The domain & all RPCs: **`proto/lms/v1/*.proto`** (start with `identity.proto`,
  `submission.proto`, `grading.proto`, `vcs.proto`, then `gateway.proto`).
- A complete vertical slice to imitate: **`services/identity`**
  (cmd → app → handler → service → repo → domain).
- Shared infrastructure patterns: **`pkg/`** and **`deploy/docker-compose.yml`**.
- A frontend brief to feed a design tool: **`.claude/design-prompt.md`**.
