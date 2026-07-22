# Food Delivery Microservices System

A distributed food delivery backend system built using a microservices architecture. This project demonstrates real-world patterns including event-driven communication, an API Gateway with Circuit Breaking, gRPC for inter-service communication, JWT authentication, and comprehensive observability (Prometheus/Grafana).

## Architecture Overview

```text
                  [ Web UI (Glassmorphism SPA) ]
                               |
                       [ API Gateway ]
        (Circuit Breaker + Prometheus Metrics + JWT Auth)
                               |
       -------------------------------------------------
       |         |           |           |             |
     Auth      User      Restaurant    Order        Delivery
                                         |
                                      Payment
                                         |
                                    Notification

                     (RabbitMQ Event Bus)
```

## Tech Stack

- **Go (Golang)** - All microservices (7 services + API Gateway)
- **gRPC & REST** - Service-to-service communication and external APIs
- **RabbitMQ** - Event-driven messaging for async workflows
- **MongoDB** - Database per service architecture
- **Docker & Docker Compose** - Containerization
- **Prometheus & Grafana** - Metrics and Observability
- **Circuit Breaker** - Implemented via `gobreaker` in the API Gateway
- **JWT** - Authentication
- **Vanilla JS + HTML + CSS** - Frontend UI with Glassmorphism and Live Polling

## Project Structure

```text
food_delivery_microservices_system/
├── services/
│   ├── auth-service/          # Authentication & user credentials (Port 8081)
│   ├── user-service/          # User profiles (Port 8088 HTTP, 9088 gRPC)
│   ├── restaurant-service/    # Restaurant & menus (Port 8082 HTTP, 9082 gRPC)
│   ├── order-service/         # Orders (Port 8083 HTTP, 8084 gRPC)
│   ├── payment-service/       # Payment processing (Port 8085)
│   ├── delivery-service/      # Delivery assignment (Port 8086)
│   ├── notification-service/  # Notifications (Port 8087)
│   └── common/                # Shared utilities
├── gateway/                   # API Gateway (Port 8080)
├── proto/                     # Protocol Buffers definitions (.proto)
├── deployments/               # Kubernetes manifests
├── web/                       # Web UI (Modern SPA)
├── docker-compose.yml         # Docker Compose configuration
└── README.md
```

## Quick Start (Local Development)

The entire microservices ecosystem runs via Docker Compose.

### 1. Clone and Enter Project

```bash
git clone https://github.com/Ravi-Wijerathne/food_delivery_microservices_system.git
cd food_delivery_microservices_system
```

### 2. Start the System

Run Docker Compose with the `--build` flag to build the Go binaries inside the containers:

```bash
# Start all services, databases, and message brokers in the background
docker-compose up --build -d

# View logs for all services
docker-compose logs -f

# Check the status of all containers
docker-compose ps
```

### 3. Verify Health

All services expose a health endpoint:
```bash
curl http://localhost:8080/api/health  # API Gateway
curl http://localhost:8081/health      # Auth
curl http://localhost:8082/health      # Restaurant
curl http://localhost:8083/health      # Order
curl http://localhost:8085/health      # Payment
curl http://localhost:8086/health      # Delivery
curl http://localhost:8087/health      # Notification
curl http://localhost:8088/health      # User
```

## Testing the Application

### Using the Web UI

Open the modern, interactive frontend directly in your browser:
```bash
open web/index.html  # On Mac
start web/index.html # On Windows
```
The UI allows you to register, log in, browse restaurants, place orders, and watch live order status updates via polling.

### Using the API (cURL)

#### 1. Register a User
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123","name":"John Doe"}'
```

#### 2. Login (Save the token!)
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

#### 3. View User Profile
```bash
curl http://localhost:8080/api/users/user-1 \
  -H "Authorization: Bearer <YOUR_TOKEN>"
```

#### 4. Get Restaurants & Menus
```bash
curl http://localhost:8080/api/restaurants
curl http://localhost:8080/api/menu/rest-1
```

#### 5. Place Order
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

#### 6. Check Order Status
```bash
curl http://localhost:8080/api/orders/<ORDER_ID> \
  -H "Authorization: Bearer <YOUR_TOKEN>"
```

## Event Flow

The system uses **RabbitMQ** for event-driven asynchronous communication. 

```
Order Created → Payment Processed → Delivery Assigned → Notifications Sent
```

Check the logs of the background workers to see the events firing in real-time:
```bash
docker-compose logs -f payment-service delivery-service notification-service
```

## Internal Communication (gRPC)

The system relies on gRPC for fast, strictly-typed synchronous communication between internal services.
- **Order Service** calls **Restaurant Service** (via gRPC) to validate menu items and prices before accepting an order.
- **Order Service** calls **User Service** (via gRPC) to validate user accounts and retrieve default delivery addresses.

## Observability & DevOps

- **Circuit Breaker:** The API Gateway implements the Circuit Breaker pattern (using `gobreaker`) to prevent cascading failures when upstream services are down.
- **Prometheus Metrics:** The API Gateway exposes standard Prometheus HTTP request metrics at `http://localhost:8080/metrics`.
- **Strict Database Separation:** Every single service connects to its own isolated database namespace in MongoDB, preventing direct database coupling.

## Kubernetes Deployment (Kind)

*Note: Ensure your Docker images are built and pushed/loaded to your Kubernetes cluster.*

```bash
# Create cluster
kind create cluster --name food-delivery

# Deploy infrastructure
kubectl apply -f deployments/mongodb.yaml
kubectl apply -f deployments/rabbitmq.yaml

# Deploy services
kubectl apply -f deployments/
```

## License

MIT License
