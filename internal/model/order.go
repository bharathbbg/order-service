package model

import (
	"time"
)

type Order struct {
	ID              string         `json:"id" db:"id"`
	CustomerID      string         `json:"customer_id" db:"customer_id"`
	Items           []OrderItem    `json:"items"`
	ShippingAddress Address        `json:"shipping_address"`
	Status          string         `json:"status" db:"status"`
	TotalPrice      Money          `json:"total_price"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
}

type OrderItem struct {
	ProductID   string `json:"product_id" db:"product_id"`
	Quantity    int    `json:"quantity" db:"quantity"`
	UnitPrice   Money  `json:"unit_price"`
	OrderID     string `json:"-" db:"order_id"`
}

type Money struct {
	Currency string `json:"currency" db:"currency"`
	Amount   int64  `json:"amount" db:"amount"`
}

type Address struct {
	Street  string `json:"street" db:"street"`
	City    string `json:"city" db:"city"`
	State   string `json:"state" db:"state"`
	Country string `json:"country" db:"country"`
	ZipCode string `json:"zip_code" db:"zip_code"`
}

// Request/Response models
type CreateOrderRequest struct {
	CustomerID      string      `json:"customer_id" binding:"required"`
	Items           []OrderItem `json:"items" binding:"required"`
	ShippingAddress Address     `json:"shipping_address" binding:"required"`
}

type UpdateOrderRequest struct {
	ID              string  `json:"-"`
	Status          string  `json:"status"`
	ShippingAddress Address `json:"shipping_address"`
}