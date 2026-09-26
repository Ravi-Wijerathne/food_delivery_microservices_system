/**
 * API Client with Intelligent Gateway Fallback & Event Bus Simulator
 * Seamlessly interfaces with:
 * - Real API Gateway (http://localhost:8080/api)
 * - Or Simulated In-Browser Microservices Engine with RabbitMQ Event Pipeline
 */

class ApiService {
    constructor() {
        this.baseUrl = 'http://localhost:8080/api';
        this.mode = 'simulated'; // 'gateway' | 'simulated'
        this.token = localStorage.getItem('crave_token') || '';
        this.lastPingTime = null;
        this.subscribers = [];
        this.orderTimers = new Map();
    }

    setMode(mode) {
        this.mode = mode;
        if (mode === 'docker') {
            this.baseUrl = 'http://localhost:8080/api';
        } else if (mode === 'k8s') {
            this.baseUrl = 'http://localhost:8089/api';
        }
        this.notify('modeChanged', { mode: this.mode });
    }

    isLive() {
        return this.mode === 'docker' || this.mode === 'k8s' || this.mode === 'gateway';
    }

    on(event, callback) {
        this.subscribers.push({ event, callback });
    }

    notify(event, data) {
        this.subscribers
            .filter(sub => sub.event === event)
            .forEach(sub => sub.callback(data));
    }

    async testConnection() {
        // Probe Docker Compose Gateway (:8080)
        try {
            const controller = new AbortController();
            const timeoutId = setTimeout(() => controller.abort(), 800);
            const res = await fetch(`http://localhost:8080/api/health`, { signal: controller.signal });
            clearTimeout(timeoutId);
            if (res.ok) {
                this.baseUrl = 'http://localhost:8080/api';
                this.mode = 'docker';
                this.lastPingTime = Date.now();
                this.notify('connectionStatus', { online: true, mode: 'docker', port: 8080 });
                return { online: true, mode: 'docker', port: 8080 };
            }
        } catch (e) {}

        // Probe Kubernetes Cluster Gateway (:8089)
        try {
            const controller = new AbortController();
            const timeoutId = setTimeout(() => controller.abort(), 800);
            const res = await fetch(`http://localhost:8089/api/health`, { signal: controller.signal });
            clearTimeout(timeoutId);
            if (res.ok) {
                this.baseUrl = 'http://localhost:8089/api';
                this.mode = 'k8s';
                this.lastPingTime = Date.now();
                this.notify('connectionStatus', { online: true, mode: 'k8s', port: 8089 });
                return { online: true, mode: 'k8s', port: 8089 };
            }
        } catch (e) {}

        this.mode = 'simulated';
        this.notify('connectionStatus', { online: false, mode: 'simulated' });
        return { online: false, mode: 'simulated' };
    }

    getHeaders() {
        const headers = { 'Content-Type': 'application/json' };
        if (this.token) {
            headers['Authorization'] = `Bearer ${this.token}`;
        }
        return headers;
    }

    // --- Authentication ---

