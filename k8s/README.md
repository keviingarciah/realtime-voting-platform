# Kubernetes Manifests

Declarative manifests for the `so1-proyecto2` namespace.

## Apply Order

Manifests must be applied in this order (handled automatically by `../scripts/deploy.sh`):

| Order | File | Resources |
|---|---|---|
| 1 | `namespace.yaml` | Namespace `so1-proyecto2` |
| 2 | `redis.yaml` | Redis Deployment + LoadBalancer Service |
| 3 | `mongo.yaml` | MongoDB Deployment + LoadBalancer Service |
| 4 | `grpc.yaml` | gRPC producer (client + server) |
| 5 | `rust.yaml` | Rust producer (client + server) |
| 6 | `consumer.yaml` | Kafka consumer (2 replicas) |
| 7 | `grafana.yaml` | Grafana with Redis datasource ConfigMap |
| 8 | `hpa.yaml` | HPA for consumer (2–5 replicas, 40% CPU) |
| 9 | `ingress.yaml` | NGINX Ingress routing `/grpc` and `/rust` |

## Prerequisites

- Kafka must be deployed separately via [Strimzi](https://strimzi.io/) before the consumer can process messages.
- NGINX Ingress Controller must be installed (`../scripts/setup-ingress.sh`).

## Configuration

Before deploying, update:

- **`ingress.yaml`** — Replace the `host` field with your cluster's external IP (e.g., `34.x.x.x.nip.io`).
- **`grafana.yaml`** — Update the Redis datasource URL and password in the ConfigMap.

## Quick Apply

```bash
# From repo root
make deploy
# or
./scripts/deploy.sh
```

## Documentation

- [Deployment Guide (EN)](../docs/en/deployment.md)
- [Guía de Despliegue (ES)](../docs/es/deployment.md)
