package main

type Restaurant struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Cuisine string  `json:"cuisine"`
	Rating  float64 `json:"rating"`
	Address string  `json:"address"`
}

type MenuItem struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
}

type Menu struct {
	RestaurantID string     `json:"restaurant_id"`
	Restaurant   string     `json:"restaurant_name"`
	Items        []MenuItem `json:"items"`
}

var restaurants = []Restaurant{
	{ID: "rest-1", Name: "Pizza Palace", Cuisine: "Italian", Rating: 4.5, Address: "123 Main St"},
	{ID: "rest-2", Name: "Burger Barn", Cuisine: "American", Rating: 4.2, Address: "456 Oak Ave"},
	{ID: "rest-3", Name: "Sushi Supreme", Cuisine: "Japanese", Rating: 4.8, Address: "789 Elm Blvd"},
	{ID: "rest-4", Name: "Taco Town", Cuisine: "Mexican", Rating: 4.3, Address: "321 Pine Rd"},
	{ID: "rest-5", Name: "Curry Corner", Cuisine: "Indian", Rating: 4.6, Address: "654 Maple Dr"},
}

var menus = map[string][]MenuItem{
	"rest-1": {
		{ID: "item-1", Name: "Margherita Pizza", Description: "Classic tomato and mozzarella", Price: 12.99, Category: "Pizza"},
		{ID: "item-2", Name: "Pepperoni Pizza", Description: "With pepperoni slices", Price: 14.99, Category: "Pizza"},
		{ID: "item-3", Name: "Garlic Bread", Description: "Crispy with garlic butter", Price: 4.99, Category: "Sides"},
	},
	"rest-2": {
		{ID: "item-4", Name: "Classic Burger", Description: "Beef patty with lettuce and tomato", Price: 9.99, Category: "Burgers"},
		{ID: "item-5", Name: "Cheese Burger", Description: "With melted cheddar", Price: 11.99, Category: "Burgers"},
		{ID: "item-6", Name: "French Fries", Description: "Crispy golden fries", Price: 3.99, Category: "Sides"},
	},
	"rest-3": {
		{ID: "item-7", Name: "Salmon Roll", Description: "Fresh salmon with avocado", Price: 15.99, Category: "Rolls"},
		{ID: "item-8", Name: "Tuna Sashimi", Description: "Premium tuna slices", Price: 18.99, Category: "Sashimi"},
		{ID: "item-9", Name: "Miso Soup", Description: "Traditional miso", Price: 3.99, Category: "Soup"},
	},
	"rest-4": {
		{ID: "item-10", Name: "Beef Tacos", Description: "Three tacos with seasoned beef", Price: 10.99, Category: "Tacos"},
		{ID: "item-11", Name: "Chicken Burrito", Description: "Large flour tortilla wrap", Price: 12.99, Category: "Burritos"},
		{ID: "item-12", Name: "Guacamole & Chips", Description: "Fresh guacamole", Price: 6.99, Category: "Sides"},
	},
	"rest-5": {
		{ID: "item-13", Name: "Chicken Curry", Description: "Creamy butter chicken", Price: 13.99, Category: "Curry"},
		{ID: "item-14", Name: "Palak Paneer", Description: "Spinach and cheese curry", Price: 11.99, Category: "Vegetarian"},
		{ID: "item-15", Name: "Garlic Naan", Description: "Soft bread with garlic", Price: 2.99, Category: "Breads"},
	},
}

func InitData() {}

func GetRestaurants() []Restaurant {
	return restaurants
}

func GetMenu(restaurantID string) (*Menu, error) {
	items, exists := menus[restaurantID]
	if !exists {
		return nil, ErrMenuNotFound
	}

	restaurantName := ""
	for _, r := range restaurants {
		if r.ID == restaurantID {
			restaurantName = r.Name
			break
		}
	}

	return &Menu{
		RestaurantID: restaurantID,
		Restaurant:   restaurantName,
		Items:        items,
	}, nil
}

var ErrMenuNotFound = &RestaurantError{"menu not found"}

type RestaurantError struct {
	Message string
}

func (e *RestaurantError) Error() string {
	return e.Message
}
