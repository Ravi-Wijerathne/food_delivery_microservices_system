# Food Delivery Microservices System

A production-grade distributed food delivery platform built using Go (Golang) microservices. The system showcases real-world distributed architectures: **Event-Driven Asynchronous Messaging via RabbitMQ**, **Synchronous gRPC Validation**, **Database-per-Service Isolation on MongoDB**, an **API Gateway with Circuit Breaking and Prometheus Metrics**, **Kubernetes (Kind) Cluster Deployments**, and a **Modern Glassmorphic Web GUI (CravePulse)**.

---

## 🏛️ Architecture Overview

```text
                  ┌────────────────────────────────────────────────────────┐
                  │          CravePulse Web GUI (Port 3333)                │
                  │   [ Seamless Target Switcher: Docker / K8s / Local ]   │
                  └───────────────┬────────────────────────┬───────────────┘
                                  │                        │
                     (Port 8080)  │                        │ (Port 8089)
                                  ▼                        ▼
                    ┌─────────────────────────┐  ┌─────────────────────────┐
                    │      Docker Compose     │  │   Kubernetes Cluster    │
                    │       API Gateway       │  │       API Gateway       │
                    │  (Circuit Breaker+CORS) │  │   (LoadBalancer Service)│
                    └────────────┬────────────┘  └────────────┬────────────┘
                                 │                            │
        ┌────────────────────────┼────────────────────────────┴────────────────────────┐
        ▼                        ▼                            ▼                        ▼
 ┌──────────────┐         ┌──────────────┐             ┌──────────────┐         ┌──────────────┐
 │ Auth Service │         │ User Service │             │  Restaurant  │         │ Order Service│
 │  (Port 8081) │         │ (8088/g9088) │             │  (8082/g9082)│         │ (8083/g8084) │
 └──────────────┘         └──────────────┘             └──────────────┘         └──────┬───────┘
                                                                                       │
 ┌─────────────────────────────────────────────────────────────────────────────────────┘
 │  RabbitMQ Event Bus (Port 5672 / Management 15672)
 ├──► OrderCreated ──────────► [ Payment Service ] (Port 8085)
 ├──► PaymentProcessed ──────► [ Delivery Service ] (Port 8086)
 └──► DeliveryAssigned ──────► [ Notification Service ] (Port 8087)
                                 │
                                 ▼
                   ┌───────────────────────────────┐
                   │        MongoDB Cluster        │
                   │  (auth-db, user-db, order-db, │
                   │   restaurant-db, payment-db,  │
                   │         delivery-db)          │
                   └───────────────────────────────┘
```

---

## 🛠️ Tech Stack & Microservices Registry

| Service | Port(s) | Protocol | Role & Functionality |
| :--- | :--- | :--- | :--- |
| **API Gateway** | `8080` (HTTP) | REST / CORS | Reverse proxy, Circuit Breaker (`gobreaker`), JWT auth, Prometheus metrics |
| **Auth Service** | `8081` (HTTP) | REST | Password hashing (`bcrypt`), JWT token generation & verification |
| **User Service** | `8088` (HTTP), `9088` (gRPC) | REST & gRPC | Customer profiles, addresses, synchronous gRPC identity verification |
| **Restaurant Service** | `8082` (HTTP), `9082` (gRPC) | REST & gRPC | Catalogs, menus, synchronous gRPC menu & price validation |
| **Order Service** | `8083` (HTTP), `8084` (gRPC) | REST, gRPC & AMQP | Order orchestration, gRPC client, RabbitMQ publisher |
| **Payment Service** | `8085` (HTTP) | AMQP & REST | Consumes `OrderCreated`, processes payment, publishes `PaymentProcessed` |
| **Delivery Service** | `8086` (HTTP) | AMQP & REST | Consumes `PaymentProcessed`, assigns couriers, publishes `DeliveryAssigned` |
| **Notification Service** | `8087` (HTTP) | AMQP & REST | Subscribes to lifecycle events, dispatches notifications |
| **RabbitMQ** | `5672` (AMQP), `15672` (Web UI)| AMQP 0-9-1 | Distributed message broker for async choreography |
| **MongoDB** | `27017` | BSON / WiredTiger | Database-per-service pattern (isolated namespaces) |
| **CravePulse Web UI** | `3333` | HTML / CSS / JS | Modern dark glassmorphic SPA with live event bus tracking |

