/**
 * CravePulse UI Controller & Application Orchestrator
 */

class AppController {
    constructor() {
        this.api = window.apiService;
        this.store = window.store;
        this.orderPollingInterval = null;
        this.activeAuthTab = 'login';
        this.selectedCategory = 'all';
    }

    async init() {
        // Subscribe to store changes
        this.store.subscribe(state => this.handleStateChange(state));

        // Connect API event listeners
        this.api.on('busEvent', event => this.handleBusEvent(event));
        this.api.on('orderUpdated', order => this.handleOrderUpdate(order));
        this.api.on('connectionStatus', status => this.updateConnectionUI(status));

        // Initialize user session
        await this.initSession();

        // Check backend gateway connection
        await this.api.testConnection();

        // Load restaurants
        await this.loadRestaurants();

        // Check for active order
        const activeOrder = this.store.getState().activeOrder;
        if (activeOrder) {
            this.startOrderTracking(activeOrder);
        }

        // Render initial UI
        this.setupEventListeners();
        this.renderAll();

        console.log("CravePulse Microservices GUI Initialized Successfully");
    }

    async initSession() {
        const profile = await this.api.getProfile();
        this.store.setState({ currentUser: profile });
    }

    // --- Data Fetching ---

    async loadRestaurants() {
        const restaurants = await this.api.getRestaurants();
        this.store.setState({ restaurants });
        this.renderRestaurants();
    }

    // --- UI State & Navigation ---

    switchTab(tabName) {
        this.store.setState({ activeTab: tabName });
        document.querySelectorAll('.view-section').forEach(el => el.classList.remove('active'));
        document.querySelectorAll('.nav-btn').forEach(el => el.classList.remove('active'));

        const targetSection = document.getElementById(`${tabName}View`);
        if (targetSection) targetSection.classList.add('active');

        const navBtn = document.getElementById(`navBtn_${tabName}`);
        if (navBtn) navBtn.classList.add('active');

        if (tabName === 'observability') {
            this.renderObservabilityHub();
        } else if (tabName === 'orders') {
            this.renderOrdersView();
        }
    }

    // --- Rendering Restaurants ---

    renderRestaurants() {
        const { restaurants, selectedCuisine, searchQuery } = this.store.getState();
        const container = document.getElementById('restaurantsGrid');
        if (!container) return;

        let filtered = restaurants.filter(r => {
            const matchesCuisine = selectedCuisine === 'all' || r.cuisine.toLowerCase() === selectedCuisine.toLowerCase();
            const matchesSearch = !searchQuery || 
                r.name.toLowerCase().includes(searchQuery.toLowerCase()) || 
                r.cuisine.toLowerCase().includes(searchQuery.toLowerCase());
            return matchesCuisine && matchesSearch;
        });

        document.getElementById('restaurantCount').textContent = `${filtered.length} places near you`;

        if (filtered.length === 0) {
            container.innerHTML = `
                <div style="grid-column: 1 / -1; text-align: center; padding: 60px 20px;">
                    <p style="font-size: 2.5rem; margin-bottom: 12px;">🔍</p>
                    <h3 style="font-size: 1.4rem; margin-bottom: 8px;">No restaurants found</h3>
                    <p style="color: var(--text-muted);">Try adjusting your search query or cuisine filter.</p>
                </div>
            `;
            return;
        }

        container.innerHTML = filtered.map(r => `
            <div class="restaurant-card" onclick="app.openMenuModal('${r.id}')">
                <div class="restaurant-card-image-box">
                    <img src="${r.image}" alt="${r.name}" class="restaurant-card-image" loading="lazy">
                    <div class="card-top-badges">
                        <span class="badge-tag highlight">${r.badge || r.cuisine}</span>
                        <span class="rating-badge">★ ${r.rating}</span>
                    </div>
                </div>
                <div class="restaurant-card-body">
                    <h3 class="restaurant-card-title">${r.name}</h3>
                    <p class="restaurant-card-tagline">${r.tag || r.cuisine}</p>
                    <div class="restaurant-meta-row">
                        <span class="meta-item">⏱️ ${r.deliveryTime}</span>
                        <span class="meta-item">📍 ${r.distance}</span>
                        <span class="meta-item" style="margin-left: auto; font-weight: 700; color: var(--text-primary);">${r.priceTier}</span>
                    </div>
                </div>
            </div>
        `).join('');
    }

