#!/bin/bash
set -e

# Configuration
DOCKER_REGISTRY="your-registry.io"  # Replace with your Docker registry
IMAGE_NAME="go-authentication"
IMAGE_TAG="latest"
NAMESPACE="auth-system"

# Database credentials - DO NOT COMMIT THESE VALUES
# Set these variables in your environment or enter them when prompted
DB_USER=${DB_USER:-}
DB_PASSWORD=${DB_PASSWORD:-}
JWT_SECRET=${JWT_SECRET:-}

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Building and deploying Go Authentication Service${NC}"

# Prompt for credentials if not provided
if [ -z "$DB_USER" ]; then
  read -p "Enter database username: " DB_USER
fi

if [ -z "$DB_PASSWORD" ]; then
  read -s -p "Enter database password: " DB_PASSWORD
  echo
fi

if [ -z "$JWT_SECRET" ]; then
  read -s -p "Enter JWT secret key: " JWT_SECRET
  echo
fi

# Validate inputs
if [ -z "$DB_USER" ] || [ -z "$DB_PASSWORD" ] || [ -z "$JWT_SECRET" ]; then
  echo -e "${RED}Error: Database credentials and JWT secret are required${NC}"
  exit 1
fi

# Step 1: Build Docker image
echo -e "${GREEN}Building Docker image...${NC}"
cd ..
docker build -t ${DOCKER_REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG} .

# Step 2: Push Docker image
echo -e "${GREEN}Pushing Docker image to registry...${NC}"
docker push ${DOCKER_REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}

# Step 3: Create namespace if it doesn't exist
echo -e "${GREEN}Creating namespace if it doesn't exist...${NC}"
kubectl apply -f namespace.yaml

# Step 4: Create the secret with actual credentials
echo -e "${GREEN}Creating Kubernetes secret with credentials...${NC}"
kubectl create secret generic go-auth-secrets \
  --from-literal=DB_USER="$DB_USER" \
  --from-literal=DB_PASSWORD="$DB_PASSWORD" \
  --from-literal=JWT_SECRET="$JWT_SECRET" \
  -n ${NAMESPACE} \
  --dry-run=client -o yaml | kubectl apply -f -

# Step 5: Update kustomization.yaml with the correct image
echo -e "${GREEN}Updating kustomization.yaml...${NC}"
sed -i "s|#images:|images:|g" kustomization.yaml
sed -i "s|#- name: \${DOCKER_REGISTRY}/go-authentication|- name: \${DOCKER_REGISTRY}/go-authentication|g" kustomization.yaml
sed -i "s|#  newName: your-registry.io/go-authentication|  newName: ${DOCKER_REGISTRY}/${IMAGE_NAME}|g" kustomization.yaml
sed -i "s|#  newTag: latest|  newTag: ${IMAGE_TAG}|g" kustomization.yaml

# Step 6: Deploy with kustomize
echo -e "${GREEN}Deploying with kustomize...${NC}"
kubectl apply -k .

# Step 7: Wait for deployment to be ready
echo -e "${GREEN}Waiting for deployment to be ready...${NC}"
kubectl rollout status deployment/go-authentication -n ${NAMESPACE}

# Step 8: Get service information
echo -e "${GREEN}Getting service information...${NC}"
echo "Go Authentication Service: $(kubectl get svc go-authentication -n ${NAMESPACE} -o jsonpath='{.metadata.name}')"

# Step 9: Get Kong Ingress information
echo -e "${GREEN}Getting Kong Ingress information...${NC}"
KONG_IP=$(kubectl get svc kong-kong-proxy -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
echo "Kong Proxy IP: ${KONG_IP}"
echo "Go Authentication API is available at: http://${KONG_IP}/auth"

echo -e "${YELLOW}Deployment completed successfully!${NC}" 