# Guía de Despliegue

Instrucciones paso a paso para desplegar la plataforma de votación en Google Kubernetes Engine (GKE).

## Prerrequisitos

Instala y configura las siguientes herramientas:

```bash
# Google Cloud SDK
curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-cli-linux-x86_64.tar.gz
tar -xf google-cloud-cli-linux-x86_64.tar.gz
./google-cloud-sdk/install.sh
gcloud init

# kubectl (vía gcloud)
gcloud components install kubectl

# Helm 3
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
chmod 700 get_helm.sh
./get_helm.sh
```

## 1. Crear o Conectar a un Cluster GKE

```bash
# Opción A: Usar un cluster existente
gcloud container clusters get-credentials <nombre-cluster> --zone <zona> --project <project-id>

# Opción B: Crear un cluster nuevo (ejemplo)
gcloud container clusters create voting-platform-cluster \
  --zone us-central1-a \
  --num-nodes 3 \
  --machine-type e2-medium \
  --project <project-id>
```

## 2. Instalar NGINX Ingress Controller

```bash
./scripts/setup-ingress.sh
```

O manualmente:

```bash
kubectl create namespace nginx-ingress
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm install nginx-ingress ingress-nginx/ingress-nginx -n nginx-ingress
kubectl get services -n nginx-ingress
```

## 3. Instalar Kafka (Strimzi)

```bash
export NAMESPACE=so1-proyecto2

kubectl create namespace $NAMESPACE
kubectl create -f "https://strimzi.io/install/latest?namespace=$NAMESPACE" -n $NAMESPACE
kubectl apply -f https://strimzi.io/examples/latest/kafka/kafka-persistent-single.yaml -n $NAMESPACE

# Esperar a que Kafka esté listo
kubectl wait kafka/my-cluster --for=condition=Ready --timeout=300s -n $NAMESPACE
```

## 4. Desplegar Manifiestos de la Aplicación

```bash
./scripts/deploy.sh
```

Esto aplica los manifiestos en el orden correcto:

1. Namespace
2. Redis, MongoDB
3. Productores gRPC y Rust
4. Consumer
5. Grafana
6. HPA
7. Ingress

## 5. Configurar Host del Ingress

Actualiza `k8s/ingress.yaml` con la IP externa de tu cluster:

```bash
# Obtener la IP externa del Ingress Controller
kubectl get svc -n nginx-ingress

# Editar el host del ingress (usa nip.io para DNS rápido)
# Ejemplo: 34.72.55.202.nip.io
kubectl apply -f k8s/ingress.yaml
```

## 6. Verificar el Despliegue

```bash
kubectl get all -n so1-proyecto2
kubectl get ingress -n so1-proyecto2
kubectl get hpa -n so1-proyecto2
```

## 7. Ejecutar Pruebas de Carga

```bash
cd locust
pip install locust

# Apuntar al host del Ingress
locust -f main.py --host http://<IP_INGRESS>.nip.io
```

Abre `http://localhost:8089`:
- Configura **Number of users** y **Spawn rate**
- Haz clic en **Start swarming**
- Monitorea Grafana y el escalado de pods del consumer

## 8. Acceder a los Servicios

| Servicio | Cómo Acceder |
|---|---|
| Productor gRPC | `http://<HOST_INGRESS>/grpc` |
| Productor Rust | `http://<HOST_INGRESS>/rust` |
| Grafana | `kubectl get svc grafana -n so1-proyecto2` → IP externa:3000 |
| MongoDB | Vía app web o conexión directa a IP LoadBalancer |

## Eliminar Recursos

```bash
kubectl delete namespace so1-proyecto2
kubectl delete namespace nginx-ingress
helm uninstall nginx-ingress -n nginx-ingress
```

## Solución de Problemas

| Síntoma | Verificar |
|---|---|
| Pods en `CrashLoopBackOff` | `kubectl logs <pod> -n so1-proyecto2` |
| Consumer no recibe mensajes | Verificar que el topic Kafka existe: `kubectl exec -it my-cluster-kafka-0 -n so1-proyecto2 -- bin/kafka-topics.sh --list --bootstrap-server localhost:9092` |
| HPA no escala | Verificar que metrics-server está instalado: `kubectl get apiservice v1beta1.metrics.k8s.io` |
| Ingress 404 | Confirmar que el host del ingress coincide con la IP externa y las rutas `/grpc`, `/rust` son correctas |

## Imágenes Docker

Las imágenes pre-construidas están disponibles en Docker Hub:

| Imagen | Tag |
|---|---|
| `keviingarciah/so1_grpc_client` | v1.4 |
| `keviingarciah/so1_grpc_server` | v1.4 |
| `keviingarciah/so1_rust_client` | v1.1 |
| `keviingarciah/so1_rust_server` | v1.1 |
| `keviingarciah/so1_consumer` | v1.6 |

Para reconstruir localmente:

```bash
# Ejemplo: consumer
cd consumer
docker build -t keviingarciah/so1_consumer:latest .
docker push keviingarciah/so1_consumer:latest
```
