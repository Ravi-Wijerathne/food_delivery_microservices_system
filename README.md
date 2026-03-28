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
└── README.md
```

## First-Time Setup (Recommended)

If this is your first time running the project, use this checklist before any other section.

### 1. Clone and Enter Project

```bash
git clone https://github.com/Ravi-Wijerathne/food_delivery_microservices_system.git
cd food_delivery_microservices_system
```

### 2. Build Local Docker Images (Required for Kind)

```bash
docker build -t auth-service:latest ./services/auth-service
docker build -t restaurant-service:latest ./services/restaurant-service
docker build -t order-service:latest -f services/order-service/Dockerfile .
docker build -t payment-service:latest -f services/payment-service/Dockerfile ./services
docker build -t delivery-service:latest -f services/delivery-service/Dockerfile ./services
docker build -t notification-service:latest -f services/notification-service/Dockerfile ./services
docker build -t gateway:latest ./gateway
```

### 3. Create Kind Cluster and Load Images

```bash
kind create cluster --name food-delivery

kind load docker-image auth-service:latest --name food-delivery
kind load docker-image restaurant-service:latest --name food-delivery
kind load docker-image order-service:latest --name food-delivery
kind load docker-image payment-service:latest --name food-delivery
kind load docker-image delivery-service:latest --name food-delivery
kind load docker-image notification-service:latest --name food-delivery
kind load docker-image gateway:latest --name food-delivery
```

### 4. Deploy Everything

```bash
kubectl apply -f deployments
kubectl get pods
kubectl get svc
```

### 5. Quick Functional Check

```bash
kubectl port-forward svc/gateway 8080:80
curl http://localhost:8080/api/health
```

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

Note: On some machines, `rabbitmq` may take a few seconds longer to become ready.
If `payment-service`, `delivery-service`, or `notification-service` exits on first startup, run:

```bash
docker-compose up -d
```

again and re-check with:

```bash
docker-compose ps --all
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

# Build images first (skip if already built)
docker build -t auth-service:latest ./services/auth-service
docker build -t restaurant-service:latest ./services/restaurant-service
docker build -t order-service:latest -f services/order-service/Dockerfile .
docker build -t payment-service:latest -f services/payment-service/Dockerfile ./services
docker build -t delivery-service:latest -f services/delivery-service/Dockerfile ./services
docker build -t notification-service:latest -f services/notification-service/Dockerfile ./services
docker build -t gateway:latest ./gateway

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

# Wait for rollouts (recommended)
kubectl rollout status deployment/mongodb --timeout=180s
kubectl rollout status deployment/rabbitmq --timeout=180s
kubectl rollout status deployment/auth-service --timeout=180s
kubectl rollout status deployment/restaurant-service --timeout=180s
kubectl rollout status deployment/order-service --timeout=180s
kubectl rollout status deployment/payment-service --timeout=180s
kubectl rollout status deployment/delivery-service --timeout=180s
kubectl rollout status deployment/notification-service --timeout=180s
kubectl rollout status deployment/gateway --timeout=180s

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
kubectl port-forward svc/rabbitmq-service 15672:15672
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

### Reset Kind Cluster (If API Server Is Unreachable)
```bash
kind delete cluster --name food-delivery
kind create cluster --name food-delivery
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
