# Plataforma de Votación en Tiempo Real sobre Kubernetes

> Sistema de votación orientado a eventos para un concurso de bandas de música guatemalteca — diseñado para demostrar orquestación de contenedores, streaming de mensajes, autoescalado y observabilidad.

**Idioma:** [English](../../README.md) | Español

[![CI](https://github.com/keviingarciah/realtime-voting-platform/actions/workflows/ci.yml/badge.svg)](https://github.com/keviingarciah/realtime-voting-platform/actions/workflows/ci.yml)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-GKE-326CE5?logo=kubernetes&logoColor=white)](../../k8s/)
[![Kafka](https://img.shields.io/badge/Kafka-Strimzi-231F20?logo=apache-kafka&logoColor=white)](architecture.md)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](../../LICENSE)

---

## Inicio Rápido para Reclutadores

| | |
|---|---|
| **Qué hace** | Recibe tráfico de votación a través de productores gRPC y Rust, transmite eventos vía Kafka, persiste resultados en Redis/MongoDB y visualiza rankings en vivo en Grafana. |
| **Por qué importa** | Demuestra ownership DevOps de punta a punta: containerización, despliegue en K8s, enrutamiento con ingress, autoescalado HPA, pruebas de carga y monitoreo. |
| **Stack** | Go, Rust, Vue.js, Kafka (Strimzi), Redis, MongoDB, Grafana, GKE, NGINX Ingress, Locust |
| **Tiempo de deploy** | ~20 min en un cluster GKE con imágenes Docker pre-construidas |

---

## Arquitectura

```mermaid
flowchart LR
    subgraph Carga["Pruebas de Carga"]
        L[Locust]
    end

    subgraph Ingress["NGINX Ingress"]
        IN[/grpc · /rust]
    end

    subgraph Productores["Servicios Productores"]
        G[gRPC Client/Server]
        R[Rust Client/Server]
    end

    subgraph Mensajeria["Streaming de Eventos"]
        K[(Kafka / Strimzi)]
    end

    subgraph Consumidores["Capa Consumidora"]
        C[Go Consumer x2–5]
    end

    subgraph Almacenamiento["Capa de Datos"]
        RD[(Redis)]
        MG[(MongoDB)]
    end

    subgraph Observabilidad["Observabilidad"]
        GF[Grafana]
    end

    subgraph Web["Aplicación Web"]
        API[API Go]
        UI[Cliente Vue.js]
    end

    L --> IN
    IN --> G
    IN --> R
    G --> K
    R --> K
    K --> C
    C --> RD
    C --> MG
    RD --> GF
    MG --> API
    API --> UI
```

Consulta la [Guía de Arquitectura](architecture.md) para detalles de componentes y decisiones de diseño.

---

## Demo

### Pruebas de carga con Locust

![Configuración de Locust](../assets/screenshots/1.png)

![Resultados de tráfico Locust](../assets/screenshots/2.png)

### Dashboards en tiempo real y persistencia

![Dashboard Grafana](../assets/screenshots/3.png)

![Logs en MongoDB](../assets/screenshots/4.png)

---

## Destacados DevOps

- **Microservicios containerizados** — Cada componente incluye su `Dockerfile` y se publica en Docker Hub.
- **Manifiestos Kubernetes** — Despliegues declarativos con límites de recursos, servicios, ingress y HPA en [`k8s/`](../../k8s/).
- **Pipeline orientado a eventos** — Productores y consumidores desacoplados usando Kafka para resiliencia bajo carga.
- **Horizontal Pod Autoscaler** — Los pods del consumer escalan de 2 a 5 según utilización de CPU (objetivo 40%).
- **Enrutamiento Ingress** — Punto de entrada único exponiendo rutas `/grpc` y `/rust` vía NGINX Ingress Controller.
- **Observabilidad** — Dashboards Grafana conectados a Redis para visualización de votos en vivo.
- **Pruebas de carga** — Locust simula tráfico realista desde datos JSON de fixture.

---

## Inicio Rápido

### Prerrequisitos

- [Google Cloud SDK](https://cloud.google.com/sdk/docs/install) (`gcloud`)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [Helm 3](https://helm.sh/docs/intro/install/)
- Un cluster GKE (o cualquier cluster Kubernetes con soporte LoadBalancer)

### Desplegar

```bash
# 1. Clonar y configurar el contexto del cluster
git clone https://github.com/keviingarciah/realtime-voting-platform.git
cd realtime-voting-platform
gcloud container clusters get-credentials <nombre-cluster> --zone <zona>

# 2. Instalar NGINX Ingress Controller
./scripts/setup-ingress.sh

# 3. Desplegar todos los manifiestos de la aplicación
./scripts/deploy.sh

# 4. Instalar Kafka vía Strimzi (ver guía de despliegue)
```

Instrucciones completas: [Guía de Despliegue](deployment.md)

### Ejecutar pruebas de carga

```bash
cd locust
pip install locust
locust -f main.py --host http://<HOST_INGRESS>
```

Abre `http://localhost:8089`, configura usuarios/tasa de spawn y apunta el tráfico a `/grpc` o `/rust`.

---

## Estructura del Proyecto

```
realtime-voting-platform/
├── .github/workflows/   # Pipelines CI (validación de manifiestos, checks Go)
├── app/                 # Frontend Vue.js + API REST Go (logs MongoDB)
├── consumer/            # Consumer Kafka → Redis + MongoDB
├── producer/
│   ├── grpc/            # Productor de votos gRPC (Go)
│   └── rust/            # Productor de votos HTTP (Rust)
├── locust/              # Scripts de pruebas de carga
├── k8s/                 # Manifiestos Kubernetes (deployments, services, HPA, ingress)
├── scripts/             # Automatización de despliegue
└── docs/
    ├── en/              # Documentación en inglés
    ├── es/              # Documentación en español
    └── assets/          # Capturas y diagramas
```

---

## Decisiones Operacionales

| Decisión | Justificación |
|---|---|
| Kafka vía operador Strimzi | Kafka production-grade en K8s sin gestionar ZooKeeper manualmente |
| Doble productor (gRPC + Rust) | Comparar rendimiento de protocolos bajo patrones de carga idénticos |
| HPA solo en consumer | El consumer es el cuello de botella durante picos de votación |
| Redis para datos en vivo, MongoDB para logs | Separación de almacenamiento caliente (dashboard) vs. frío (auditoría) |
| Imágenes Docker pre-construidas | Bootstrap más rápido del cluster; CI puede reconstruir en cada merge |
| Servicios LoadBalancer | Simplifica acceso en lab/demo; en producción se usaría ClusterIP + Ingress |

---

## Documentación

| Documento | English | Español |
|---|---|---|
| README completo | [README.md](../../README.md) | Este archivo |
| Arquitectura | [docs/en/architecture.md](../en/architecture.md) | [architecture.md](architecture.md) |
| Despliegue | [docs/en/deployment.md](../en/deployment.md) | [deployment.md](deployment.md) |
| Cheatsheet K8s | [docs/en/kubernetes-cheatsheet.md](../en/kubernetes-cheatsheet.md) | [kubernetes-cheatsheet.md](kubernetes-cheatsheet.md) |

---

## Consejos para tu Portfolio

Este repositorio está estructurado como un **caso de estudio**, no como un volcado de código. Al presentarlo:

1. Comienza con el diagrama de arquitectura y explica el flujo de datos en una entrevista.
2. Destaca los **tradeoffs** (por qué Strimzi, por qué HPA en consumer, por qué doble productor).
3. Recorre un despliegue desde `kubectl apply` hasta ver Grafana actualizarse en vivo.
4. Menciona qué agregarías después: CI/CD para builds de imágenes, Terraform para GKE, sealed secrets, métricas Prometheus.

---

## Autor

**Kevin Garcia** — aspirante a DevOps Engineer

- GitHub: [@keviingarciah](https://github.com/keviingarciah)
- Docker Hub: [keviingarciah](https://hub.docker.com/u/keviingarciah)

---

## Licencia

Este proyecto está bajo la licencia MIT — ver [LICENSE](../../LICENSE) para más detalles.
