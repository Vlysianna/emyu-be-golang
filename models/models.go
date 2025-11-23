package models

import "time"

// Role
type Role struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	Permissions []string  `json:"permissions,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RolePermission
type RolePermission struct {
	ID         int       `json:"id"`
	RoleID     int       `json:"role_id"`
	Permission string    `json:"permission"`
	CreatedAt  time.Time `json:"created_at"`
}

// User
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	RoleID    int       `json:"role_id"`
	Role      *Role     `json:"role,omitempty"`
	Password  string    `json:"-"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Category
type Category struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Product
type Product struct {
	ID             string           `json:"id"`
	Name           string           `json:"name"`
	Description    string           `json:"description"`
	Price          float64          `json:"price"`
	CategoryID     string           `json:"category_id"`
	IsCustomizable bool             `json:"is_customizable"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	Category       *Category        `json:"category,omitempty"`
	Images         []ProductImage   `json:"images,omitempty"`
	Variants       []ProductVariant `json:"variants,omitempty"`
}

// ProductImage
type ProductImage struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
}

// ProductVariant
type ProductVariant struct {
	ID              string    `json:"id"`
	ProductID       string    `json:"product_id"`
	Name            string    `json:"name"`
	PriceAdjustment float64   `json:"price_adjustment"`
	CreatedAt       time.Time `json:"created_at"`
}

// Cart
type Cart struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Items     []CartItem `json:"items,omitempty"`
}

// CartItem
type CartItem struct {
	ID               string          `json:"id"`
	CartID           string          `json:"cart_id"`
	ProductVariantID string          `json:"product_variant_id"`
	Quantity         int             `json:"quantity"`
	CustomName       string          `json:"custom_name"`
	CustomNumber     string          `json:"custom_number"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	ProductVariant   *ProductVariant `json:"product_variant,omitempty"`
}

// ShippingAddress
type ShippingAddress struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Address    string    `json:"address"`
	City       string    `json:"city"`
	Province   string    `json:"province"`
	PostalCode string    `json:"postal_code"`
	Phone      string    `json:"phone"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Order
type Order struct {
	ID                string           `json:"id"`
	UserID            string           `json:"user_id"`
	OrderNumber       string           `json:"order_number"`
	TotalAmount       float64          `json:"total_amount"`
	ShippingCost      float64          `json:"shipping_cost"`
	Status            string           `json:"status"`         // pending, paid, packed, shipped, delivered, canceled
	PaymentMethod     string           `json:"payment_method"` // qris, bank_transfer, ewallet
	ShippingAddressID string           `json:"shipping_address_id"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
	User              *User            `json:"user,omitempty"`
	ShippingAddress   *ShippingAddress `json:"shipping_address,omitempty"`
	Items             []OrderItem      `json:"items,omitempty"`
	Payment           *Payment         `json:"payment,omitempty"`
}

// OrderItem
type OrderItem struct {
	ID               string    `json:"id"`
	OrderID          string    `json:"order_id"`
	ProductVariantID string    `json:"product_variant_id"`
	Quantity         int       `json:"quantity"`
	Price            float64   `json:"price"`
	CreatedAt        time.Time `json:"created_at"`
}

// Payment
type Payment struct {
	ID            string     `json:"id"`
	OrderID       string     `json:"order_id"`
	PaymentStatus string     `json:"payment_status"` // pending, success, failed
	PaymentCode   string     `json:"payment_code"`
	PaidAt        *time.Time `json:"paid_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// Review
type Review struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ProductID string    `json:"product_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	User      *User     `json:"user,omitempty"`
}

// Auth DTOs
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}
