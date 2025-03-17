package repository

import (
	"context"

	"github.com/valeriouberti/order-service-test/internal/domain"
)

type ProductRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Product, error)
}

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) (*domain.Order, error)
	GetByID(ctx context.Context, id int64) (*domain.Order, error)
}
