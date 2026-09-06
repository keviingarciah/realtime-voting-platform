.PHONY: deploy ingress validate test-locust help

NAMESPACE := so1-proyecto2

help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

deploy: ## Deploy all Kubernetes manifests
	./scripts/deploy.sh

ingress: ## Install NGINX Ingress Controller
	./scripts/setup-ingress.sh

validate: ## Validate Kubernetes manifests locally
	@for f in k8s/*.yaml; do echo "==> $$f"; kubectl apply --dry-run=client -f $$f; done

status: ## Show deployment status
	kubectl get all -n $(NAMESPACE)
	kubectl get ingress -n $(NAMESPACE)
	kubectl get hpa -n $(NAMESPACE)

logs-consumer: ## Tail consumer pod logs
	kubectl logs -l app=consumer -n $(NAMESPACE) -f --tail=50

test-locust: ## Run Locust load test (set HOST=http://your-ingress)
	locust -f locust/main.py --host $(HOST)

build-consumer: ## Build consumer Docker image
	docker build -t keviingarciah/so1_consumer:latest consumer/

build-grpc: ## Build gRPC Docker images
	docker build -t keviingarciah/so1_grpc_server:latest producer/grpc/server/
	docker build -t keviingarciah/so1_grpc_client:latest producer/grpc/client/

build-rust: ## Build Rust Docker images
	docker build -t keviingarciah/so1_rust_server:latest producer/rust/server/
	docker build -t keviingarciah/so1_rust_client:latest producer/rust/client/
