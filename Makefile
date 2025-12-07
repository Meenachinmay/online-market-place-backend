#----------------
# Helper variables
#----------------
CLUSTER_NAME := nominomi-local
NAMESPACE := nominomi
KUBE_CONTEXT := kind-$(CLUSTER_NAME)
SERVICE_NAME := soda-service
IMAGE_NAME := soda-interview-grpc

#----------------
# General
#----------------
all: proto sqlc build

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*##"; printf "\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  %%-15s %%s\n", $$1, $$2 }' $(MAKEFILE_LIST)

#----------------
# Code Generation & Building
#----------------
proto: ## Generate protobuf files
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative foundation/proto/product/v1/product.proto foundation/proto/referral-blog/v1/referral_blog.proto foundation/proto/order/v1/order.proto foundation/proto/soda-finance/v1/finance.proto

sqlc: ## Generate database code
	sqlc generate

build: ## Build the Go binary
	go build ./app/services/soda-interview-grpc/...

run: ## Run the service locally
	go run ./app/services/soda-interview-grpc/main.go

tidy: ## Tidy go modules
	go mod tidy
	go mod vendor

#----------------
# Testing
#----------------
test: ## Run all tests
	@echo "▶ Running tests..."
	go test ./...

test-verbose: ## Run tests with verbose output
	@echo "▶ Running tests (verbose)..."
	go test -v ./...

test-coverage: ## Run tests with coverage report
	@echo "▶ Running tests with coverage..."
	go test -v -cover ./...
	@echo ""
	@echo "▶ Generating coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "▶ Coverage report generated: coverage.html"

#----------------
# Database (Docker Compose)
#----------------
db-start: ## Start database using Docker Compose
	docker compose -p $(SERVICE_NAME) -f zarf/docker/docker-compose.local.yaml up -d --remove-orphans

db-stop: ## Stop database
	@echo "▶ Stopping Docker Compose..."
	@docker compose -p $(SERVICE_NAME) -f zarf/docker/docker-compose.local.yaml down

#----------------
# Kubernetes (Kind)
#----------------
k8s-namespace-ensure: ## Ensure Kubernetes namespace exists
	@kubectl get ns $(NAMESPACE) >/dev/null 2>&1 || kubectl create namespace $(NAMESPACE)

k8s-apply-secrets: ## Apply local secrets manifest
	@kubectl apply -f k8s/secrets.yaml

k8s-build-nocache: ## Build Docker image without cache and load into Kind
	@echo "▶ Building Docker image without cache..."
	$(eval TIMESTAMP := $(shell date +%s))
	@docker build --no-cache --progress=plain -f zarf/docker/k8s.Dockerfile \
		-t $(IMAGE_NAME):latest \
		-t $(IMAGE_NAME):$(TIMESTAMP) .
	@echo "▶ Loading image into Kind cluster $(CLUSTER_NAME)..."
	@kind load docker-image $(IMAGE_NAME):latest --name $(CLUSTER_NAME)
	@echo "✅ Built and loaded image with tag: latest and $(TIMESTAMP)"

k8s-deploy: ## Apply Kubernetes manifests
	@echo "▶ Applying Kubernetes manifests..."
	@kubectl config use-context $(KUBE_CONTEXT)
	@kubectl apply -f k8s/

k8s-rollout: ## Restart deployment
	@echo "▶ Rolling out restart for $(SERVICE_NAME)..."
	@kubectl rollout restart deployment/$(SERVICE_NAME) -n $(NAMESPACE)
	@kubectl rollout status deployment/$(SERVICE_NAME) -n $(NAMESPACE) --timeout=120s || true

k8s-redeploy: ## Re-build and re-deploy
	$(MAKE) k8s-build-nocache
	$(MAKE) k8s-rollout

k8s-start: ## Start everything in K8s (Cluster must exist)
	@kubectl config use-context $(KUBE_CONTEXT)
	$(MAKE) k8s-namespace-ensure
	$(MAKE) k8s-apply-secrets
	$(MAKE) k8s-build-nocache
	$(MAKE) k8s-deploy
	$(MAKE) k8s-rollout
	@echo "▶ Port-forwarding service/$(SERVICE_NAME) 50055 -> localhost:50055"
	@kubectl port-forward -n $(NAMESPACE) service/$(SERVICE_NAME) 50055:50055 >/dev/null 2>&1 &

k8s-stop: ## Stop port-forward
	@echo "▶ Stopping port-forward (if any)..."
	@pkill -f "kubectl port-forward -n $(NAMESPACE) service/$(SERVICE_NAME) 50055:50055" || true

logs: ## View logs of the service
	@kubectl logs -n $(NAMESPACE) -l app=$(SERVICE_NAME) -f

k8s-status: ## Check status of all resources in namespace
	@echo "Cluster Info:"
	@kubectl cluster-info --context $(KUBE_CONTEXT)
	@echo "\nPods in $(NAMESPACE):"
	@kubectl get pods -n $(NAMESPACE)
	@echo "\nServices in $(NAMESPACE):"
	@kubectl get services -n $(NAMESPACE)

port-forward: ## Manual port forward
	kubectl port-forward service/$(SERVICE_NAME) 50055:50055 -n $(NAMESPACE)
