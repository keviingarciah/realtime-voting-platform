#!/usr/bin/env bash
set -euo pipefail

NAMESPACE="so1-proyecto2"
K8S_DIR="$(cd "$(dirname "$0")/../k8s" && pwd)"

echo "==> Deploying to namespace: ${NAMESPACE}"

MANIFESTS=(
  "namespace.yaml"
  "redis.yaml"
  "mongo.yaml"
  "grpc.yaml"
  "rust.yaml"
  "consumer.yaml"
  "grafana.yaml"
  "hpa.yaml"
  "ingress.yaml"
)

for manifest in "${MANIFESTS[@]}"; do
  echo "==> Applying ${manifest}"
  kubectl apply -f "${K8S_DIR}/${manifest}"
done

echo ""
echo "==> Deployment complete. Checking status..."
kubectl get all -n "${NAMESPACE}"
echo ""
echo "Next steps:"
echo "  1. Install Kafka (Strimzi) — see docs/en/deployment.md"
echo "  2. Update k8s/ingress.yaml with your external IP"
echo "  3. Run load tests: cd locust && locust -f main.py --host http://<INGRESS_HOST>"