    // --- Menu Modal ---

    async openMenuModal(restaurantId) {
        const restaurant = this.store.getState().restaurants.find(r => r.id === restaurantId);
        if (!restaurant) return;

        this.store.setState({ activeRestaurant: restaurant, isMenuModalOpen: true });
        const modal = document.getElementById('menuModal');
        modal.classList.add('open');

        document.getElementById('modalRestImg').src = restaurant.image;
        document.getElementById('modalRestName').textContent = restaurant.name;
        document.getElementById('modalRestAddress').textContent = `📍 ${restaurant.address} • ⭐ ${restaurant.rating} (${restaurant.reviewsCount} reviews)`;

        const menuContainer = document.getElementById('modalMenuItems');
        menuContainer.innerHTML = `
            <div style="text-align: center; padding: 40px;">
                <p style="color: var(--text-secondary);">Querying Restaurant Service menu catalog...</p>
            </div>
        `;

        const menuData = await this.api.getMenu(restaurantId);
        this.store.setState({ currentMenu: menuData });
        this.renderMenuItems(restaurant, menuData);
    }

    renderMenuItems(restaurant, menuData) {
        const container = document.getElementById('modalMenuItems');
        if (!container || !menuData || !menuData.items) return;

        // Group by category
        const categories = {};
        menuData.items.forEach(item => {
            const cat = item.category || 'Specialties';
            if (!categories[cat]) categories[cat] = [];
            categories[cat].push(item);
        });

        container.innerHTML = Object.entries(categories).map(([catName, items]) => `
            <div class="menu-category-section">
                <h4 class="menu-category-title">${catName}</h4>
                <div class="menu-items-grid">
                    ${items.map(item => `
                        <div class="menu-item-row">
                            <div class="dish-info">
                                <div class="dish-header-row">
                                    <h5 class="dish-title">${item.name}</h5>
                                    ${item.isVeg ? '<span style="font-size: 0.75rem; background: rgba(16, 185, 129, 0.2); color: #34d399; padding: 2px 8px; border-radius: 4px; font-weight: 700;">VEG</span>' : ''}
                                    ${item.isSpicy ? '<span style="font-size: 0.75rem; background: rgba(239, 68, 68, 0.2); color: #f87171; padding: 2px 8px; border-radius: 4px; font-weight: 700;">🌶️ SPICY</span>' : ''}
                                </div>
                                <p class="dish-desc">${item.description}</p>
                                <div class="dish-price">$${item.price.toFixed(2)}</div>
                            </div>
                            <button class="add-cart-btn" onclick="app.addItemToCart('${restaurant.id}', '${item.id}')">
                                <span>Add</span>
                                <strong>+</strong>
                            </button>
                        </div>
                    `).join('')}
                </div>
            </div>
        `).join('');
    }

    closeMenuModal() {
        document.getElementById('menuModal').classList.remove('open');
        this.store.setState({ isMenuModalOpen: false, activeRestaurant: null });
    }

    // --- Cart & Checkout Drawer ---

    addItemToCart(restaurantId, itemId) {
        const restaurant = this.store.getState().restaurants.find(r => r.id === restaurantId);
        const menu = this.store.getState().currentMenu;
        if (!restaurant || !menu) return;

        const item = menu.items.find(i => i.id === itemId);
        if (!item) return;

        const success = this.store.addToCart(restaurant, item, 1);
        if (success) {
            this.showToast(`Added ${item.name} to cart`, 'success');
            this.renderCartUI();
        }
    }

    toggleCart(open = null) {
        const overlay = document.getElementById('cartOverlay');
        const isOpen = open !== null ? open : !overlay.classList.contains('open');
        overlay.classList.toggle('open', isOpen);
        this.store.setState({ isCartOpen: isOpen });
        if (isOpen) {
            this.renderCartUI();
        }
    }

