package services

import (
	"context"
	"fmt"

	"github.com/valeriouberti/order-service-test/internal/domain"
	"github.com/valeriouberti/order-service-test/internal/repository"
)

type OrderService struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
}

func NewOrderService(orderRepo repository.OrderRepository, productRepo repository.ProductRepository) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

// CreateOrder creates a new order based on the provided request.
//
// It performs the following steps:
// 1. Initializes an order with items from the request
// 2. For each item:
//   - Retrieves the product details from repository
//   - Calculates the price and VAT based on product information and quantity
//   - Updates the item with calculated values
//
// 3. Calculates total price and VAT for the entire order
// 4. Persists the order in the database
// 5. Maps the created order to a response object
//
// The function returns the order response containing ID, price, VAT, and items.
// If a product is not found or if there's an error saving the order, an error is returned.
//
// Parameters:
//   - ctx: context.Context for the operation
//   - req: *domain.CreateOrderRequest containing order details
//
// Returns:
//   - *domain.OrderResponse: the created order with calculated totals
//   - error: if any step fails during order creation
func (s *OrderService) CreateOrder(ctx context.Context, req *domain.CreateOrderRequest) (*domain.OrderResponse, error) {
	// Initialize order
	order := &domain.Order{
		Items: req.Order.Items,
	}

	// Calculate price and VAT for each item
	var totalPrice, totalVAT float64

	for i, item := range order.Items {
		// Get product details
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product with ID %d not found: %w", item.ProductID, err)
		}

		// Calculate item price and VAT
		itemPrice := product.Price * float64(item.Quantity)
		itemVAT := product.VAT * float64(item.Quantity)

		// Update the item with price and VAT
		order.Items[i].Price = itemPrice
		order.Items[i].VAT = itemVAT

		// Add to totals
		totalPrice += itemPrice
		totalVAT += itemVAT
	}

	// Set order totals
	order.Price = totalPrice
	order.VAT = totalVAT

	// Save order to database
	createdOrder, err := s.orderRepo.Create(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Map to response
	response := &domain.OrderResponse{
		OrderID:    createdOrder.ID,
		OrderPrice: createdOrder.Price,
		OrderVAT:   createdOrder.VAT,
		Items:      createdOrder.Items,
	}

	return response, nil
}

// GetOrder retrieves an order by its ID.
// It fetches the order from the repository and maps it to an OrderResponse type.
//
// Parameters:
//   - ctx: The context for the operation, which can be used for cancellation, timeouts, etc.
//   - id: The unique identifier of the order to retrieve.
//
// Returns:
//   - *domain.OrderResponse: The order data formatted as a response object, or nil if an error occurs.
//   - error: Any error encountered during the operation. If the order is not found, the repository error is returned.
func (s *OrderService) GetOrder(ctx context.Context, id int64) (*domain.OrderResponse, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := &domain.OrderResponse{
		OrderID:    order.ID,
		OrderPrice: order.Price,
		OrderVAT:   order.VAT,
		Items:      order.Items,
	}

	return response, nil
}
