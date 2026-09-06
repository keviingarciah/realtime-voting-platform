# Kubernetes Cheatsheet

Quick reference for common `kubectl` operations used in this project.

## Essential Commands

### Apply / Delete

```bash
kubectl apply -f <file>.yaml
kubectl delete -f <file>.yaml
```

### Get Resources

```bash
# Namespaces
kubectl get namespaces

# Deployments
kubectl get deployments
kubectl get deployments --all-namespaces

# Services
kubectl get services
kubectl get services --all-namespaces

# Pods
kubectl get pods
kubectl get pods --all-namespaces
kubectl get pods -n so1-proyecto2 -o wide

# Ingress
kubectl get ingress -n so1-proyecto2

# HPA
kubectl get hpa -n so1-proyecto2
```

### Set Default Namespace

```bash
kubectl config set-context --current --namespace=so1-proyecto2
```

### Logs and Debugging

```bash
kubectl logs <pod-name> -n so1-proyecto2
kubectl logs <pod-name> -n so1-proyecto2 -f          # follow
kubectl describe pod <pod-name> -n so1-proyecto2
kubectl exec -it <pod-name> -n so1-proyecto2 -- /bin/sh
```

## Kafka (Strimzi)

```bash
# Install operator
kubectl create -f 'https://strimzi.io/install/latest?namespace=so1-proyecto2' -n so1-proyecto2

# Deploy Kafka cluster
kubectl apply -f https://strimzi.io/examples/latest/kafka/kafka-persistent-single.yaml -n so1-proyecto2

# Check status
kubectl get kafka -n so1-proyecto2
kubectl wait kafka/my-cluster --for=condition=Ready --timeout=300s -n so1-proyecto2
```

## Ingress (NGINX)

```bash
kubectl create namespace nginx-ingress
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm install nginx-ingress ingress-nginx/ingress-nginx -n nginx-ingress
kubectl get services -n nginx-ingress
```

## Project-Specific

```bash
# Deploy everything
./scripts/deploy.sh

# Check all resources in project namespace
kubectl get all -n so1-proyecto2

# Watch HPA scaling during load test
kubectl get hpa -n so1-proyecto2 -w
```
