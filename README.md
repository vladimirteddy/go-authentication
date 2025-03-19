# Go Authentication Service

A robust authentication and authorization service built with Go, JWT, and RBAC, designed to work with Traefik API Gateway and Kubernetes.

## Features

- JWT-based authentication
- Role-based access control (RBAC)
- Database-stored permissions
- RESTful API
- Integration with Traefik API Gateway
- Kubernetes deployment ready

## Architecture

This service is designed as part of a microservices architecture:

- **Authentication**: Handles user login, signup, and token generation
- **Authorization**: Handles permission checks for users based on roles and resources
- **Traefik Integration**: Provides a Forward Auth endpoint for Traefik API Gateway

## Setup & Installation

### Prerequisites

- Go 1.16+
- PostgreSQL database
- Kubernetes cluster (optional)
- Traefik API Gateway (for authorization integration)

### Local Development

1. Clone the repository
2. Create a `.env` file with the following variables:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=authentication
SECRET_JWT=your_secret_key
PORT=8080
```

3. Run database migrations:

```bash
psql -U postgres -d authentication -f migrations/create_users_table.sql
psql -U postgres -d authentication -f migrations/create_roles_permissions_tables.sql
```

4. Start the server:

```bash
go run main.go
```

### Docker

Build the Docker image:

```bash
docker build -t go-authentication:latest .
```

Run the container:

```bash
docker run -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=your_password \
  -e DB_NAME=authentication \
  -e SECRET_JWT=your_secret_key \
  go-authentication:latest
```

## API Endpoints

### Authentication

- `POST /auth/signup` - Create a new user
- `POST /auth/login` - Login and get JWT token
- `GET /user/profile` - Get user profile (requires auth)

### Role Management

- `POST /roles` - Create a new role
- `GET /roles` - Get all roles
- `GET /roles/:id` - Get role by ID
- `PUT /roles/:id` - Update role
- `DELETE /roles/:id` - Delete role
- `POST /roles/assign` - Assign role to user
- `POST /roles/remove` - Remove role from user

### Permission Management

- `POST /permissions` - Create a new permission
- `GET /permissions` - Get all permissions
- `GET /permissions/:id` - Get permission by ID
- `GET /permissions/resource/:resource` - Get permissions by resource
- `PUT /permissions/:id` - Update permission
- `DELETE /permissions/:id` - Delete permission
- `POST /permissions/assign` - Assign permission to role
- `POST /permissions/remove` - Remove permission from role
- `POST /permissions/check` - Check if user has permission

### Traefik Integration

- `GET /traefik/auth` - Forward auth endpoint for Traefik

## Integrating with Traefik API Gateway

### Traefik Configuration

1. Configure Traefik to use the Forward Auth middleware:

```yaml
# Static configuration (traefik.yaml)
providers:
  file:
    directory: "/etc/traefik/dynamic"
    watch: true
```

2. Create a dynamic configuration file:

```yaml
# /etc/traefik/dynamic/auth.yaml
http:
  middlewares:
    auth-middleware:
      forwardAuth:
        address: "http://go-authentication:8080/traefik/auth"
        authResponseHeaders:
          - "X-User-ID"
          - "X-Username"
          - "X-User-Roles"
        trustForwardHeader: true

  routers:
    api:
      rule: "Host(`api.example.com`)"
      service: "api-service"
      middlewares:
        - "auth-middleware"
      tls: {}
```

### How it Works

1. When a request is made to a protected route:

   - Traefik forwards the request details to the auth endpoint
   - The auth service validates the JWT token
   - The auth service checks if the user has the required permissions
   - If authorized, the request proceeds to the target service
   - If unauthorized, a 401 or 403 response is returned

2. User information is passed to backend services:
   - `X-User-ID`: ID of the authenticated user
   - `X-Username`: Username of the authenticated user
   - `X-User-Roles`: Comma-separated list of user roles

## Kubernetes Deployment

See the `k8s` directory for Kubernetes deployment manifests.

## License

MIT
