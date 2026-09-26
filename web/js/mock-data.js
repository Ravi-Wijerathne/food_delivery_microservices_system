/**
 * Mock data matching exactly with Go Microservices MongoDB Seed Data:
 * - services/restaurant-service/data.go
 * - services/user-service/models.go
 * - services/delivery-service/main.go
 */

window.MOCK_DATA = {
    restaurants: [
        {
            id: "rest-1",
            name: "Pizza Palace",
            cuisine: "Italian",
            rating: 4.8,
            reviewsCount: 342,
            deliveryTime: "25-35 min",
            distance: "1.8 km",
            priceTier: "$$",
            minOrder: 15.00,
            address: "123 Main St, Little Italy",
            image: "assets/images/pizza.jpg",
            tag: "Wood-Fired Neapolitan",
            featured: true,
            badge: "Top Rated"
        },
        {
            id: "rest-2",
            name: "Burger Barn",
            cuisine: "American",
            rating: 4.6,
            reviewsCount: 289,
            deliveryTime: "20-30 min",
            distance: "2.4 km",
            priceTier: "$$",
            minOrder: 12.00,
            address: "456 Oak Ave, Downtown",
            image: "assets/images/burger.jpg",
            tag: "Gourmet Smash Burgers",
            featured: true,
            badge: "Staff Pick"
        },
        {
            id: "rest-3",
            name: "Sushi Supreme",
            cuisine: "Japanese",
            rating: 4.9,
            reviewsCount: 512,
            deliveryTime: "30-40 min",
            distance: "3.1 km",
            priceTier: "$$$",
            minOrder: 25.00,
            address: "789 Elm Blvd, Marina Bay",
            image: "assets/images/sushi.jpg",
            tag: "Wild-Caught Sashimi & Omakase",
            featured: true,
            badge: "Michelin Guide"
        },
        {
            id: "rest-4",
            name: "Taco Town",
            cuisine: "Mexican",
            rating: 4.5,
            reviewsCount: 198,
            deliveryTime: "15-25 min",
            distance: "1.2 km",
            priceTier: "$",
            minOrder: 10.00,
            address: "321 Pine Rd, Sunset District",
            image: "assets/images/tacos.jpg",
            tag: "Authentic Street Taqueria",
            featured: false,
            badge: "Fastest Delivery"
        },
        {
            id: "rest-5",
            name: "Curry Corner",
            cuisine: "Indian",
            rating: 4.7,
            reviewsCount: 420,
            deliveryTime: "30-45 min",
            distance: "2.8 km",
            priceTier: "$$",
            minOrder: 18.00,
            address: "654 Maple Dr, Spice Quarter",
            image: "assets/images/curry.jpg",
            tag: "Traditional Mughlai & Clay Oven",
            featured: true,
            badge: "Trending"
        }
    ],

    menus: {
        "rest-1": {
            restaurant_id: "rest-1",
            restaurant_name: "Pizza Palace",
            categories: ["Pizza", "Sides", "Beverages"],
            items: [
                {
                    id: "item-1",
                    name: "Margherita Pizza",
                    description: "Classic San Marzano tomato sauce, fior di latte mozzarella, fresh fragrant basil & virgin olive oil.",
                    price: 12.99,
                    category: "Pizza",
                    isVeg: true,
                    isSpicy: false,
                    popular: true
                },
                {
                    id: "item-2",
                    name: "Pepperoni Pizza",
                    description: "Hand-tossed crust topped with artisan spicy salami pepperoni and bubbling mozzarella blend.",
                    price: 14.99,
                    category: "Pizza",
                    isVeg: false,
                    isSpicy: true,
                    popular: true
                },
                {
                    id: "item-3",
                    name: "Garlic Bread",
                    description: "Crispy artisan ciabatta with roasted garlic butter, rosemary and parmesan shavings.",
                    price: 4.99,
                    category: "Sides",
                    isVeg: true,
                    isSpicy: false,
                    popular: false
                },
                {
                    id: "item-101",
                    name: "Truffle Mushroom Pizza",
                    description: "Wild foraged forest mushrooms, black truffle glaze, creamy fontina cheese and fresh thyme.",
                    price: 17.99,
                    category: "Pizza",
                    isVeg: true,
                    isSpicy: false,
                    popular: true
                },
                {
                    id: "item-102",
                    name: "San Pellegrino Sparkling Aranciata",
                    description: "Sparkling natural mineral water with Sicilian blood orange essence (330ml).",
                    price: 3.50,
                    category: "Beverages",
                    isVeg: true,
                    isSpicy: false,
                    popular: false
                }
            ]
        },
        "rest-2": {
            restaurant_id: "rest-2",
            restaurant_name: "Burger Barn",
            categories: ["Burgers", "Sides", "Beverages"],
            items: [
                {
                    id: "item-4",
                    name: "Classic Burger",
                    description: "100% Angus beef patty with crispy butter lettuce, vine tomato, house pickles and secret smash sauce.",
                    price: 9.99,
                    category: "Burgers",
                    isVeg: false,
                    isSpicy: false,
                    popular: false
                },
                {
                    id: "item-5",
                    name: "Cheese Burger",
                    description: "Double smashed prime patties layered with double molten aged Wisconsin cheddar and caramelized onions.",
                    price: 11.99,
                    category: "Burgers",
                    isVeg: false,
                    isSpicy: false,
                    popular: true
                },
                {
                    id: "item-6",
                    name: "French Fries",
                    description: "Crispy skin-on golden potato fries dusted with sea salt and smoked paprika.",
                    price: 3.99,
                    category: "Sides",
                    isVeg: true,
                    isSpicy: false,
                    popular: true
                },
                {
                    id: "item-201",
                    name: "Smoky Bacon BBQ Burger",
                    description: "Applewood smoked thick-cut bacon, crispy fried onion ring, tangy bourbon BBQ glaze and pepper jack.",
                    price: 13.99,
                    category: "Burgers",
                    isVeg: false,
                    isSpicy: true,
                    popular: true
                },
                {
                    id: "item-202",
                    name: "Craft Vanilla Milkshake",
                    description: "Thick hand-spun artisanal Madagascar vanilla bean ice cream with whipped cream.",
                    price: 5.99,
                    category: "Beverages",
                    isVeg: true,
                    isSpicy: false,
                    popular: false
                }
            ]
        },
        "rest-3": {
            restaurant_id: "rest-3",
            restaurant_name: "Sushi Supreme",
            categories: ["Rolls", "Sashimi", "Soup", "Sides"],
            items: [
                {
                    id: "item-7",
                    name: "Salmon Roll",
                    description: "Fresh Scottish salmon, ripe Hass avocado, Japanese cucumber and toasted sesame seeds.",
                    price: 15.99,
                    category: "Rolls",
                    isVeg: false,
                    isSpicy: false,
                    popular: true
                },
                {
                    id: "item-8",
                    name: "Tuna Sashimi",
                    description: "Five thick-cut slices of wild bluefin sashimi grade tuna served with fresh grated wasabi.",
                    price: 18.99,
                    category: "Sashimi",
                    isVeg: false,
                    isSpicy: false,
                    popular: true
                },
                {
                    id: "item-9",
                    name: "Miso Soup",
                    description: "Traditional red dashi broth with silken tofu cubes, wakame seaweed and spring scallions.",
                    price: 3.99,
                    category: "Soup",
                    isVeg: true,
                    isSpicy: false,
                    popular: false
                },
                {
                    id: "item-301",
                    name: "Dragon Special Roll",
                    description: "Crispy tempura prawn, BBQ unagi eel, avocado slices, tobiko roe and sweet tare glaze.",
                    price: 19.99,
                    category: "Rolls",
                    isVeg: false,
                    isSpicy: false,
                    popular: true
                },
                {
                    id: "item-302",
                    name: "Spicy Edamame",
                    description: "Steamed young soybean pods tossed in garlic chili oil and Maldon sea salt flakes.",
                    price: 5.49,
                    category: "Sides",
                    isVeg: true,
                    isSpicy: true,
                    popular: false
                }
            ]
        },
        "rest-4": {
            restaurant_id: "rest-4",
            restaurant_name: "Taco Town",
            categories: ["Tacos", "Burritos", "Sides"],
            items: [
                {
                    id: "item-10",
                    name: "Beef Tacos",
                    description: "Trio of grilled corn tortillas loaded with marinated grilled carne asada, cilantro, onion and cotija cheese.",
                    price: 10.99,
                    category: "Tacos",
                    isVeg: false,
                    isSpicy: true,
                    popular: true
                },
                {
                    id: "item-11",
                    name: "Chicken Burrito",
                    description: "Jumbo flour tortilla stuffed with citrus-grilled chicken, cilantro-lime rice, black beans and pico de gallo.",
                    price: 12.99,
                    category: "Burritos",
                    isVeg: false,
                    isSpicy: false,
                    popular: true
                },
                {
                    id: "item-12",
                    name: "Guacamole & Chips",
                    description: "Freshly crushed Hass avocados with lime, serrano pepper and warm crispy house tortilla chips.",
                    price: 6.99,
                    category: "Sides",
                    isVeg: true,
                    isSpicy: false,
                    popular: true
                },
                {
                    id: "item-401",
                    name: "Baja Crispy Fish Tacos",
                    description: "Beer-battered cod fish, chipotle crema, shredded red cabbage and pickled red onions.",
                    price: 11.99,
                    category: "Tacos",
                    isVeg: false,
                    isSpicy: true,
                    popular: false
                }
            ]
        },
        "rest-5": {
            restaurant_id: "rest-5",
            restaurant_name: "Curry Corner",
            categories: ["Curry", "Vegetarian", "Breads", "Sides"],
            items: [
                {
                    id: "item-13",
                    name: "Chicken Curry",
                    description: "Tender tandoor roasted chicken breast simmered in a velvety buttery tomato fenugreek gravy.",
                    price: 13.99,
                    category: "Curry",
                    isVeg: false,
                    isSpicy: true,
                    popular: true
                },
                {
                    id: "item-14",
                    name: "Palak Paneer",
                    description: "Fresh artisanal cottage cheese cubes gently cooked in seasoned spiced pureed spinach and garlic.",
                    price: 11.99,
                    category: "Vegetarian",
                    isVeg: true,
                    isSpicy: false,
                    popular: true
                },
                {
                    id: "item-15",
                    name: "Garlic Naan",
                    description: "Puffy leavened flatbread baked in clay tandoor oven, brushed with melted ghee and fresh minced garlic.",
                    price: 2.99,
                    category: "Breads",
                    isVeg: true,
                    isSpicy: false,
                    popular: true
                },
                {
                    id: "item-501",
                    name: "Lamb Rogan Josh",
                    description: "Slow-braised New Zealand lamb shoulder infused with Kashmiri chili, cardamom and aromatic spices.",
                    price: 16.99,
                    category: "Curry",
                    isVeg: false,
                    isSpicy: true,
                    popular: true
                },
                {
                    id: "item-502",
                    name: "Saffron Basmati Rice",
                    description: "Long grain fragrant aged basmati rice infused with Persian saffron strands and whole cloves.",
                    price: 4.49,
                    category: "Sides",
                    isVeg: true,
                    isSpicy: false,
                    popular: false
                }
            ]
        }
    },

    drivers: [
        { name: "Sarah Jenkins", vehicle: "Honda PCX 160 (Silver)", rating: 4.95, trips: 1420, phone: "+1 (555) 234-8901", avatar: "👩🏻" },
        { name: "John Miller", vehicle: "Yamaha NMAX (Matte Black)", rating: 4.88, trips: 980, phone: "+1 (555) 345-6712", avatar: "👨🏼" },
        { name: "Mike Chen", vehicle: "Vespa Primavera (Red)", rating: 4.92, trips: 1150, phone: "+1 (555) 456-7823", avatar: "👨🏻" },
        { name: "David Kim", vehicle: "Toyota Prius (Hybrid)", rating: 4.90, trips: 840, phone: "+1 (555) 567-8934", avatar: "👨🏻" },
        { name: "Emma Watson", vehicle: "Electric Super Soco Bike", rating: 4.97, trips: 2100, phone: "+1 (555) 678-9045", avatar: "👩🏼" }
    ],

    defaultUser: {
        id: "user-1",
        email: "test@example.com",
        name: "Alex Vance",
        phone: "+1 (555) 987-6543",
        address: "742 Evergreen Terrace, Apt 4B, Springfield",
        notes: "Leave package at door. Ring smart bell."
    },

    microservices: [
        { id: "gateway", name: "API Gateway", port: 8080, protocol: "HTTP / REST", status: "ONLINE", icon: "🌐", role: "Routing, Circuit Breaking, JWT Auth & Prometheus Metrics" },
        { id: "auth-service", name: "Auth Service", port: 8081, protocol: "HTTP / JWT", status: "ONLINE", icon: "🔐", role: "User Credentials, Hashing & JWT Token Generation" },
        { id: "restaurant-service", name: "Restaurant Service", port: 8082, grpcPort: 9082, protocol: "REST & gRPC", status: "ONLINE", icon: "🍽️", role: "Catalogs, Menus, gRPC Menu Validator" },
        { id: "order-service", name: "Order Service", port: 8083, grpcPort: 8084, protocol: "REST & gRPC", status: "ONLINE", icon: "📦", role: "Order Lifecycle Orchestrator & RabbitMQ Publisher" },
        { id: "payment-service", name: "Payment Service", port: 8085, protocol: "HTTP & AMQP", status: "ONLINE", icon: "💳", role: "Asynchronous Payment Processing via RabbitMQ" },
        { id: "delivery-service", name: "Delivery Service", port: 8086, protocol: "HTTP & AMQP", status: "ONLINE", icon: "🛵", role: "Driver Assignment & Geo Fleet Dispatching" },
        { id: "notification-service", name: "Notification Service", port: 8087, protocol: "AMQP Consumer", status: "ONLINE", icon: "🔔", role: "Event Subscriber & Customer SMS/Push Alerts" },
        { id: "user-service", name: "User Service", port: 8088, grpcPort: 9088, protocol: "REST & gRPC", status: "ONLINE", icon: "👤", role: "User Profiles & gRPC Account Verification" },
        { id: "rabbitmq", name: "RabbitMQ Message Bus", port: 5672, mgmtPort: 15672, protocol: "AMQP 0-9-1", status: "ONLINE", icon: "🐇", role: "Event-Driven Async Messaging Engine" },
        { id: "mongodb", name: "MongoDB Cluster", port: 27017, protocol: "BSON / WiredTiger", status: "ONLINE", icon: "🍃", role: "Isolated Database-per-Service Architecture" }
    ]
};
