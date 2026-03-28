# Food Delivery Microservices System — Project Specification

## 🧩 Overview

This project is a distributed food delivery backend system built using a microservices architecture. It simulates real-world platforms like Uber Eats, focusing on scalability, service isolation, and event-driven communication.

## 🎯 Goals

* Demonstrate Microservices Architecture
* Implement Event-Driven Communication
* Use Kubernetes for orchestration
* Apply production-level patterns (API Gateway, Auth, messaging)

---

## 🧱 Core Services

### 1. API Gateway

* Entry point for all clients
* Routes requests to internal services
* Handles authentication (JWT validation)

---

### 2. Auth Service

* User registration/login
* JWT token generation
* Password hashing

---

### 3. User Service

* Stores user profiles
* Manages user data

---

### 4. Restaurant Service

* Manages restaurants and menus
* Provides menu browsing APIs

---

### 5. Order Service (Core)

* Handles order creation
* Maintains order lifecycle
* Publishes events

---

### 6. Payment Service

* Simulates payment processing
* Emits success/failure events

---

### 7. Delivery Service

* Assigns delivery agents
* Tracks delivery status

---

### 8. Notification Service

* Sends system notifications (logs/mock)

---

## 🔄 Communication Model

### External

* REST (Client → API Gateway)

### Internal

* gRPC (Service-to-service)
* RabbitMQ (event-driven async communication)

---

## 📦 Event Flow

* OrderCreated
* PaymentProcessed
* OrderConfirmed
* DeliveryAssigned
* OrderDelivered

---

## 🗄️ Database Strategy

* Each service has its own MongoDB database
* No shared database
* Loose coupling ensured

---

## 🔐 Authentication

* JWT-based authentication
* API Gateway validates tokens
* Services trust Gateway

---

## 🌐 Client

* Simple Web UI (basic dashboard)
* Features:

  * Browse restaurants
  * Place orders
  * Track status

---

## 🚀 Deployment

* Dockerized services
* Kubernetes deployment:

  * Deployments
  * Services
  * Ingress
  * ConfigMaps & Secrets

---

## 📈 Non-Functional Requirements

* Scalability (horizontal scaling via K8s)
* Fault tolerance (retry mechanisms)
* Observability-ready (logs + metrics later)

---

## 🧠 Key Concepts Demonstrated

* Service isolation
* Event-driven architecture
* API Gateway pattern
* Database per service
* Container orchestration

---