---

## 🚀 Quick Launch: Docker Compose

Docker Compose builds all Go binaries from source and boots the full 10-container ecosystem.

### 1. Build and Start the Entire System

```bash
# Build and start all services, databases, and message brokers in the background
docker compose up --build -d

# Check the status of all containers (all should be 'Up')
docker compose ps
```

### 2. Verify Health Endpoints

Each microservice exposes a dedicated health check endpoint:

```bash
curl http://localhost:8080/api/health  # API Gateway
curl http://localhost:8081/health      # Auth Service
curl http://localhost:8082/health      # Restaurant Service
curl http://localhost:8083/health      # Order Service
curl http://localhost:8085/health      # Payment Service
curl http://localhost:8086/health      # Delivery Service
curl http://localhost:8087/health      # Notification Service
curl http://localhost:8088/health      # User Service
```

### 3. Monitor Live Logs

Stream live async events across the event-driven workers:

```bash
docker compose logs -f payment-service delivery-service notification-service
```

---

## ☸️ Launching on Kubernetes (Kind Cluster)

The project includes production-ready Kubernetes manifests in the `deployments/` directory with full service discovery and gRPC inter-pod communication.

### 1. Create the Kind Cluster

```bash
# Install kind if not already installed (macOS Homebrew)
brew install kind

# Create cluster named 'food-delivery'
kind create cluster --name food-delivery

# Verify cluster node is Ready
kubectl get nodes
```

### 2. Tag and Load Images into the Cluster

Build the local images with Docker and import them directly into Kind's containerd image store:

```bash
# Build Docker images
docker compose build

# Tag images for Kubernetes manifests
docker tag food_delivery_microservices_system-auth-service:latest auth-service:latest
docker tag food_delivery_microservices_system-delivery-service:latest delivery-service:latest
docker tag food_delivery_microservices_system-gateway:latest gateway:latest
docker tag food_delivery_microservices_system-notification-service:latest notification-service:latest
docker tag food_delivery_microservices_system-order-service:latest order-service:latest
docker tag food_delivery_microservices_system-payment-service:latest payment-service:latest
docker tag food_delivery_microservices_system-restaurant-service:latest restaurant-service:latest
docker tag food_delivery_microservices_system-user-service:latest user-service:latest

# Load images into the Kind cluster
for img in auth-service:latest delivery-service:latest gateway:latest notification-service:latest order-service:latest payment-service:latest restaurant-service:latest user-service:latest mongo:latest rabbitmq:3-management; do
  echo "Loading $img into cluster..."
  docker save "$img" | docker exec -i food-delivery-control-plane ctr -n k8s.io images import -
done
```

### 3. Deploy Workloads to Kubernetes

```bash
# 1. Deploy Infrastructure (MongoDB & RabbitMQ)
kubectl apply -f deployments/mongodb.yaml
kubectl apply -f deployments/rabbitmq.yaml

# 2. Deploy All Microservices & Gateway
kubectl apply -f deployments/

# 3. Verify all pods are in 'Running' status (1/1)
kubectl get pods
kubectl get svc
```

### 4. Access the Kubernetes API Gateway

Port-forward the Kubernetes Gateway service to a local port (e.g. `8089`):

```bash
kubectl port-forward svc/gateway 8089:80
```

Now you can send requests to Kubernetes at `http://localhost:8089/api/health`.

---

## 🎨 CravePulse Web GUI (Frontend)

The application includes a luxury, responsive glassmorphic Web UI located in `web/`.

### 1. Launch the Web UI

```bash
# Option A: Start a lightweight local server (Port 3333)
python3 -m http.server 3333 --directory web

# Option B: Open directly in your browser
open web/index.html   # On macOS
start web/index.html  # On Windows
```

