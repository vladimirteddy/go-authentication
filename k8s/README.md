# Kubernetes Deployment Guide

This directory contains Kubernetes manifests for deploying the Go Authentication Service with Traefik integration.

## Components

- `deployment.yaml`: Defines the authentication service deployment with security contexts and resource limits
- `service.yaml`: Defines the ClusterIP service with Traefik-specific annotations
- `ingress.yaml`: Configures Traefik IngressRoute with TLS and middleware
- `configmap.yaml`: Contains non-sensitive configuration
- `secret.yaml`: Contains sensitive data (credentials, JWT secret)
- `namespace.yaml`: Defines the service namespace
- `consul-defaults.yaml`: Consul service defaults
- `consul-intentions.yaml`: Consul service-to-service communication rules

## Prerequisites

1. Kubernetes cluster with Traefik installed
2. Consul running in the cluster
3. TLS certificate secret (referenced in ingress.yaml)

## Deployment

1. Create the namespace:

   ```bash
   kubectl apply -f namespace.yaml
   ```

2. Create secrets and configmap:

   ```bash
   kubectl apply -f secret.yaml
   kubectl apply -f configmap.yaml
   ```

3. Deploy Consul configurations:

   ```bash
   kubectl apply -f consul-defaults.yaml
   kubectl apply -f consul-intentions.yaml
   ```

4. Deploy the service:

   ```bash
   kubectl apply -f service.yaml
   kubectl apply -f deployment.yaml
   ```

5. Configure Traefik routing:
   ```bash
   kubectl apply -f ingress.yaml
   ```

## Security Features

- Non-root container execution
- Read-only root filesystem
- Dropped capabilities
- Rate limiting middleware
- CORS protection
- TLS encryption

## Traefik Integration

The service is configured to work with Traefik through:

1. IngressRoute with:

   - TLS termination
   - Path-based routing
   - Rate limiting middleware
   - CORS headers middleware
   - Path stripping

2. Service annotations for:

   - Load balancing method (DRR)
   - Session affinity
   - Service weights

3. Deployment annotations for:
   - Traefik middleware association
   - TLS configuration

## Health Checks

The service includes:

- Readiness probe: `/health` endpoint (5s initial delay, 10s period)
- Liveness probe: `/health` endpoint (15s initial delay, 20s period)

## Monitoring

The service exposes metrics for Prometheus at `/metrics` endpoint.

## Scaling

The deployment is configured with:

- 2 replicas by default
- Resource requests and limits
- Horizontal Pod Autoscaling (HPA) support

## Troubleshooting

1. Check pod status:

   ```bash
   kubectl get pods -l app=go-authentication
   ```

2. View pod logs:

   ```bash
   kubectl logs -l app=go-authentication
   ```

3. Check Traefik routing:

   ```bash
   kubectl get ingressroute
   kubectl get middleware
   ```

4. Verify Consul service registration:
   ```bash
   kubectl exec -it consul-server-0 -- consul catalog services
   ```

```

```
