package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valeriouberti/order-service-test/internal/domain"
)

// Mock repositories
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	args := m.Called(ctx, order)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *MockOrderRepository) GetByID(ctx context.Context, id int64) (*domain.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

// Test CreateOrder
func TestCreateOrder(t *testing.T) {
	// Setup
	mockOrderRepo := new(MockOrderRepository)
	mockProductRepo := new(MockProductRepository)
	orderService := NewOrderService(mockOrderRepo, mockProductRepo)

	ctx := context.Background()

	// Test case 1: Successful order creation
	t.Run("Successful order creation", func(t *testing.T) {
		// Mock input
		req := &domain.CreateOrderRequest{
			Order: struct {
				Items []domain.OrderItem `json:"items"`
			}{
				Items: []domain.OrderItem{
					{ProductID: 1, Quantity: 2},
					{ProductID: 2, Quantity: 3},
				},
			},
		}

		// Mock expected products
		product1 := &domain.Product{ID: 1, Name: "Product 1", Price: 10.0, VAT: 1.0}
		product2 := &domain.Product{ID: 2, Name: "Product 2", Price: 5.0, VAT: 0.5}

		// Mock expected order to be returned from repository
		expectedOrder := &domain.Order{
			ID: 1,
			Items: []domain.OrderItem{
				{ProductID: 1, Quantity: 2, Price: 20.0, VAT: 2.0},
				{ProductID: 2, Quantity: 3, Price: 15.0, VAT: 1.5},
			},
			Price: 35.0,
			VAT:   3.5,
		}

		// Set up mocks
		mockProductRepo.On("GetByID", ctx, int64(1)).Return(product1, nil)
		mockProductRepo.On("GetByID", ctx, int64(2)).Return(product2, nil)
		mockOrderRepo.On("Create", ctx, mock.AnythingOfType("*domain.Order")).Return(expectedOrder, nil)

		// Call the service
		result, err := orderService.CreateOrder(ctx, req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.OrderID)
		assert.Equal(t, 35.0, result.OrderPrice)
		assert.Equal(t, 3.5, result.OrderVAT)
		assert.Len(t, result.Items, 2)

		// Verify mocks
		mockProductRepo.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t)
	})

	// Test case 2: Product not found
	t.Run("Product not found", func(t *testing.T) {
		// Mock input
		req := &domain.CreateOrderRequest{
			Order: struct {
				Items []domain.OrderItem `json:"items"`
			}{
				Items: []domain.OrderItem{
					{ProductID: 999, Quantity: 1},
				},
			},
		}

		// Set up mocks
		mockProductRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("product not found"))

		// Call the service
		result, err := orderService.CreateOrder(ctx, req)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "product not found")

		// Verify mocks
		mockProductRepo.AssertExpectations(t)
	})
}

// Test GetOrder
func TestGetOrder(t *testing.T) {
	// Setup
	mockOrderRepo := new(MockOrderRepository)
	mockProductRepo := new(MockProductRepository)
	orderService := NewOrderService(mockOrderRepo, mockProductRepo)

	ctx := context.Background()

	// Test case 1: Successful order retrieval
	t.Run("Successful order retrieval", func(t *testing.T) {
		// Mock expected order
		expectedOrder := &domain.Order{
			ID: 1,
			Items: []domain.OrderItem{
				{ProductID: 1, Quantity: 2, Price: 20.0, VAT: 2.0},
				{ProductID: 2, Quantity: 3, Price: 15.0, VAT: 1.5},
			},
			Price: 35.0,
			VAT:   3.5,
		}

		// Set up mocks
		mockOrderRepo.On("GetByID", ctx, int64(1)).Return(expectedOrder, nil)

		// Call the service
		result, err := orderService.GetOrder(ctx, 1)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.OrderID)
		assert.Equal(t, 35.0, result.OrderPrice)
		assert.Equal(t, 3.5, result.OrderVAT)
		assert.Len(t, result.Items, 2)

		// Verify mocks
		mockOrderRepo.AssertExpectations(t)
	})

	// Test case 2: Order not found
	t.Run("Order not found", func(t *testing.T) {
		// Set up mocks
		mockOrderRepo.On("GetByID", ctx, int64(999)).Return(nil, errors.New("order not found"))

		// Call the service
		result, err := orderService.GetOrder(ctx, 999)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "order not found")

		// Verify mocks
		mockOrderRepo.AssertExpectations(t)
	})
}
