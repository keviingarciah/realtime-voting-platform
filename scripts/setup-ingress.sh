#!/usr/bin/env bash
set -euo pipefail

echo "==> Installing NGINX Ingress Controller"

kubectl create namespace nginx-ingress --dry-run=client -o yaml | kubectl apply -f -

helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx 2>/dev/null || true
helm repo update

helm upgrade --install nginx-ingress ingress-nginx/ingress-nginx \
  --namespace nginx-ingress \
  --set controller.publishService.enabled=true

echo ""
echo "==> Waiting for Ingress Controller..."
kubectl wait --namespace nginx-ingress \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=120s

echo ""
echo "==> Ingress Controller external IP:"
kubectl get services -n nginx-ingress