    renderCartUI() {
        const { cart, currentUser } = this.store.getState();
        const itemsContainer = document.getElementById('cartItemsList');
        const emptyState = document.getElementById('cartEmptyState');
        const cartFooter = document.getElementById('cartFooter');
        const totals = this.store.getCartTotals();

        // Update badge counter
        document.getElementById('cartBadgeCount').textContent = totals.itemCount;

        if (totals.itemCount === 0) {
            itemsContainer.style.display = 'none';
            cartFooter.style.display = 'none';
            emptyState.style.display = 'block';
            return;
        }

        emptyState.style.display = 'none';
        itemsContainer.style.display = 'flex';
        cartFooter.style.display = 'block';

        document.getElementById('cartRestTitle').textContent = `From ${cart.restaurantName}`;

        itemsContainer.innerHTML = cart.items.map(item => `
            <div class="cart-item-card">
                <div class="cart-item-left">
                    <h4>${item.name}</h4>
                    <p>$${(item.price * item.quantity).toFixed(2)}</p>
                </div>
                <div class="qty-control">
                    <button class="qty-btn" onclick="app.updateItemQty('${item.id}', -1)">−</button>
                    <span class="qty-display">${item.quantity}</span>
                    <button class="qty-btn" onclick="app.updateItemQty('${item.id}', 1)">+</button>
                </div>
            </div>
        `).join('');

        document.getElementById('cartSubtotal').textContent = `$${totals.subtotal.toFixed(2)}`;
        document.getElementById('cartDeliveryFee').textContent = totals.deliveryFee === 0 ? 'FREE' : `$${totals.deliveryFee.toFixed(2)}`;
        document.getElementById('cartTax').textContent = `$${totals.tax.toFixed(2)}`;
        document.getElementById('cartTotal').textContent = `$${totals.total.toFixed(2)}`;

        // Delivery address input
        const addressInput = document.getElementById('checkoutAddressInput');
        if (addressInput && !addressInput.value) {
            addressInput.value = currentUser?.address || '742 Evergreen Terrace, Apt 4B, Springfield';
        }
    }

    updateItemQty(itemId, delta) {
        this.store.updateCartItemQuantity(itemId, delta);
        this.renderCartUI();
    }

    async handleCheckout() {
        const { cart, currentUser } = this.store.getState();
        if (cart.items.length === 0) return;

        const addressInput = document.getElementById('checkoutAddressInput').value;
        const notesInput = document.getElementById('checkoutNotesInput')?.value || '';

        const checkoutBtn = document.getElementById('checkoutSubmitBtn');
        checkoutBtn.innerHTML = `<span>Publishing to RabbitMQ...</span>`;
        checkoutBtn.disabled = true;

        const payload = {
            user_id: currentUser?.id || 'user-1',
            restaurant_id: cart.restaurantId,
            restaurant_name: cart.restaurantName,
            delivery_address: addressInput,
            delivery_notes: notesInput,
            items: cart.items.map(i => ({
                menu_item_id: i.id,
                name: i.name,
                quantity: i.quantity,
                price: i.price
            }))
        };

        try {
            const order = await this.api.createOrder(payload);
            this.store.clearCart();
            this.store.setActiveOrder(order);
            this.toggleCart(false);

            this.showToast(`Order #${order.id} placed! Tracking started.`, 'success');
            this.switchTab('orders');
            this.startOrderTracking(order);
        } catch (err) {
            this.showToast("Order placement failed", "warning");
        } finally {
            checkoutBtn.innerHTML = `<span>Place Order (RabbitMQ Bus)</span> <span>→</span>`;
            checkoutBtn.disabled = false;
        }
    }

    // --- Order Tracking & Pipeline Visualizer ---

    startOrderTracking(order) {
        this.renderActiveOrderUI(order);

        if (this.orderPollingInterval) {
            clearInterval(this.orderPollingInterval);
        }

        // Live polling if connected to real Gateway
        if (this.api.mode === 'gateway') {
            this.orderPollingInterval = setInterval(async () => {
                const updated = await this.api.getOrder(order.id);
                if (updated) {
                    this.store.setActiveOrder(updated);
                    this.renderActiveOrderUI(updated);
                    if (updated.status === 'DELIVERED' || updated.status === 'CANCELLED') {
                        clearInterval(this.orderPollingInterval);
                    }
                }
            }, 2000);
        }
    }

