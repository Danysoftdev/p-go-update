#!/bin/bash

set -e  # Detener el script si ocurre un error

echo "📁 Desplegando microservicio p-go-update..."

# Namespace
kubectl apply -f k8s/update/namespace-update.yaml

# Secret
kubectl apply -f k8s/update/secrets-update.yaml

# Deployment
kubectl apply -f k8s/update/deployment-update.yaml

# Espera hasta que el Deployment esté listo
echo "⏳ Esperando a que p-go-update esté listo..."
kubectl wait --namespace=p-go-update \
  --for=condition=available deployment/update-deployment \
  --timeout=90s

# Service
kubectl apply -f k8s/update/service-update.yaml

# Ingress
kubectl apply -f k8s/update/ingress.yaml

echo "✅ p-go-update desplegado correctamente."

echo -e "\n🔍 Estado actual:"
kubectl get all -n p-go-update
kubectl get ingress -n p-go-update
