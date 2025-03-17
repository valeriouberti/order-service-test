package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valeriouberti/order-service-test/internal/domain"
)

// MockOrderService is a mock implementation of the OrderServicer interface
type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) CreateOrder(ctx context.Context, req *domain.CreateOrderRequest) (*domain.OrderResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.OrderResponse), args.Error(1)
}

func (m *MockOrderService) GetOrder(ctx context.Context, id int64) (*domain.OrderResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.OrderResponse), args.Error(1)
}

func TestCreateOrder(t *testing.T) {
	// Setup
	mockService := new(MockOrderService)
	handler := NewOrderHandler(mockService)

	t.Run("Successful order creation", func(t *testing.T) {
		// Reset mock before each test.  Good practice!
		mockService.ExpectedCalls = nil
		mockService.Calls = nil

		// Create request body
		reqBody := domain.CreateOrderRequest{
			Order: struct {
				Items []domain.OrderItem `json:"items"`
			}{
				Items: []domain.OrderItem{
					{ProductID: 1, Quantity: 2},
					{ProductID: 2, Quantity: 3},
				},
			},
		}

		// Marshal request to JSON
		body, _ := json.Marshal(reqBody)

		// Create HTTP request
		req := httptest.NewRequest("POST", "/orders", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		w := httptest.NewRecorder()

		// Set up mock response
		expectedResp := &domain.OrderResponse{
			OrderID:    1,
			OrderPrice: 35.0,
			OrderVAT:   3.5,
			Items: []domain.OrderItem{
				{ProductID: 1, Quantity: 2, Price: 20.0, VAT: 2.0},
				{ProductID: 2, Quantity: 3, Price: 15.0, VAT: 1.5},
			},
		}

		// Set up mock expectation
		mockService.On("CreateOrder", mock.Anything, mock.MatchedBy(func(req *domain.CreateOrderRequest) bool {
			return len(req.Order.Items) == 2
		})).Return(expectedResp, nil)

		// Call handler
		handler.CreateOrder(w, req)

		// Assertions
		assert.Equal(t, http.StatusCreated, w.Code)

		var response domain.OrderResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, int64(1), response.OrderID)
		assert.Equal(t, 35.0, response.OrderPrice)
		assert.Equal(t, 3.5, response.OrderVAT)
		assert.Len(t, response.Items, 2)

		// Verify mock
		mockService.AssertExpectations(t)
	})

	t.Run("Empty order items", func(t *testing.T) {
		// Reset mock
		mockService.ExpectedCalls = nil
		mockService.Calls = nil
		// Create request body with no items
		reqBody := domain.CreateOrderRequest{
			Order: struct {
				Items []domain.OrderItem `json:"items"`
			}{
				Items: []domain.OrderItem{},
			},
		}

		// Marshal request to JSON
		body, _ := json.Marshal(reqBody)

		// Create HTTP request
		req := httptest.NewRequest("POST", "/orders", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		w := httptest.NewRecorder()

		// Call handler
		handler.CreateOrder(w, req)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Order must contain at least one item")
	})

	t.Run("Invalid JSON request", func(t *testing.T) {
		// Reset mock
		mockService.ExpectedCalls = nil
		mockService.Calls = nil

		// Create invalid JSON
		invalidJSON := []byte(`{"order":{"items":[{"product_id":1,}]}`)

		// Create HTTP request
		req := httptest.NewRequest("POST", "/orders", bytes.NewBuffer(invalidJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		w := httptest.NewRecorder()

		// Call handler
		handler.CreateOrder(w, req)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request body")
	})

	t.Run("Service error", func(t *testing.T) {
		// Reset mock
		mockService.ExpectedCalls = nil
		mockService.Calls = nil

		// Create request body
		reqBody := domain.CreateOrderRequest{
			Order: struct {
				Items []domain.OrderItem `json:"items"`
			}{
				Items: []domain.OrderItem{
					{ProductID: 1, Quantity: 2},
				},
			},
		}

		// Marshal request to JSON
		body, _ := json.Marshal(reqBody)

		// Create HTTP request
		req := httptest.NewRequest("POST", "/orders", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		w := httptest.NewRecorder()

		// Set up mock to return error
		mockService.On("CreateOrder", mock.Anything, mock.MatchedBy(func(req *domain.CreateOrderRequest) bool {
			return len(req.Order.Items) == 1
		})).Return(nil, errors.New("service error"))

		// Call handler
		handler.CreateOrder(w, req)

		// Assertions
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "service error")

		// Verify mock
		mockService.AssertExpectations(t)
	})
}

func TestGetOrder(t *testing.T) {
	// Setup
	mockService := new(MockOrderService)
	handler := NewOrderHandler(mockService)

	t.Run("Successful order retrieval", func(t *testing.T) {
		// Reset mock
		mockService.ExpectedCalls = nil
		mockService.Calls = nil

		// Create HTTP request
		req := httptest.NewRequest("GET", "/orders/1", nil)
		// *** KEY CHANGE:  Use mux.SetURLVars to simulate route variables ***
		req = mux.SetURLVars(req, map[string]string{"id": "1"})

		// Create response recorder
		w := httptest.NewRecorder()

		// Set up mock response
		expectedResp := &domain.OrderResponse{
			OrderID:    1,
			OrderPrice: 35.0,
			OrderVAT:   3.5,
			Items: []domain.OrderItem{
				{ProductID: 1, Quantity: 2, Price: 20.0, VAT: 2.0},
				{ProductID: 2, Quantity: 3, Price: 15.0, VAT: 1.5},
			},
		}

		// Set up mock expectation
		mockService.On("GetOrder", mock.Anything, int64(1)).Return(expectedResp, nil)

		// Call handler
		handler.GetOrder(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		var response domain.OrderResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, int64(1), response.OrderID)
		assert.Equal(t, 35.0, response.OrderPrice)
		assert.Equal(t, 3.5, response.OrderVAT)
		assert.Len(t, response.Items, 2)

		// Verify mock
		mockService.AssertExpectations(t)
	})

	t.Run("Order not found", func(t *testing.T) {
		// Reset mock
		mockService.ExpectedCalls = nil
		mockService.Calls = nil
		// Create HTTP request
		req := httptest.NewRequest("GET", "/orders/999", nil)
		// *** KEY CHANGE: Use mux.SetURLVars to simulate route variables ***
		req = mux.SetURLVars(req, map[string]string{"id": "999"})

		// Create response recorder
		w := httptest.NewRecorder()

		// Set up mock to return error
		mockService.On("GetOrder", mock.Anything, int64(999)).Return(nil, errors.New("order not found"))

		// Call handler
		handler.GetOrder(w, req)

		// Assertions
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "order not found")

		// Verify mock
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid order ID", func(t *testing.T) {
		// Reset mock
		mockService.ExpectedCalls = nil
		mockService.Calls = nil

		// Create HTTP request with an invalid ID (non-numeric)
		req := httptest.NewRequest("GET", "/orders/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"}) // Set the invalid ID

		// Create response recorder
		w := httptest.NewRecorder()

		// Call handler (no mock setup needed, as the error occurs before service call)
		handler.GetOrder(w, req)

		// Assertions
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid order ID")
	})
}
