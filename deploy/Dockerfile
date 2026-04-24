ARG GO_VERSION=1.26

FROM golang:${GO_VERSION}-alpine AS builder
RUN apk add --no-cache git ca-certificates

WORKDIR /src
COPY go.work go.work.sum* ./
COPY gen ./gen
COPY pkg ./pkg
COPY services ./services
COPY test ./test

ARG SERVICE
ENV CGO_ENABLED=0
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    cd services/${SERVICE} && \
    go build -trimpath -ldflags="-s -w" -o /out/svc ./cmd

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
ARG SERVICE

COPY --from=builder /src/services/${SERVICE}/migrations /app/migrations
COPY --from=builder /out/svc /app/svc

ENTRYPOINT ["/app/svc"]