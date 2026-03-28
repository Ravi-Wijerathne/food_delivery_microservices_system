# Food Delivery Microservices System

A distributed food delivery backend system built using microservices architecture. This project demonstrates real-world patterns including event-driven communication, API Gateway, JWT authentication, Kubernetes orchestration, and more.

## Architecture Overview

```
                  [ Web UI ]
                       |
                 [ API Gateway ]
                       |
      -----------------------------------------
      |        |         |        |           |
    Auth     User   Restaurant   Order     Delivery
                                  |
                              Payment
                                  |
                          Notification

                 (RabbitMQ Event Bus)
```

## Tech Stack

- **Go (Golang)** - All microservices
- **gRPC** - Service-to-service communication
- **RabbitMQ** - Event-driven messaging
- **MongoDB** - Database per service
- **Docker** - Containerization
- **Kubernetes** - Orchestration
- **JWT** - Authentication

## Prerequisites

### Required Tools

| Tool | Version | Installation |
|------|---------|--------------|
| Go | 1.26+ | [Download](https://go.dev/dl/) |
| Docker | Latest | [Download](https://www.docker.com/) |
| kubectl | Latest | [Download](https://kubernetes.io/docs/tasks/tools/) |
| Kind | 0.20+ | `choco install kind` |
| protoc | 3.12+ | `choco install protoc` |

### Verify Installation

```bash
# Check all prerequisites
go version
docker --version
kubectl version --client
kind --version
protoc --version
```

## Project Structure

```
food_delivery_microservices_system/
├── services/
│   ├── auth-service/          # User authentication (Port 8081)
│   ├── restaurant-service/   # Restaurant & menu APIs (Port 8082)
│   ├── order-service/        # Order management (Port 8083 HTTP, 8084 gRPC)
│   ├── payment-service/      # Payment processing (Port 8085)
│   ├── delivery-service/     # Delivery assignment (Port 8086)
│   ├── notification-service/ # Notifications (Port 8087)
│   └── common/               # Shared utilities
├── gateway/                  # API Gateway (Port 8080)
├── proto/                   # Protocol Buffers definitions
├── deployments/              # Kubernetes manifests
├── web/                     # Web UI
├── docker-compose.yml       # Docker Compose configuration
├── start.ps1                # Automation script - starts all services
├── cleanup.ps1              # Automation script - stops all services
├── deploy-k8s.ps1           # Automation script - deploys to Kubernetes
└── README.md
```

## Automation Scripts

We provide PowerShell scripts to automate the entire startup process.

### Available Scripts

| Script | Description |
|--------|-------------|
| `start.ps1` | Full automation - checks prerequisites, starts infrastructure, builds and runs all services |
| `cleanup.ps1` | Stops all services and removes Docker containers |
| `deploy-k8s.ps1` | Deploys the entire system to Kubernetes (Kind cluster) |

### Usage

#### Option 1: Automated Startup (Recommended)
```powershell
# Run the full startup script
.\start.ps1
```

The script will:
1. Check all prerequisites (Go, Docker, kubectl)
2. Start MongoDB and RabbitMQ containers
3. Build all Go services
4. Start all services in the correct order
5. Verify all health endpoints

#### Option 2: Cleanup
```powershell
# Stop all services and remove containers
.\cleanup.ps1
```

#### Option 3: Kubernetes Deployment
```powershell
# Deploy to Kubernetes (requires Kind)
.\deploy-k8s.ps1
```

#### Script Options

```powershell
# Skip building (use existing binaries)
.\start.ps1 -SkipBuild

# Use Docker Compose instead
.\start.ps1 -UseDockerCompose

# Deploy to Kubernetes
.\deploy-k8s.ps1

# Delete Kubernetes cluster
.\deploy-k8s.ps1 -DeleteCluster
```

---

## Quick Start (Local Development)

### Option 1: Run Services Manually

#### 1. Start MongoDB
```bash
docker run -d --name mongodb -p 27017:27017 mongo:latest
```

#### 2. Start RabbitMQ
```bash
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management
```

#### 3. Build All Services
```bash
# Auth Service
cd services/auth-service
go build -o auth-service .
./auth-service &

# Restaurant Service
cd services/restaurant-service
go build -o restaurant-service .
./restaurant-service &

# Order Service
cd services/order-service
go build -o order-service .
./order-service &

# Payment Service
cd services/payment-service
go build -o payment-service .
./payment-service &

# Delivery Service
cd services/delivery-service
go build -o delivery-service .
./delivery-service &

# Notification Service
cd services/notification-service
go build -o notification-service .
./notification-service &

# API Gateway
cd gateway
go build -o gateway .
./gateway &
```

#### 4. Verify Services
```bash
# Check health endpoints
curl http://localhost:8081/health  # Auth
curl http://localhost:8082/health  # Restaurant
curl http://localhost:8083/health  # Order
curl http://localhost:8085/health  # Payment
curl http://localhost:8086/health  # Delivery
curl http://localhost:8087/health  # Notification
curl http://localhost:8080/api/health  # Gateway
```

### Option 2: Use Docker Compose

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down
```

## Testing the Application

### 1. Register a User
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123","name":"John Doe"}'
```

### 2. Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### 3. Get Restaurants
```bash
curl http://localhost:8080/api/restaurants
```

### 4. Get Menu
```bash
curl http://localhost:8080/api/menu/rest-1
```

### 5. Place Order (Requires JWT Token)
```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <YOUR_TOKEN>" \
  -d '{
    "user_id": "user-1",
    "restaurant_id": "rest-1",
    "items": [
      {
        "menu_item_id": "item-1",
        "name": "Margherita Pizza",
        "quantity": 2,
        "price": 12.99
      }
    ],
    "delivery_address": "123 Main St"
  }'
```

### 6. Check Order Status
```bash
curl http://localhost:8080/api/orders/<ORDER_ID> \
  -H "Authorization: Bearer <YOUR_TOKEN>"
```

## Event Flow

```
Order Created → Payment Processed → Delivery Assigned → Notifications
```

Check the logs to see the event-driven flow:
```bash
# Order Service
tail -f services/order-service/order.log

# Payment Service
tail -f services/payment-service/payment.log

# Delivery Service
tail -f services/delivery-service/delivery.log

# Notification Service
tail -f services/notification-service/notification.log
```

## Kubernetes Deployment

### Prerequisites
- Docker running
- Kind installed: `kind create cluster`
- kubectl configured

### Deploy to Kubernetes

```bash
# Create Kind cluster
kind create cluster --name food-delivery

# Build and load images
kind load docker-image auth-service:latest
kind load docker-image restaurant-service:latest
kind load docker-image order-service:latest
kind load docker-image payment-service:latest
kind load docker-image delivery-service:latest
kind load docker-image notification-service:latest
kind load docker-image gateway:latest

# Deploy infrastructure
kubectl apply -f deployments/mongodb.yaml
kubectl apply -f deployments/rabbitmq.yaml

# Deploy services
kubectl apply -f deployments/auth-service.yaml
kubectl apply -f deployments/restaurant-service.yaml
kubectl apply -f deployments/order-service.yaml
kubectl apply -f deployments/payment-service.yaml
kubectl apply -f deployments/delivery-service.yaml
kubectl apply -f deployments/notification-service.yaml
kubectl apply -f deployments/gateway.yaml

# Check status
kubectl get pods
kubectl get services

# View logs
kubectl logs -l app=auth-service
```

### Verify Kubernetes Deployment

```bash
# Check all pods are running
kubectl get pods -o wide

# Check services
kubectl get svc

# Test gateway
kubectl port-forward svc/gateway 8080:80
curl http://localhost:8080/api/health
```

### Access RabbitMQ Management UI
```bash
kubectl port-forward svc/rabbitmq 15672:15672
# Open http://localhost:15672 (guest/guest)
```

### Access MongoDB
```bash
kubectl port-forward svc/mongodb-service 27017:27017
# Connect with MongoDB Compass
```

## Web UI

Open `web/index.html` in a browser to use the simple web interface:

1. Register/Login
2. Browse Restaurants
3. View Menu
4. Place Orders

## Troubleshooting

### Port Already in Use
```bash
# Find process using port
netstat -ano | findstr :8080

# Kill process
taskkill /PID <PID> /F
```

### Docker Issues
```bash
# Restart Docker
docker restart mongodb
docker restart rabbitmq

# Check logs
docker logs mongodb
docker logs rabbitmq
```

### MongoDB Connection Issues
```bash
# Verify MongoDB is running
docker ps | grep mongo

# Test connection
docker exec -it mongodb mongosh
```

### RabbitMQ Issues
```bash
# Check RabbitMQ status
docker exec rabbitmq rabbitmqctl status

# View management UI logs
docker logs rabbitmq
```

## API Reference

### Auth Service
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/register` | POST | Register new user |
| `/login` | POST | Login and get JWT |
| `/health` | GET | Health check |

### Restaurant Service
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/restaurants` | GET | List all restaurants |
| `/menu/{id}` | GET | Get restaurant menu |
| `/health` | GET | Health check |

### Order Service
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/orders` | POST | Create new order |
| `/orders/{id}` | GET | Get order by ID |
| `/health` | GET | Health check |

### Payment Service
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Health check |

### Delivery Service
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Health check |

### Notification Service
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Health check |

## Environment Variables

| Service | Variable | Default |
|---------|----------|---------|
| Gateway | PORT | 8080 |
| Auth | MONGO_URI | mongodb://localhost:27017 |
| Order | MONGO_URI | mongodb://localhost:27017 |
| All | RABBITMQ_URI | amqp://guest:guest@localhost:5672/ |

## License

MIT License
