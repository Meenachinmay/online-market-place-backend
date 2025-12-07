# syntax=docker/dockerfile:1
# ─────────────────────────────────────────────
# Stage 1: Builder
FROM golang:1.24-alpine AS builder
LABEL stage=builder

RUN apk add --no-cache git

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /src
COPY go.mod go.sum ./

COPY vendor ./vendor

COPY . .

# Build the application
RUN go build -mod=vendor -ldflags="-w -s" -o /out/server ./app/services/soda-interview-grpc

# Download grpc-health-probe
RUN GRPC_HEALTH_PROBE_VERSION=v0.4.19 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

# ─────────────────────────────────────────────
# Stage 2: Runner
FROM alpine:3.19
WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata
RUN apk add --no-cache docker-cli

# Create non-root user
RUN addgroup -g 1000 -S appuser && \
    adduser -u 1000 -S appuser -G appuser

# Copy binaries and files
COPY --from=builder /out/server ./server
COPY --from=builder /bin/grpc_health_probe /bin/grpc_health_probe
COPY --from=builder /src/foundation/config ./foundation/config
COPY --from=builder /src/business/data/schema/migrations ./business/data/schema/migrations

# Change ownership
RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 50055

CMD ["./server"]