    handleOrderUpdate(order) {
        const active = this.store.getState().activeOrder;
        if (active && active.id === order.id) {
            this.store.setActiveOrder(order);
            this.renderActiveOrderUI(order);
            this.renderOrdersView();
        }
    }

    renderActiveOrderUI(order) {
        const trackingCard = document.getElementById('activeOrderCard');
        if (!trackingCard || !order) return;

        document.getElementById('noActiveOrderMsg').style.display = 'none';
        trackingCard.style.display = 'block';

        document.getElementById('trackOrderId').textContent = order.id;
        document.getElementById('trackRestName').textContent = order.restaurant_name || 'Restaurant';
        document.getElementById('trackTotal').textContent = `$${order.total_amount ? order.total_amount.toFixed(2) : '0.00'}`;
        document.getElementById('trackItemsSummary').textContent = (order.items || []).map(i => `${i.quantity}x ${i.name}`).join(', ');

        const statusPill = document.getElementById('trackStatusPill');
        statusPill.textContent = order.status;
        statusPill.className = `live-status-pill ${order.status.toLowerCase()}`;

        // Render Microservices Pipeline
        const pipelineStages = [
            { name: "Order Created", service: "order-service (Port 8083)", trigger: "OrderCreated event to RabbitMQ" },
            { name: "Payment Processed", service: "payment-service (Port 8085)", trigger: "Consumed OrderCreated, debited balance" },
            { name: "Delivery Assigned", service: "delivery-service (Port 8086)", trigger: "Assigned available courier agent" },
            { name: "Out For Delivery", service: "fleet-dispatch", trigger: "Driver en route with thermal bag" },
            { name: "Order Delivered", service: "customer-confirmed", trigger: "Handed over at delivery address" }
        ];

        let currentIdx = 0;
        if (order.status === 'PaymentProcessed' || order.status === 'CONFIRMED') currentIdx = 1;
        else if (order.status === 'DeliveryAssigned') currentIdx = 2;
        else if (order.status === 'OUT_FOR_DELIVERY') currentIdx = 3;
        else if (order.status === 'DELIVERED') currentIdx = 4;

        const pipelineContainer = document.getElementById('eventPipelineList');
        pipelineContainer.innerHTML = pipelineStages.map((stage, idx) => {
            let stateClass = '';
            let icon = idx + 1;
            if (idx < currentIdx || (idx === 4 && order.status === 'DELIVERED')) {
                stateClass = 'done';
                icon = '✓';
            } else if (idx === currentIdx && order.status !== 'DELIVERED') {
                stateClass = 'current';
            }

            return `
                <div class="pipeline-step ${stateClass}">
                    <div class="pipeline-node">${icon}</div>
                    <div class="pipeline-details">
                        <h5>${stage.name}</h5>
                        <p style="color: var(--accent-cyan); font-weight: 600;">${stage.service}</p>
                        <p>${stage.trigger}</p>
                    </div>
                </div>
            `;
        }).join('');

        // Driver section
        const driverBox = document.getElementById('driverCardBox');
        if (order.driver) {
            driverBox.style.display = 'flex';
            document.getElementById('driverAvatar').textContent = order.driver.avatar || '🛵';
            document.getElementById('driverName').textContent = order.driver.name;
            document.getElementById('driverVehicle').textContent = `${order.driver.vehicle} • ⭐ ${order.driver.rating}`;
            document.getElementById('driverPhone').textContent = order.driver.phone;
        } else {
            driverBox.style.display = 'none';
        }

        // Render Animated SVG Route
        this.renderRouteMap(currentIdx);
    }

