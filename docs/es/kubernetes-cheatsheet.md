# Cheatsheet de Kubernetes

Referencia rápida de operaciones comunes de `kubectl` usadas en este proyecto.

## Comandos Esenciales

### Crear / Eliminar

```bash
kubectl apply -f <archivo>.yaml
kubectl delete -f <archivo>.yaml
```

### Obtener Recursos

```bash
# Namespaces
kubectl get namespaces

# Deployments
kubectl get deployments
kubectl get deployments --all-namespaces

# Servicios
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

### Configurar Namespace Predeterminado

```bash
kubectl config set-context --current --namespace=so1-proyecto2
```

### Logs y Debugging

```bash
kubectl logs <nombre-pod> -n so1-proyecto2
kubectl logs <nombre-pod> -n so1-proyecto2 -f          # seguir en vivo
kubectl describe pod <nombre-pod> -n so1-proyecto2
kubectl exec -it <nombre-pod> -n so1-proyecto2 -- /bin/sh
```

## Kafka (Strimzi)

```bash
# Instalar operador
kubectl create -f 'https://strimzi.io/install/latest?namespace=so1-proyecto2' -n so1-proyecto2

# Desplegar cluster Kafka
kubectl apply -f https://strimzi.io/examples/latest/kafka/kafka-persistent-single.yaml -n so1-proyecto2

# Verificar estado
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

## Específico del Proyecto

```bash
# Desplegar todo
./scripts/deploy.sh

# Ver todos los recursos en el namespace del proyecto
kubectl get all -n so1-proyecto2

# Observar escalado HPA durante prueba de carga
kubectl get hpa -n so1-proyecto2 -w
```
