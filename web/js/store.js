/**
 * Global Reactive State Store for CravePulse Food Delivery UI
 */

class Store {
    constructor() {
        this.state = {
            currentUser: null,
            token: localStorage.getItem('crave_token') || '',
            connectionMode: 'simulated', // 'gateway' | 'simulated'
            restaurants: [],
            selectedCuisine: 'all',
            searchQuery: '',
            sortBy: 'rating',
            activeRestaurant: null,
            currentMenu: null,
            cart: this.loadCart(),
            activeOrder: this.loadActiveOrder(),
            orderHistory: this.loadOrderHistory(),
            eventLogs: [],
            systemHealth: {},
            activeTab: 'restaurants', // 'restaurants' | 'orders' | 'observability'
            isCartOpen: false,
            isAuthModalOpen: false,
            isProfileModalOpen: false,
            isMenuModalOpen: false
        };

        this.listeners = [];
    }

    subscribe(callback) {
        this.listeners.push(callback);
        return () => {
            this.listeners = this.listeners.filter(cb => cb !== callback);
        };
    }

    setState(partialState) {
        this.state = { ...this.state, ...partialState };
        this.listeners.forEach(cb => cb(this.state));
    }

    getState() {
        return this.state;
    }

    // --- Cart Management ---

    loadCart() {
        try {
            return JSON.parse(localStorage.getItem('crave_cart') || '{"items":[], "restaurantId": null, "restaurantName": ""}');
        } catch (e) {
            return { items: [], restaurantId: null, restaurantName: "" };
        }
    }

    saveCart(cart) {
        localStorage.setItem('crave_cart', JSON.stringify(cart));
        this.setState({ cart });
    }

    addToCart(restaurant, item, quantity = 1) {
        const cart = { ...this.state.cart };

        // Check if cart has items from another restaurant
        if (cart.restaurantId && cart.restaurantId !== restaurant.id && cart.items.length > 0) {
            const confirmSwitch = confirm(
                `Your cart contains items from "${cart.restaurantName}". Would you like to clear your cart and start a new order from "${restaurant.name}"?`
            );
            if (!confirmSwitch) return false;
            cart.items = [];
        }

        cart.restaurantId = restaurant.id;
        cart.restaurantName = restaurant.name;

        const existingIdx = cart.items.findIndex(i => i.id === item.id);
        if (existingIdx >= 0) {
            cart.items[existingIdx].quantity += quantity;
        } else {
            cart.items.push({
                id: item.id,
                name: item.name,
                price: item.price,
                quantity: quantity,
                isVeg: item.isVeg,
                category: item.category
            });
        }

        this.saveCart(cart);
        return true;
    }

    updateCartItemQuantity(itemId, delta) {
        const cart = { ...this.state.cart };
        const idx = cart.items.findIndex(i => i.id === itemId);
        if (idx === -1) return;

        cart.items[idx].quantity += delta;
        if (cart.items[idx].quantity <= 0) {
            cart.items.splice(idx, 1);
        }

        if (cart.items.length === 0) {
            cart.restaurantId = null;
            cart.restaurantName = "";
        }

        this.saveCart(cart);
    }

    removeCartItem(itemId) {
        const cart = { ...this.state.cart };
        cart.items = cart.items.filter(i => i.id !== itemId);
        if (cart.items.length === 0) {
            cart.restaurantId = null;
            cart.restaurantName = "";
        }
        this.saveCart(cart);
    }

    clearCart() {
        this.saveCart({ items: [], restaurantId: null, restaurantName: "" });
    }

    getCartTotals() {
        const items = this.state.cart.items || [];
        const subtotal = items.reduce((sum, i) => sum + (i.price * i.quantity), 0);
        const deliveryFee = subtotal > 35 || subtotal === 0 ? 0.00 : 3.99;
        const serviceFee = subtotal > 0 ? 1.49 : 0.00;
        const tax = subtotal * 0.0825; // 8.25%
        const total = subtotal + deliveryFee + serviceFee + tax;
        const itemCount = items.reduce((sum, i) => sum + i.quantity, 0);

        return {
            subtotal,
            deliveryFee,
            serviceFee,
            tax,
            total,
            itemCount
        };
    }

    // --- Orders Management ---

    loadActiveOrder() {
        try {
            return JSON.parse(localStorage.getItem('crave_active_order') || 'null');
        } catch(e) {
            return null;
        }
    }

    setActiveOrder(order) {
        if (order) {
            localStorage.setItem('crave_active_order', JSON.stringify(order));
            this.addOrderToHistory(order);
        } else {
            localStorage.removeItem('crave_active_order');
        }
        this.setState({ activeOrder: order });
    }

    loadOrderHistory() {
        try {
            return JSON.parse(localStorage.getItem('crave_order_history') || '[]');
        } catch(e) {
            return [];
        }
    }

    addOrderToHistory(order) {
        const history = this.loadOrderHistory();
        const existingIdx = history.findIndex(o => o.id === order.id);
        if (existingIdx >= 0) {
            history[existingIdx] = order;
        } else {
            history.unshift(order);
        }
        localStorage.setItem('crave_order_history', JSON.stringify(history));
        this.setState({ orderHistory: history });
    }
}

window.store = new Store();