    renderRouteMap(progressStage) {
        const mapContainer = document.getElementById('routeMapContainer');
        if (!mapContainer) return;

        const percentage = Math.min(100, Math.max(10, progressStage * 25));

        mapContainer.innerHTML = `
            <svg width="100%" height="100%" viewBox="0 0 500 240" fill="none" xmlns="http://www.w3.org/2000/svg">
                <!-- Map Background Grid -->
                <defs>
                    <pattern id="grid" width="30" height="30" patternUnits="userSpaceOnUse">
                        <path d="M 30 0 L 0 0 0 30" fill="none" stroke="rgba(255,255,255,0.03)" stroke-width="1"/>
                    </pattern>
                    <linearGradient id="routeGradient" x1="0" y1="0" x2="1" y2="0">
                        <stop offset="0%" stop-color="#ff4b2b" />
                        <stop offset="100%" stop-color="#10b981" />
                    </linearGradient>
                </defs>
                <rect width="100%" height="100%" fill="#0a0e17" />
                <rect width="100%" height="100%" fill="url(#grid)" />

                <!-- City Roads -->
                <path d="M 30 180 Q 150 90 250 150 T 470 70" stroke="rgba(255,255,255,0.08)" stroke-width="14" stroke-linecap="round" fill="none"/>
                
                <!-- Active Route -->
                <path d="M 30 180 Q 150 90 250 150 T 470 70" stroke="url(#routeGradient)" stroke-width="5" stroke-linecap="round" fill="none" stroke-dasharray="8 6"/>

                <!-- Restaurant Origin Pin -->
                <circle cx="30" cy="180" r="14" fill="#ff4b2b" />
                <text x="30" y="184" text-anchor="middle" font-size="12" fill="white">🏪</text>
                <text x="30" y="210" text-anchor="middle" font-size="11" font-weight="700" fill="#cbd5e1">Restaurant</text>

                <!-- Customer Destination Pin -->
                <circle cx="470" cy="70" r="14" fill="#10b981" />
                <text x="470" y="74" text-anchor="middle" font-size="12" fill="white">🏠</text>
                <text x="470" y="100" text-anchor="middle" font-size="11" font-weight="700" fill="#cbd5e1">Delivery Spot</text>

                <!-- Moving Courier Vehicle Pin -->
                <g transform="translate(${30 + (percentage / 100) * 440}, ${180 - (percentage / 100) * 110})">
                    <circle cx="0" cy="0" r="16" fill="#f59e0b" filter="drop-shadow(0 0 8px rgba(245, 158, 11, 0.8))"/>
                    <text x="0" y="4" text-anchor="middle" font-size="13">🛵</text>
                </g>
            </svg>
        `;
    }

    renderOrdersView() {
        const history = this.store.getState().orderHistory;
        const container = document.getElementById('pastOrdersGrid');
        if (!container) return;

        if (history.length === 0) {
            container.innerHTML = `<p style="color: var(--text-muted); padding: 20px;">No previous orders found.</p>`;
            return;
        }

        container.innerHTML = history.map(order => `
            <div class="glass-card" style="padding: 20px; display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px;">
                <div>
                    <div style="display: flex; align-items: center; gap: 10px; margin-bottom: 4px;">
                        <strong style="font-family: var(--font-heading); font-size: 1.1rem;">${order.restaurant_name || 'Restaurant'}</strong>
                        <span class="live-status-pill ${order.status.toLowerCase()}" style="font-size: 0.72rem; padding: 2px 10px;">${order.status}</span>
                    </div>
                    <p style="font-size: 0.82rem; color: var(--text-muted);">${order.id} • ${new Date(order.created_at).toLocaleDateString()}</p>
                    <p style="font-size: 0.85rem; color: var(--text-secondary); margin-top: 4px;">
                        ${(order.items || []).map(i => `${i.quantity}x ${i.name}`).join(', ')}
                    </p>
                </div>
                <div style="text-align: right;">
                    <div style="font-size: 1.2rem; font-weight: 800; color: var(--accent-orange);">$${order.total_amount ? order.total_amount.toFixed(2) : '0.00'}</div>
                    <button class="nav-btn" style="margin-top: 8px; font-size: 0.8rem; padding: 6px 14px;" onclick="app.viewPastOrder('${order.id}')">View</button>
                </div>
            </div>
        `).join('');
    }

    viewPastOrder(orderId) {
        const order = this.store.getState().orderHistory.find(o => o.id === orderId);
        if (order) {
            this.store.setActiveOrder(order);
            this.startOrderTracking(order);
            window.scrollTo({ top: 0, behavior: 'smooth' });
        }
    }

