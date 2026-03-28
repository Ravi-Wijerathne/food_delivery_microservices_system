# System Architecture

## 🧭 High-Level Architecture

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

---

## 🔁 Request Flow (REST + gRPC)

```
Client → API Gateway → (REST)
                ↓
        Internal Services → (gRPC)
```

---

## 🔄 Event-Driven Flow

```
[Order Service]
    |
    | Publish → OrderCreated
    ↓
[RabbitMQ] → [Payment Service]
                |
                | Publish → PaymentProcessed
                ↓
           [Order Service]
                |
                | Publish → OrderConfirmed
                ↓
         [Delivery Service]
                |
                | Publish → DeliveryAssigned
                ↓
     [Notification Service]
```

---

## 📦 Service Isolation

Each service:

* Own codebase
* Own database
* Own container
* Independent deployment

---

## 🗄️ Database Architecture

```
Auth Service        → MongoDB (auth-db)
User Service        → MongoDB (user-db)
Restaurant Service  → MongoDB (restaurant-db)
Order Service       → MongoDB (order-db)
Payment Service     → MongoDB (payment-db)
Delivery Service    → MongoDB (delivery-db)
```

---

## 🧱 Kubernetes Architecture

```
[Kubernetes Cluster]

|-- API Gateway (Deployment + Service)
|-- Auth Service
|-- User Service
|-- Restaurant Service
|-- Order Service
|-- Payment Service
|-- Delivery Service
|-- Notification Service
|
|-- RabbitMQ (Stateful)
|-- MongoDB (per service or shared cluster)
|
|-- Ingress Controller
```

---

## 🌐 Ingress Routing

```
/api/auth        → Auth Service
/api/users       → User Service
/api/restaurants → Restaurant Service
/api/orders      → Order Service
```

---

## 🔐 Security Flow

```
Client → Login → Auth Service
       ← JWT Token

Client → Request → API Gateway
       → Validate JWT
       → Forward request
```

---

## ⚙️ Tech Stack

* Go (Golang)
* gRPC (internal communication)
* RabbitMQ (event system)
* MongoDB (databases)
* Docker (containerization)
* Kubernetes (orchestration)

---

## 📡 Future Enhancements

* Prometheus (metrics)
* Grafana (dashboard)
* Circuit breaker pattern
* Rate limiting

---
