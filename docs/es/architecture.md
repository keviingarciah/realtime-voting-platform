# Arquitectura

## Visión General

Este proyecto implementa un **pipeline de votación en tiempo real** para un concurso de bandas de música guatemalteca. Los votos ingresan a través de dos stacks de productores en paralelo (gRPC y Rust), se encolan en Kafka, son consumidos por un servicio Go y se persisten en Redis (rankings en vivo) y MongoDB (logs de auditoría). Una aplicación web Vue.js consulta MongoDB para datos históricos.

## Componentes

### Capa de Productores

| Servicio | Lenguaje | Protocolo | Rol |
|---|---|---|---|
| `producer/grpc` | Go | gRPC + HTTP | Recibe payloads de votos, publica al topic Kafka `so1-proyecto2` |
| `producer/rust` | Rust | HTTP | Misma ruta de ingesta, runtime diferente para comparación de rendimiento |

Cada despliegue de productor ejecuta un par **cliente + servidor** en el mismo pod:
- **Cliente** expone el endpoint HTTP que Locust alcanza vía Ingress.
- **Servidor** maneja el procesamiento interno gRPC/HTTP y la publicación a Kafka.

### Capa de Mensajería

- **Kafka** desplegado vía el operador [Strimzi](https://strimzi.io/) en Kubernetes.
- Topic: `so1-proyecto2`
- Broker bootstrap: `my-cluster-kafka-bootstrap:9092` (DNS in-cluster)

### Capa de Consumidores

| Servicio | Lenguaje | Réplicas | Rol |
|---|---|---|---|
| `consumer` | Go | 2–5 (HPA) | Lee mensajes Kafka, escribe en Redis y MongoDB |

El consumer usa `segmentio/kafka-go` con un consumer group por instancia de pod. HPA escala réplicas cuando la CPU supera 40%.

### Capa de Datos

| Almacén | Propósito | Patrón de Acceso |
|---|---|---|
| **Redis** | Conteos de votos / rankings en vivo | Leído por Grafana para dashboards en tiempo real |
| **MongoDB** | Logs de votos con timestamps | Consultado por la API REST Go para la UI web |

### Observabilidad

- **Grafana 8.4** con plugin de datasource Redis para visualización de votos en vivo.
- Desplegado como servicio LoadBalancer para acceso en demos.

### Aplicación Web

| Componente | Stack | Rol |
|---|---|---|
| `app/server` | Go (Fiber) | API REST que consulta MongoDB |
| `app/client` | Vue.js 2 | Frontend para explorar logs de votos |

### Pruebas de Carga

- **Locust** lee fixtures de votos desde `locust/traffic.json` y envía payloads aleatorios via POST a las rutas Ingress `/grpc` y `/rust`.

## Flujo de Datos

```
Locust → Ingress (/grpc|/rust) → Cliente Productor → Servidor Productor → Kafka
  → Consumer → Redis (vivo) + MongoDB (logs) → Grafana / App Web
```

## Topología de Red (Kubernetes)

```
Namespace: so1-proyecto2
├── grpc-deployment     (cliente:3000, servidor:3001)
├── rust-deployment     (cliente:8000, servidor:8080)
├── consumer-deployment (2-5 réplicas, HPA)
├── redis-deployment    (6379)
├── mongodb-deployment  (27017)
├── grafana             (3000)
└── Ingress             (/grpc → grpc-service:3000, /rust → rust-service:8000)
```

## Tradeoffs de Diseño

### ¿Por qué dos productores?

Demuestra comparación de protocolos bajo carga idéntica. En producción se estandarizaría en un solo protocolo.

### ¿Por qué Strimzi en lugar de Kafka managed?

- Demuestra capacidad de ejecutar workloads stateful en Kubernetes.
- Sin vendor lock-in en la capa de mensajería.
- Tradeoff: más complejidad operacional vs. Confluent Cloud o GCP Pub/Sub.

### ¿Por qué HPA solo en el consumer?

Los productores manejan ingesta HTTP ligera; el consumer hace escrituras a base de datos y es el cuello de botella durante picos de votación.

### ¿Por qué LoadBalancer en servicios internos?

Simplifica debugging en lab/demo. En producción:
- Usar `ClusterIP` para Redis, MongoDB, Kafka.
- Exponer solo Ingress (y opcionalmente Grafana vía VPN interna).

## Consideraciones de Seguridad

> **Nota de portfolio:** Algunos servicios usan strings de conexión hardcodeados. Para producción, migrar a:
> - Kubernetes Secrets o [External Secrets Operator](https://external-secrets.io/)
> - ConfigMaps para configuración no sensible
> - NetworkPolicies para restringir comunicación pod-a-pod

## Mejoras Futuras

- [ ] Módulos Terraform para aprovisionamiento de cluster GKE
- [ ] Pipeline GitHub Actions para build y push de imágenes Docker en cada merge
- [ ] Prometheus + Alertmanager para métricas más allá de Grafana
- [ ] Overlays Kustomize para entornos dev/staging/prod
- [ ] Sealed Secrets para gestión de credenciales