Visit: **[http://localhost:3333](http://localhost:3333)**

### 2. Frontend Features

- **Dynamic Backend Target Switcher**: The status pill in the top header auto-detects your backend. Click the badge to toggle seamlessly between:
  - 🟢 **Docker Gateway** (`http://localhost:8080/api`)
  - 🔵 **Kubernetes Gateway** (`http://localhost:8089/api`)
  - 🟡 **Simulated Local Event Bus** (In-browser testing with zero backend required)
- **Restaurant Discovery**: Filter by cuisine (Italian, American, Japanese, Mexican, Indian) with real-time dish search and curated culinary photography.
- **Interactive Cart & Checkout**: Slide-over drawer with item quantity steppers, live tax, delivery fee calculation, and delivery address editor.
- **Live Order Tracker & Dispatch Radar**:
  - Real-time visual timeline tracking the microservice pipeline:
    `Order Created` ➔ `Payment Processed` ➔ `Delivery Assigned` ➔ `Out for Delivery` ➔ `Delivered`
  - Animated SVG vector radar showing courier movement from restaurant to destination.
  - Assigned driver contact card with avatar, rating, and vehicle model.
- **Architecture & Observability Explorer**: View live status cards for all 7 microservices, gRPC ports, protocols, and a real-time RabbitMQ event stream terminal.

---

## 🧪 End-to-End API Testing (cURL Workflow)

You can run the complete automated lifecycle test via `curl` against either Docker (`:8080`) or Kubernetes (`:8089`):

```bash
# Set Gateway Port (8080 for Docker, 8089 for Kubernetes)
GATEWAY_URL="http://localhost:8080/api"

# 1. Register a New Customer
REGISTER_RES=$(curl -s -X POST "$GATEWAY_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"email":"customer@example.com","password":"securepassword123","name":"Alex Vance"}')
echo "$REGISTER_RES"

# 2. Authenticate & Obtain JWT Token
LOGIN_RES=$(curl -s -X POST "$GATEWAY_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"customer@example.com","password":"securepassword123"}')
TOKEN=$(echo "$LOGIN_RES" | grep -o '"token":"[^"]*' | cut -d'"' -f4)
echo "JWT Token: Bearer ${TOKEN:0:30}..."

# 3. Fetch User Profile
curl -s "$GATEWAY_URL/users/profile" \
  -H "Authorization: Bearer $TOKEN"

# 4. Discover Restaurants (MongoDB)
curl -s "$GATEWAY_URL/restaurants"

# 5. Fetch Restaurant Menu
curl -s "$GATEWAY_URL/menu/rest-1"

# 6. Place Order (gRPC Menu Validation + RabbitMQ Event Publish)
ORDER_RES=$(curl -s -X POST "$GATEWAY_URL/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "user_id": "user-1",
    "restaurant_id": "rest-1",
    "items": [
      {
        "menu_item_id": "item-1",
        "name": "Margherita Pizza",
        "quantity": 2,
        "price": 12.99
      },
      {
        "menu_item_id": "item-2",
        "name": "Pepperoni Pizza",
        "quantity": 1,
        "price": 14.99
      }
    ],
    "delivery_address": "742 Evergreen Terrace, Apt 4B"
  }')
echo "Order Created: $ORDER_RES"
ORDER_ID=$(echo "$ORDER_RES" | grep -o '"id":"[^"]*' | cut -d'"' -f4)

# 7. Check Order Status (Watch status transition: CREATED -> DELIVERED)
sleep 2
curl -s "$GATEWAY_URL/orders/$ORDER_ID" \
  -H "Authorization: Bearer $TOKEN"
```

---

## 📊 Observability & Metrics

- **Circuit Breaker**: The API Gateway uses `gobreaker` to prevent cascading failures if any downstream service is unavailable.
- **Prometheus Metrics**: Request count and status code counters are exposed at:
  ```bash
  curl http://localhost:8080/metrics
  ```
- **RabbitMQ Management Dashboard**:
  Access the RabbitMQ UI at **[http://localhost:15672](http://localhost:15672)** (Username: `guest`, Password: `guest`) to inspect queues:
  - `OrderCreated`
  - `PaymentProcessed`
  - `DeliveryAssigned`

---

## 🧹 Tear Down & Cleanup

### Stop Docker Compose

```bash
# Stop and remove containers and networks
docker compose down

# Stop and purge volumes (resets MongoDB database)
docker compose down -v
```

### Delete Kind Kubernetes Cluster

```bash
kind delete cluster --name food-delivery
```

---

## 📄 License

MIT License © 2026 Food Delivery Microservices System
