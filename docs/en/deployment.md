# Deployment Guide

Step-by-step instructions to deploy the voting platform on Google Kubernetes Engine (GKE).

## Prerequisites

Install and configure the following tools:

```bash
# Google Cloud SDK
curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-cli-linux-x86_64.tar.gz
tar -xf google-cloud-cli-linux-x86_64.tar.gz
./google-cloud-sdk/install.sh
gcloud init

# kubectl (via gcloud)
gcloud components install kubectl

# Helm 3
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
chmod 700 get_helm.sh
./get_helm.sh
```

## 1. Create or Connect to a GKE Cluster

```bash
# Option A: Use an existing cluster
gcloud container clusters get-credentials <cluster-name> --zone <zone> --project <project-id>

# Option B: Create a new cluster (example)
gcloud container clusters create voting-platform-cluster \
  --zone us-central1-a \
  --num-nodes 3 \
  --machine-type e2-medium \
  --project <project-id>
```

## 2. Install NGINX Ingress Controller

```bash
./scripts/setup-ingress.sh
```

Or manually:

```bash
kubectl create namespace nginx-ingress
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm install nginx-ingress ingress-nginx/ingress-nginx -n nginx-ingress
kubectl get services -n nginx-ingress
```

## 3. Install Kafka (Strimzi)

```bash
export NAMESPACE=so1-proyecto2

kubectl create namespace $NAMESPACE
kubectl create -f "https://strimzi.io/install/latest?namespace=$NAMESPACE" -n $NAMESPACE
kubectl apply -f https://strimzi.io/examples/latest/kafka/kafka-persistent-single.yaml -n $NAMESPACE

# Wait for Kafka to be ready
kubectl wait kafka/my-cluster --for=condition=Ready --timeout=300s -n $NAMESPACE
```

## 4. Deploy Application Manifests

```bash
./scripts/deploy.sh
```

This applies manifests in the correct order:

1. Namespace
2. Redis, MongoDB
3. gRPC and Rust producers
4. Consumer
5. Grafana
6. HPA
7. Ingress

## 5. Configure Ingress Host

Update `k8s/ingress.yaml` with your cluster's external IP:

```bash
# Get the Ingress Controller external IP
kubectl get svc -n nginx-ingress

# Edit ingress host (use nip.io for quick DNS)
# Example: 34.72.55.202.nip.io
kubectl apply -f k8s/ingress.yaml
```

## 6. Verify Deployment

```bash
kubectl get all -n so1-proyecto2
kubectl get ingress -n so1-proyecto2
kubectl get hpa -n so1-proyecto2
```

## 7. Run Load Tests

```bash
cd locust
pip install locust

# Point at your Ingress host
locust -f main.py --host http://<INGRESS_IP>.nip.io
```

Open `http://localhost:8089`:
- Set **Number of users** and **Spawn rate**
- Click **Start swarming**
- Monitor Grafana and consumer pod scaling

## 8. Access Services

| Service | How to Access |
|---|---|
| gRPC producer | `http://<INGRESS_HOST>/grpc` |
| Rust producer | `http://<INGRESS_HOST>/rust` |
| Grafana | `kubectl get svc grafana -n so1-proyecto2` → external IP:3000 |
| MongoDB | Via web app or direct connection to LoadBalancer IP |

## Teardown

```bash
kubectl delete namespace so1-proyecto2
kubectl delete namespace nginx-ingress
helm uninstall nginx-ingress -n nginx-ingress
```

## Troubleshooting

| Symptom | Check |
|---|---|
| Pods in `CrashLoopBackOff` | `kubectl logs <pod> -n so1-proyecto2` |
| Consumer not receiving messages | Verify Kafka topic exists: `kubectl exec -it my-cluster-kafka-0 -n so1-proyecto2 -- bin/kafka-topics.sh --list --bootstrap-server localhost:9092` |
| HPA not scaling | Ensure metrics-server is installed: `kubectl get apiservice v1beta1.metrics.k8s.io` |
| Ingress 404 | Confirm ingress host matches external IP and paths `/grpc`, `/rust` are correct |

## Docker Images

Pre-built images are available on Docker Hub:

| Image | Tag |
|---|---|
| `keviingarciah/so1_grpc_client` | v1.4 |
| `keviingarciah/so1_grpc_server` | v1.4 |
| `keviingarciah/so1_rust_client` | v1.1 |
| `keviingarciah/so1_rust_server` | v1.1 |
| `keviingarciah/so1_consumer` | v1.6 |

To rebuild locally:

```bash
# Example: consumer
cd consumer
docker build -t keviingarciah/so1_consumer:latest .
docker push keviingarciah/so1_consumer:latest
```
