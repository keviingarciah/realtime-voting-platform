# Architecture

## Overview

This project implements a **real-time voting pipeline** for a Guatemalan music band contest. Votes enter through two parallel producer stacks (gRPC and Rust), are queued in Kafka, consumed by a Go service, and persisted to Redis (live rankings) and MongoDB (audit logs). A Vue.js web app queries MongoDB for historical vote data.

## Components

### Producer Layer

| Service | Language | Protocol | Role |
|---|---|---|---|
| `producer/grpc` | Go | gRPC + HTTP | Receives vote payloads, publishes to Kafka topic `so1-proyecto2` |
| `producer/rust` | Rust | HTTP | Same vote ingestion path, different runtime for performance comparison |

Each producer deployment runs a **client + server** pair in the same pod:
- **Client** exposes the HTTP endpoint that Locust hits via Ingress.
- **Server** handles the internal gRPC/HTTP processing and Kafka publishing.

### Messaging Layer

- **Kafka** deployed via the [Strimzi](https://strimzi.io/) operator on Kubernetes.
- Topic: `so1-proyecto2`
- Bootstrap broker: `my-cluster-kafka-bootstrap:9092` (in-cluster DNS)

### Consumer Layer

| Service | Language | Replicas | Role |
|---|---|---|---|
| `consumer` | Go | 2–5 (HPA) | Reads Kafka messages, writes to Redis and MongoDB |

The consumer uses `segmentio/kafka-go` with a consumer group per pod instance. HPA scales replicas when CPU exceeds 40%.

### Data Layer

| Store | Purpose | Access Pattern |
|---|---|---|
| **Redis** | Live vote counts / rankings | Read by Grafana for real-time dashboards |
| **MongoDB** | Vote logs with timestamps | Queried by the Go REST API for the web UI |

### Observability

- **Grafana 8.4** with Redis datasource plugin for live vote visualization.
- Deployed as a LoadBalancer service for demo access.

### Web Application

| Component | Stack | Role |
|---|---|---|
| `app/server` | Go (Fiber) | REST API querying MongoDB |
| `app/client` | Vue.js 2 | Frontend console for browsing vote logs |

### Load Testing

- **Locust** reads vote fixtures from `locust/traffic.json` and POSTs randomized payloads to Ingress paths `/grpc` and `/rust`.

## Data Flow

```
Locust → Ingress (/grpc|/rust) → Producer Client → Producer Server → Kafka
  → Consumer → Redis (live) + MongoDB (logs) → Grafana / Web App
```

## Network Topology (Kubernetes)

```
Namespace: so1-proyecto2
├── grpc-deployment     (client:3000, server:3001)
├── rust-deployment     (client:8000, server:8080)
├── consumer-deployment (2-5 replicas, HPA)
├── redis-deployment    (6379)
├── mongodb-deployment  (27017)
├── grafana             (3000)
└── Ingress             (/grpc → grpc-service:3000, /rust → rust-service:8000)
```

## Design Tradeoffs

### Why two producers?

Demonstrates protocol comparison under identical load. In production you would typically standardize on one protocol.

### Why Strimzi instead of a managed Kafka service?

- Shows ability to run stateful workloads on Kubernetes.
- No cloud vendor lock-in for the messaging layer.
- Tradeoff: more operational complexity vs. Confluent Cloud or GCP Pub/Sub.

### Why HPA only on the consumer?

Producers handle lightweight HTTP ingestion; the consumer does database writes and is the scaling bottleneck during vote spikes.

### Why LoadBalancer on internal services?

Simplifies lab/demo debugging. In production:
- Use `ClusterIP` for Redis, MongoDB, Kafka.
- Expose only Ingress (and optionally Grafana via internal VPN).

## Security Considerations

> **Portfolio note:** Some services currently use hardcoded connection strings. For a production deployment, migrate to:
> - Kubernetes Secrets or [External Secrets Operator](https://external-secrets.io/)
> - ConfigMaps for non-sensitive configuration
> - NetworkPolicies to restrict pod-to-pod communication

## Future Improvements

- [ ] Terraform modules for GKE cluster provisioning
- [ ] GitHub Actions pipeline to build and push Docker images on merge
- [ ] Prometheus + Alertmanager for metrics beyond Grafana
- [ ] Kustomize overlays for dev/staging/prod environments
- [ ] Sealed Secrets for credential management