    async login(email, password) {
        if (this.isLive()) {
            try {
                const res = await fetch(`${this.baseUrl}/auth/login`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ email, password })
                });
                const data = await res.json();
                if (res.ok && data.token) {
                    this.token = data.token;
                    localStorage.setItem('crave_token', this.token);
                    this.logEvent("AuthService", "UserLoggedIn", { email, userId: data.user?.id });
                    return { success: true, token: data.token, user: data.user };
                }
                return { success: false, error: data.error || data.message || 'Login failed' };
            } catch (err) {
                console.warn("Gateway login error, falling back to simulated:", err);
            }
        }

        // Simulated Login
        await this.delay(400);
        const fakeToken = this.generateFakeJWT(email, "user-1");
        this.token = fakeToken;
        localStorage.setItem('crave_token', this.token);
        const user = {
            id: "user-1",
            email: email,
            name: email.split('@')[0].toUpperCase(),
            phone: "+1 (555) 987-6543",
            address: "742 Evergreen Terrace, Apt 4B, Springfield"
        };
        this.logEvent("AuthService", "UserLoggedIn", { email, userId: user.id });
        return { success: true, token: fakeToken, user };
    }

    async register(name, email, password) {
        if (this.isLive()) {
            try {
                const res = await fetch(`${this.baseUrl}/auth/register`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ name, email, password })
                });
                const data = await res.json();
                if (res.ok && data.token) {
                    this.token = data.token;
                    localStorage.setItem('crave_token', this.token);
                    this.logEvent("AuthService", "UserRegistered", { email, userId: data.user?.id });
                    return { success: true, token: data.token, user: data.user };
                }
                return { success: false, error: data.error || data.message || 'Registration failed' };
            } catch (err) {
                console.warn("Gateway register error, falling back to simulated:", err);
            }
        }

        // Simulated Register
        await this.delay(500);
        const fakeToken = this.generateFakeJWT(email, "user-1");
        this.token = fakeToken;
        localStorage.setItem('crave_token', this.token);
        const user = {
            id: "user-" + Math.floor(Math.random() * 1000),
            name: name,
            email: email,
            phone: "+1 (555) 123-4567",
            address: "123 Main Boulevard, Suite 100"
        };
        this.logEvent("AuthService", "UserRegistered", { email, userId: user.id });
        return { success: true, token: fakeToken, user };
    }

    logout() {
        this.token = '';
        localStorage.removeItem('crave_token');
        this.notify('loggedOut');
    }

    // --- User Profile ---

    async getProfile() {
        if (this.isLive()) {
            try {
                const res = await fetch(`${this.baseUrl}/users/profile`, {
                    headers: this.getHeaders()
                });
                if (res.ok) {
                    return await res.json();
                }
            } catch (e) {
                console.warn("Gateway profile error, falling back:", e);
            }
        }

        // Simulated Profile
        const saved = localStorage.getItem('crave_user_profile');
        if (saved) {
            try { return JSON.parse(saved); } catch(e) {}
        }
        return window.MOCK_DATA.defaultUser;
    }

    async updateProfile(profileData) {
        if (this.isLive()) {
            try {
                const res = await fetch(`${this.baseUrl}/users/profile`, {
                    method: 'PUT',
                    headers: this.getHeaders(),
                    body: JSON.stringify(profileData)
                });
                if (res.ok) {
                    this.logEvent("UserService", "ProfileUpdated", profileData);
                    return await res.json();
                }
            } catch (e) {
                console.warn("Gateway update profile error:", e);
            }
        }

        await this.delay(350);
        localStorage.setItem('crave_user_profile', JSON.stringify(profileData));
        this.logEvent("UserService", "ProfileUpdated", { id: profileData.id, name: profileData.name });
        return { message: "Profile updated successfully" };
    }

    // --- Restaurants & Menus ---

    async getRestaurants() {
        if (this.isLive()) {
            try {
                const res = await fetch(`${this.baseUrl}/restaurants`);
                if (res.ok) {
                    const data = await res.json();
                    if (Array.isArray(data) && data.length > 0) {
                        // Merge with rich images & details from mock data
                        return data.map(r => {
                            const mockMatch = window.MOCK_DATA.restaurants.find(m => m.id === r.id);
                            return {
                                ...r,
                                image: mockMatch?.image || 'assets/images/pizza.jpg',
                                reviewsCount: mockMatch?.reviewsCount || 240,
                                deliveryTime: mockMatch?.deliveryTime || '25-35 min',
                                distance: mockMatch?.distance || '2.0 km',
                                priceTier: mockMatch?.priceTier || '$$',
                                minOrder: mockMatch?.minOrder || 15.0,
                                tag: mockMatch?.tag || 'Gourmet Cuisine',
                                badge: mockMatch?.badge || 'Popular'
                            };
                        });
                    }
                }
            } catch (e) {
                console.warn("Gateway fetch restaurants error:", e);
            }
        }

        // Simulated
        return window.MOCK_DATA.restaurants;
    }

    async getMenu(restaurantId) {
        if (this.isLive()) {
            try {
                const res = await fetch(`${this.baseUrl}/menu/${restaurantId}`);
                if (res.ok) {
                    const data = await res.json();
                    if (data && data.items && data.items.length > 0) {
                        return data;
                    }
                }
            } catch (e) {
                console.warn(`Gateway fetch menu for ${restaurantId} failed:`, e);
            }
        }

        // Simulated
        return window.MOCK_DATA.menus[restaurantId] || {
            restaurant_id: restaurantId,
            restaurant_name: "Restaurant",
            items: []
        };
    }

    // --- Orders & RabbitMQ Pipeline ---

    async createOrder(payload) {
        if (this.isLive()) {
            try {
                const res = await fetch(`${this.baseUrl}/orders`, {
                    method: 'POST',
                    headers: this.getHeaders(),
                    body: JSON.stringify(payload)
                });
                if (res.ok) {
                    const order = await res.json();
                    this.logEvent("OrderService", "OrderCreated", {
                        orderId: order.id,
                        userId: order.user_id,
                        amount: order.total_amount,
                        itemsCount: order.items?.length
                    });
                    return order;
                }
            } catch (e) {
                console.warn("Gateway createOrder error, fallback to simulated:", e);
            }
        }

        // --- Simulated Order Execution with RabbitMQ Event Pipeline ---
        await this.delay(350);
        const orderId = "order-" + new Date().toISOString().replace(/[-:T.Z]/g, '').slice(0, 14);
        const total = payload.items.reduce((sum, item) => sum + (item.price * item.quantity), 0);

        const order = {
            id: orderId,
            user_id: payload.user_id || 'user-1',
            restaurant_id: payload.restaurant_id,
            restaurant_name: payload.restaurant_name || 'Gourmet Restaurant',
            items: payload.items,
            total_amount: total,
            delivery_address: payload.delivery_address,
            delivery_notes: payload.delivery_notes || '',
            status: "CREATED",
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
            driver: null,
            pipeline: [
                { stage: "Order Created", service: "order-service", time: new Date().toLocaleTimeString(), done: true, current: false },
                { stage: "Payment Processing", service: "payment-service", time: null, done: false, current: true },
                { stage: "Delivery Assignment", service: "delivery-service", time: null, done: false, current: false },
                { stage: "Out For Delivery", service: "driver-dispatch", time: null, done: false, current: false },
                { stage: "Delivered", service: "customer-confirmed", time: null, done: false, current: false }
            ]
        };

        // Save in local storage
        this.saveSimulatedOrder(order);

        // RabbitMQ: Event 1 -> OrderCreated
        this.logEvent("OrderService", "OrderCreated", {
            orderId: order.id,
            userId: order.user_id,
            amount: `$${total.toFixed(2)}`,
            status: "CREATED"
        });

        // Trigger Async RabbitMQ Consumers simulation
        this.startSimulatedAsyncEventBus(order.id);

        return order;
    }

    async getOrder(orderId) {
        if (this.isLive()) {
            try {
                const res = await fetch(`${this.baseUrl}/orders/${orderId}`, {
                    headers: this.getHeaders()
                });
                if (res.ok) {
                    return await res.json();
                }
            } catch (e) {
                console.warn(`Gateway getOrder for ${orderId} failed:`, e);
            }
        }

        // Simulated Order
        const orders = this.getSimulatedOrders();
        return orders[orderId] || null;
    }

    // --- Simulated Microservices RabbitMQ Flow ---
    startSimulatedAsyncEventBus(orderId) {
        // Clear any existing timer for this order
        if (this.orderTimers.has(orderId)) {
            clearTimeout(this.orderTimers.get(orderId));
        }

        // Step 1: Payment Service handles OrderCreated -> PaymentProcessed (1.5s)
        const t1 = setTimeout(() => {
            const order = this.getSimulatedOrders()[orderId];
            if (!order || order.status === 'CANCELLED') return;

            order.status = "PaymentProcessed";
            order.updated_at = new Date().toISOString();
            if (order.pipeline) {
                order.pipeline[1].done = true;
                order.pipeline[1].current = false;
                order.pipeline[1].time = new Date().toLocaleTimeString();
                order.pipeline[2].current = true;
            }
            this.saveSimulatedOrder(order);

            this.logEvent("PaymentService", "PaymentProcessed", {
                orderId: order.id,
                amount: `$${order.total_amount.toFixed(2)}`,
                status: "SUCCESS",
                channel: "amqp.order.events"
            });
            this.notify('orderUpdated', order);

            // Step 2: Delivery Service handles PaymentProcessed -> DeliveryAssigned (3.2s)
            const t2 = setTimeout(() => {
                const updated = this.getSimulatedOrders()[orderId];
                if (!updated || updated.status === 'CANCELLED') return;

                const drivers = window.MOCK_DATA.drivers;
                const assignedDriver = drivers[Math.floor(Math.random() * drivers.length)];

                updated.status = "DeliveryAssigned";
                updated.driver = assignedDriver;
                updated.updated_at = new Date().toISOString();
                if (updated.pipeline) {
                    updated.pipeline[2].done = true;
                    updated.pipeline[2].current = false;
                    updated.pipeline[2].time = new Date().toLocaleTimeString();
                    updated.pipeline[3].current = true;
                }
                this.saveSimulatedOrder(updated);

                this.logEvent("DeliveryService", "DeliveryAssigned", {
                    orderId: updated.id,
                    agent: assignedDriver.name,
                    vehicle: assignedDriver.vehicle,
                    rating: `⭐ ${assignedDriver.rating}`
                });

                this.logEvent("NotificationService", "SMSNotificationSent", {
                    recipient: updated.user_id,
                    message: `Driver ${assignedDriver.name} is on the way!`
                });

                this.notify('orderUpdated', updated);

                // Step 3: Out for Delivery (6.0s)
                const t3 = setTimeout(() => {
                    const enRoute = this.getSimulatedOrders()[orderId];
                    if (!enRoute || enRoute.status === 'CANCELLED') return;

                    enRoute.status = "OUT_FOR_DELIVERY";
                    if (enRoute.pipeline) {
                        enRoute.pipeline[3].done = true;
                        enRoute.pipeline[3].current = true;
                        enRoute.pipeline[3].time = new Date().toLocaleTimeString();
                    }
                    this.saveSimulatedOrder(enRoute);
                    this.notify('orderUpdated', enRoute);

                    // Step 4: Final DELIVERED (after 20s or user manual trigger)
                    const t4 = setTimeout(() => {
                        const delivered = this.getSimulatedOrders()[orderId];
                        if (!delivered || delivered.status === 'CANCELLED') return;

                        delivered.status = "DELIVERED";
                        if (delivered.pipeline) {
                            delivered.pipeline[3].current = false;
                            delivered.pipeline[4].done = true;
                            delivered.pipeline[4].current = false;
                            delivered.pipeline[4].time = new Date().toLocaleTimeString();
                        }
                        this.saveSimulatedOrder(delivered);

                        this.logEvent("DeliveryService", "OrderDelivered", {
                            orderId: delivered.id,
                            completedAt: new Date().toLocaleTimeString()
                        });
                        this.notify('orderUpdated', delivered);
                    }, 14000);
                    this.orderTimers.set(orderId, t4);
                }, 4000);
                this.orderTimers.set(orderId, t3);
            }, 3000);
            this.orderTimers.set(orderId, t2);
        }, 1600);
        this.orderTimers.set(orderId, t1);
    }

    forceDeliverOrder(orderId) {
        const orders = this.getSimulatedOrders();
        const order = orders[orderId];
        if (order) {
            order.status = "DELIVERED";
            if (order.pipeline) {
                order.pipeline.forEach(p => { p.done = true; p.current = false; });
                order.pipeline[4].time = new Date().toLocaleTimeString();
            }
            this.saveSimulatedOrder(order);
            this.logEvent("DeliveryService", "OrderDelivered", { orderId, manualOverride: true });
            this.notify('orderUpdated', order);
        }
    }

    // --- Observability & Metrics ---

    async getGatewayMetrics() {
        if (this.isLive()) {
            try {
                const res = await fetch(`${this.baseUrl}/metrics`);
                if (res.ok) {
                    return await res.text();
                }
            } catch (e) {}
        }
        return `# HELP http_requests_total Total number of HTTP requests made.
# TYPE http_requests_total counter
http_requests_total{code="200",handler="/api/restaurants",method="get"} 438
http_requests_total{code="200",handler="/api/orders",method="post"} 129
http_requests_total{code="200",handler="/api/users/profile",method="get"} 312
# HELP gobreaker_circuit_breaker Circuit Breaker State (0=closed, 1=half-open, 2=open)
# TYPE gobreaker_circuit_breaker gauge
gobreaker_circuit_breaker{name="Gateway_Circuit_Breaker"} 0`;
    }

    async getServiceHealth(port) {
        try {
            const controller = new AbortController();
            const timeoutId = setTimeout(() => controller.abort(), 600);
            const res = await fetch(`http://localhost:${port}/health`, { signal: controller.signal });
            clearTimeout(timeoutId);
            if (res.ok) {
                return { status: "ONLINE", latency: Math.floor(Math.random() * 8) + 2 };
            }
        } catch (e) {}
        return { status: "SIMULATED", latency: Math.floor(Math.random() * 5) + 1 };
    }

    // --- Helper Utilities ---

    saveSimulatedOrder(order) {
        const orders = this.getSimulatedOrders();
        orders[order.id] = order;
        localStorage.setItem('crave_simulated_orders', JSON.stringify(orders));
    }

    getSimulatedOrders() {
        try {
            return JSON.parse(localStorage.getItem('crave_simulated_orders') || '{}');
        } catch (e) {
            return {};
        }
    }

    logEvent(service, eventType, payload) {
        const event = {
            id: 'evt-' + Math.random().toString(36).substr(2, 9),
            service,
            eventType,
            payload,
            timestamp: new Date().toLocaleTimeString(),
            rawTime: Date.now()
        };

        // Notify subscribers (Event Stream panel & Toasts)
        this.notify('busEvent', event);

        // Save event history
        try {
            const history = JSON.parse(localStorage.getItem('crave_event_log') || '[]');
            history.unshift(event);
            if (history.length > 50) history.pop();
            localStorage.setItem('crave_event_log', JSON.stringify(history));
        } catch(e) {}
    }

    getEventLog() {
        try {
            return JSON.parse(localStorage.getItem('crave_event_log') || '[]');
        } catch(e) {
            return [];
        }
    }

    generateFakeJWT(email, userId) {
        const header = btoa(JSON.stringify({ alg: "HS256", typ: "JWT" }));
        const payload = btoa(JSON.stringify({
            user_id: userId,
            email: email,
            role: "customer",
            exp: Math.floor(Date.now() / 1000) + (24 * 60 * 60)
        }));
        const signature = btoa("mock_hmac_sha256_signature_verified");
        return `${header}.${payload}.${signature}`;
    }

    delay(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }
}

window.apiService = new ApiService();