    // --- Microservices Observability Hub ---

    renderObservabilityHub() {
        const container = document.getElementById('servicesGrid');
        if (!container) return;

        const services = window.MOCK_DATA.microservices;
        container.innerHTML = services.map(s => `
            <div class="service-card">
                <div class="service-top-row">
                    <div class="service-icon-wrap">${s.icon}</div>
                    <div class="service-status-dot">
                        <span style="width: 8px; height: 8px; border-radius: 50%; background: var(--accent-emerald); box-shadow: 0 0 8px var(--accent-emerald);"></span>
                        <span>${s.status}</span>
                    </div>
                </div>
                <h4 class="service-title">${s.name}</h4>
                <p class="service-desc">${s.role}</p>
                <div class="service-meta-footer">
                    <span>Port: <strong>${s.port}${s.grpcPort ? ` (gRPC ${s.grpcPort})` : ''}</strong></span>
                    <span>${s.protocol}</span>
                </div>
            </div>
        `).join('');

        this.renderTerminalLogs();
    }

    renderTerminalLogs() {
        const logsContainer = document.getElementById('terminalLogEntries');
        if (!logsContainer) return;

        const logs = this.api.getEventLog();
        if (logs.length === 0) {
            logsContainer.innerHTML = `<div style="color: var(--text-muted);">Listening to RabbitMQ AMQP message bus [amqp://rabbitmq:5672]...</div>`;
            return;
        }

        logsContainer.innerHTML = logs.map(l => `
            <div class="log-entry">
                <span class="log-time">[${l.timestamp}]</span>
                <span class="log-service">&lt;${l.service}&gt;</span>
                <span class="log-event">${l.eventType}</span>
                <span class="log-payload">${JSON.stringify(l.payload)}</span>
            </div>
        `).join('');
    }

    handleBusEvent(event) {
        this.renderTerminalLogs();
        if (event.service === 'PaymentService' || event.service === 'DeliveryService') {
            this.showToast(`[${event.service}] ${event.eventType}`, 'info');
        }
    }

    // --- Toasts & Notifications ---

    showToast(message, type = 'info') {
        const container = document.getElementById('toastContainer');
        if (!container) return;

        const toast = document.createElement('div');
        toast.className = `toast ${type}`;
        const icon = type === 'success' ? '✅' : type === 'warning' ? '⚠️' : '⚡';
        toast.innerHTML = `<span>${icon}</span><span style="font-weight: 600; font-size: 0.9rem;">${message}</span>`;
        container.appendChild(toast);

        setTimeout(() => {
            toast.style.opacity = '0';
            toast.style.transform = 'translateY(10px)';
            setTimeout(() => toast.remove(), 300);
        }, 3200);
    }

    // --- Profile & Authentication Modal ---

    openAuthModal(tab = 'login') {
        this.activeAuthTab = tab;
        document.getElementById('authModal').classList.add('open');
        this.switchAuthTab(tab);
    }

    closeAuthModal() {
        document.getElementById('authModal').classList.remove('open');
    }

    switchAuthTab(tab) {
        this.activeAuthTab = tab;
        document.getElementById('authTabLogin').classList.toggle('active', tab === 'login');
        document.getElementById('authTabRegister').classList.toggle('active', tab === 'register');
        document.getElementById('authNameGroup').style.display = tab === 'register' ? 'block' : 'none';
        document.getElementById('authSubmitBtn').textContent = tab === 'register' ? 'Create Account' : 'Sign In';
    }

    async handleAuthSubmit(e) {
        e.preventDefault();
        const email = document.getElementById('authEmail').value;
        const password = document.getElementById('authPassword').value;
        const name = document.getElementById('authName')?.value;

        if (this.activeAuthTab === 'register') {
            const res = await this.api.register(name, email, password);
            if (res.success) {
                this.store.setState({ currentUser: res.user });
                this.closeAuthModal();
                this.showToast(`Welcome, ${res.user.name}!`, 'success');
                this.renderUserUI();
            } else {
                alert(res.error);
            }
        } else {
            const res = await this.api.login(email, password);
            if (res.success) {
                this.store.setState({ currentUser: res.user });
                this.closeAuthModal();
                this.showToast(`Logged in as ${res.user.email}`, 'success');
                this.renderUserUI();
            } else {
                alert(res.error);
            }
        }
    }

