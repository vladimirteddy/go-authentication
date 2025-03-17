# Kubernetes Deployment for Go Authentication Service

This directory contains Kubernetes manifests for deploying the Go Authentication service with Consul service mesh and Kong API Gateway integration.

## Prerequisites

- Kubernetes cluster with Consul and Kong installed (follow the main implementation guide)
- Docker registry with the go-authentication image pushed
- Existing PostgreSQL database

## Security Note

This deployment uses a secure approach for handling sensitive information:

- Database credentials and JWT secrets are not stored in the repository
- Secrets are created at deployment time using the `build-and-deploy.sh` script
- You will be prompted to enter credentials during deployment

## Configuration

Before deploying, update the following files:

1. **configmap.yaml**: Update the `DB_HOST`, `DB_PORT`, and `DB_NAME` values to point to your existing PostgreSQL database
2. **deployment.yaml**: Update `${DOCKER_REGISTRY}` with your actual Docker registry URL
3. **build-and-deploy.sh**: Update the `DOCKER_REGISTRY` variable with your Docker registry URL

## Deployment

The easiest way to deploy is using the provided script:

```bash
# Make the script executable
chmod +x build-and-deploy.sh

# Run the deployment script
./build-and-deploy.sh
```

The script will:

1. Prompt for database credentials and JWT secret
2. Build and push the Docker image
3. Create the Kubernetes namespace
4. Create the Kubernetes secret with your credentials
5. Deploy all resources using kustomize
6. Provide the URL to access your service

## Manual Deployment

If you prefer to deploy manually, follow these steps:

1. Create the namespace

```bash
kubectl apply -f namespace.yaml
```

2. Create the secret with your actual credentials

```bash
kubectl create secret generic go-auth-secrets \
  --from-literal=DB_USER=your-db-user \
  --from-literal=DB_PASSWORD=your-db-password \
  --from-literal=JWT_SECRET=your-jwt-secret \
  -n auth-system
```

3. Deploy ConfigMap

```bash
kubectl apply -f configmap.yaml
```

4. Deploy Go Authentication service

```bash
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
```

5. Deploy Consul configuration

```bash
kubectl apply -f consul-defaults.yaml
kubectl apply -f consul-intentions.yaml
```

6. Deploy Kong Ingress

```bash
kubectl apply -f ingress.yaml
```

## Accessing the Service

Once deployed, the authentication service will be available at:

```
http://<KONG_PROXY_IP>/auth
```

## Monitoring

You can monitor the service through:

1. Kubernetes dashboard

```bash
kubectl get pods -n auth-system
kubectl logs -f deployment/go-authentication -n auth-system
```

2. Consul UI

```bash
# Get Consul UI address
kubectl get svc consul-ui
```

3. Kong Manager (if installed)

```bash
# Get Kong Manager address
kubectl get svc kong-kong-manager
```
