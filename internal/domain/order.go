package domain

import "time"

type OrderItem struct {
	ProductID int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price,omitempty"`
	VAT       float64 `json:"vat,omitempty"`
}

type Order struct {
	ID        int64       `json:"order_id"`
	Items     []OrderItem `json:"items"`
	Price     float64     `json:"order_price,omitempty"`
	VAT       float64     `json:"order_vat,omitempty"`
	CreatedAt time.Time   `json:"created_at,omitempty"`
}

// Request and response structures
type CreateOrderRequest struct {
	Order struct {
		Items []OrderItem `json:"items"`
	} `json:"order"`
}

type OrderResponse struct {
	OrderID    int64       `json:"order_id"`
	OrderPrice float64     `json:"order_price"`
	OrderVAT   float64     `json:"order_vat"`
	Items      []OrderItem `json:"items"`
}
