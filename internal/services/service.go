package services

import (
	"context"

	"github.com/valeriouberti/order-service-test/internal/domain"
)

type OrderServiceInterface interface {
	CreateOrder(ctx context.Context, req *domain.CreateOrderRequest) (*domain.OrderResponse, error)
	GetOrder(ctx context.Context, id int64) (*domain.OrderResponse, error)
}
