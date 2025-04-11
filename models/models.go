package models

import "time"

type User struct {
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Phone     string    `json:"phone,omitempty"`
	Address   string    `json:"address,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Category struct {
	CategoryID  int64  `json:"category_id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type Plant struct {
	PlantID       int64     `json:"plant_id"`
	Name          string    `json:"name"`
	Description   string    `json:"description,omitempty"`
	CategoryID    int64     `json:"category_id"`
	ImageURL      string    `json:"image_url,omitempty"`
	Rating        float32    `json:"rating,omitempty"`
	Price         float64   `json:"price"`
	OriginalPrice float64   `json:"original_price"`
	BadgeText     string    `json:"badge_text,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type PlantInventory struct {
	InventoryID   int64     `json:"inventory_id"`
	PlantID       int64     `json:"plant_id"`
	Price         float64   `json:"price"`
	StockQuantity int       `json:"stock_quantity"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Service struct {
	ServiceID   int64     `json:"service_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
}

type Cart struct {
	CartID    int64     `json:"cart_id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CartItem struct {
	CartItemID int64     `json:"cart_item_id"`
	CartID     int64     `json:"cart_id"`
	PlantID    int64     `json:"plant_id,omitempty"`
	ServiceID  int64     `json:"service_id,omitempty"`
	Quantity   int       `json:"quantity"`
	Price      float64   `json:"price,omitempty"`
	Total      float64   `json:"total,omitempty"`
}

type Order struct {
	OrderID         int64     `json:"order_id"`
	UserID          int64     `json:"user_id"`
	TotalPrice      float64   `json:"total_price"`
	Status          string    `json:"status"`
	ShippingAddress string    `json:"shipping_address,omitempty"`
	OrderDate       time.Time `json:"order_date"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type OrderItem struct {
	OrderItemID int64     `json:"order_item_id"`
	OrderID     int64     `json:"order_id"`
	PlantID     int64     `json:"plant_id,omitempty"`
	ServiceID   int64     `json:"service_id,omitempty"`
	Quantity    int       `json:"quantity"`
	Price       float64   `json:"price,omitempty"`
	Total       float64   `json:"total,omitempty"`
}

type Payment struct {
	PaymentID     int64     `json:"payment_id"`
	OrderID       int64     `json:"order_id"`
	PaymentMethod string    `json:"payment_method,omitempty"`
	PaymentStatus string    `json:"payment_status"`
	PaymentDate   time.Time `json:"payment_date"`
	Amount        float64   `json:"amount"`
}

type Login struct {
	UserID 	  int64		`json:"user_id" db:"user"`
	Email     string    `json:"email" db:"email" binding:"required"`
	Password  string    `json:"password" db:"password" binding:"required"`
}