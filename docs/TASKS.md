# TASKS.md — OpenCode Execution Plan

## 🧠 Execution Strategy

* Work **service by service**
* Keep commits small and isolated
* Always verify before moving forward
* Prefer working system > perfect system

---

# 🟢 PHASE 0 — Project Initialization

## Task 0.1 — Create Root Structure

* Create root project folder
* Initialize Git repository
* Create folders:

  * /services
  * /gateway
  * /proto
  * /deployments
  * /web

---

## Task 0.2 — Setup Go Workspace

* Initialize Go modules for:

  * each service
  * gateway
* Ensure independent modules per service

---

## Task 0.3 — Setup Base Tooling

* Install:

  * protoc (Protocol Buffers)
  * Docker
  * Kubernetes (Kind)
* Verify installations

---

# 🟢 PHASE 1 — Auth Service (Start Simple)

## Task 1.1 — Create Auth Service

* Create `/services/auth-service`
* Initialize Go module
* Setup basic HTTP server

---

## Task 1.2 — Implement Endpoints

* POST /register
* POST /login

---

## Task 1.3 — JWT Logic

* Generate JWT tokens
* Validate tokens

---

## Task 1.4 — Password Handling

* Hash passwords (bcrypt)
* Store securely

---

## Task 1.5 — Test Service

* Test via curl
* Ensure login/register works

---

# 🟢 PHASE 2 — Restaurant Service

## Task 2.1 — Setup Service

* Create `/services/restaurant-service`
* Setup HTTP server

---

## Task 2.2 — Implement Features

* GET /restaurants
* GET /menu/{id}

---

## Task 2.3 — Mock Data

* Use in-memory data initially

---

## Task 2.4 — Test Endpoints

---

# 🟢 PHASE 3 — Order Service (Core)

## Task 3.1 — Setup Service

* Create `/services/order-service`

---

## Task 3.2 — Implement Order Logic

* POST /orders
* GET /orders/{id}

---

## Task 3.3 — Order State Machine

States:

* CREATED
* CONFIRMED
* PREPARING
* DELIVERED

---

## Task 3.4 — Store Orders (in-memory first)

---

# 🟡 PHASE 4 — API Gateway

## Task 4.1 — Create Gateway

* `/gateway`
* Setup HTTP reverse proxy

---

## Task 4.2 — Routing

* Route requests to services:

  * /auth → auth-service
  * /restaurants → restaurant-service
  * /orders → order-service

---

## Task 4.3 — JWT Middleware

* Validate token
* Reject unauthorized requests

---

## Task 4.4 — Test Full Flow

* Client → Gateway → Services

---

# 🟡 PHASE 5 — gRPC Setup

## Task 5.1 — Define Protobuf

* Create `/proto`
* Define:

  * order.proto
  * payment.proto

---

## Task 5.2 — Generate Code

* Compile proto → Go code

---

## Task 5.3 — Implement gRPC Server

* Add gRPC to Order Service

---

## Task 5.4 — Implement gRPC Client

* Enable service-to-service calls

---

# 🟠 PHASE 6 — RabbitMQ Integration

## Task 6.1 — Setup RabbitMQ (Docker)

* Run RabbitMQ container
* Verify dashboard

---

## Task 6.2 — Event Publisher

* Order Service publishes:

  * OrderCreated

---

## Task 6.3 — Event Consumer

* Payment Service consumes events

---

## Task 6.4 — Define Event Schema

* JSON format:

  * order_id
  * status
  * timestamp

---

# 🟠 PHASE 7 — Payment Service

## Task 7.1 — Create Service

* `/services/payment-service`

---

## Task 7.2 — Consume Events

* Listen for OrderCreated

---

## Task 7.3 — Process Payment

* Simulate success/failure

---

## Task 7.4 — Publish Result

* Emit PaymentProcessed event

---

# 🟠 PHASE 8 — Delivery Service

## Task 8.1 — Create Service

* `/services/delivery-service`

---

## Task 8.2 — Consume Events

* Listen for PaymentProcessed

---

## Task 8.3 — Assign Delivery

* Simulate rider assignment

---

## Task 8.4 — Emit Event

* DeliveryAssigned

---

# 🔴 PHASE 9 — Notification Service

## Task 9.1 — Create Service

* `/services/notification-service`

---

## Task 9.2 — Listen to Events

* OrderConfirmed
* DeliveryAssigned

---

## Task 9.3 — Output Notifications

* Log messages to console

---

# 🔴 PHASE 10 — MongoDB Integration

## Task 10.1 — Setup MongoDB

* Run MongoDB container

---

## Task 10.2 — Integrate per Service

* Replace in-memory storage

---

## Task 10.3 — Schema Design

* Define collections per service

---

# 🔵 PHASE 11 — Dockerization

## Task 11.1 — Dockerfile per Service

* Write minimal Dockerfile

---

## Task 11.2 — Build Images

* Build all services

---

## Task 11.3 — Test Containers

* Run locally

---

# 🔵 PHASE 12 — Kubernetes Deployment

## Task 12.1 — Create YAMLs

* Deployment
* Service

---

## Task 12.2 — Deploy Services

* Apply configs

---

## Task 12.3 — Setup Ingress

* Route external traffic

---

## Task 12.4 — Config Management

* ConfigMaps
* Secrets

---

# 🟣 PHASE 13 — Web UI

## Task 13.1 — Setup Frontend

* Basic UI (React or static)

---

## Task 13.2 — Features

* View restaurants
* Place order
* View status

---

## Task 13.3 — Connect Gateway

---

# ⚫ PHASE 14 — Final Enhancements

## Task 14.1 — Error Handling

* Retry failed events

---

## Task 14.2 — Logging

* Structured logs

---

## Task 14.3 — Cleanup

* Refactor code
* Improve structure

---

# 🚀 FINAL RESULT

* Fully working microservices system
* Event-driven architecture
* Kubernetes deployment ready
* Strong portfolio project

---