    openProfileModal() {
        const user = this.store.getState().currentUser;
        if (!user) {
            this.openAuthModal();
            return;
        }

        document.getElementById('profileNameInput').value = user.name || '';
        document.getElementById('profileEmailInput').value = user.email || '';
        document.getElementById('profilePhoneInput').value = user.phone || '';
        document.getElementById('profileAddressInput').value = user.address || '';
        document.getElementById('profileModal').classList.add('open');
    }

    closeProfileModal() {
        document.getElementById('profileModal').classList.remove('open');
    }

    async handleProfileSave(e) {
        e.preventDefault();
        const user = { ...this.store.getState().currentUser };
        user.name = document.getElementById('profileNameInput').value;
        user.phone = document.getElementById('profilePhoneInput').value;
        user.address = document.getElementById('profileAddressInput').value;

        await this.api.updateProfile(user);
        this.store.setState({ currentUser: user });
        this.closeProfileModal();
        this.showToast("Profile details updated successfully", "success");
        this.renderUserUI();
    }

    renderUserUI() {
        const user = this.store.getState().currentUser;
        const userDisplay = document.getElementById('userNameDisplay');
        const userAvatar = document.getElementById('userAvatarLetter');
        if (user && userDisplay) {
            userDisplay.textContent = user.name || user.email;
            userAvatar.textContent = (user.name || user.email).charAt(0).toUpperCase();
        }
    }

    updateConnectionUI(status) {
        const badge = document.getElementById('gatewayBadge');
        if (!badge) return;

        if (status.mode === 'docker') {
            badge.className = 'gateway-badge';
            badge.innerHTML = `<span class="gateway-badge-dot"></span><span>Docker Gateway (:8080)</span>`;
            badge.title = "Connected to Docker Compose. Click to switch to Kubernetes or Simulated.";
        } else if (status.mode === 'k8s') {
            badge.className = 'gateway-badge';
            badge.innerHTML = `<span class="gateway-badge-dot" style="background:#38bdf8; box-shadow:0 0 8px #38bdf8;"></span><span style="color:#38bdf8;">K8s Gateway (:8089)</span>`;
            badge.title = "Connected to Kubernetes Kind Cluster. Click to switch to Simulated.";
        } else {
            badge.className = 'gateway-badge simulated';
            badge.innerHTML = `<span class="gateway-badge-dot"></span><span>Simulated Event Bus</span>`;
            badge.title = "Running in Simulated mode. Click to switch to Docker (:8080).";
        }
    }

    toggleConnectionMode() {
        const current = this.api.mode;
        let next = 'docker';
        if (current === 'docker') next = 'k8s';
        else if (current === 'k8s') next = 'simulated';
        else next = 'docker';

        this.api.setMode(next);
        this.updateConnectionUI({ online: next !== 'simulated', mode: next });
        this.showToast(`Switched backend target to: ${next.toUpperCase()}`, 'info');
        this.loadRestaurants();
    }

    // --- Global Setup ---

    handleStateChange(state) {
        // Automatically sync UI on state updates
    }

    setupEventListeners() {
        // Search bar
        const searchInput = document.getElementById('searchInput');
        if (searchInput) {
            searchInput.addEventListener('input', e => {
                this.store.setState({ searchQuery: e.target.value });
                this.renderRestaurants();
            });
        }

        // Cuisine pills
        document.querySelectorAll('.filter-pill').forEach(pill => {
            pill.addEventListener('click', e => {
                document.querySelectorAll('.filter-pill').forEach(p => p.classList.remove('active'));
                pill.classList.add('active');
                const cuisine = pill.getAttribute('data-cuisine');
                this.store.setState({ selectedCuisine: cuisine });
                this.renderRestaurants();
            });
        });
    }

    renderAll() {
        this.renderRestaurants();
        this.renderCartUI();
        this.renderUserUI();
    }
}

window.app = new AppController();
document.addEventListener('DOMContentLoaded', () => window.app.init());
