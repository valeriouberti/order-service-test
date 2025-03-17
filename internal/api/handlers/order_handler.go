package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/valeriouberti/order-service-test/internal/domain"
	"github.com/valeriouberti/order-service-test/internal/services"
)

type OrderHandler struct {
	orderService services.OrderServiceInterface
}

func NewOrderHandler(orderService services.OrderServiceInterface) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// CreateOrder handles HTTP POST requests to create a new order.
// It validates the incoming JSON payload to ensure it contains at least one item,
// processes the order through the order service, and returns the created order details.
//
// The handler expects a request body containing a JSON representation of domain.CreateOrderRequest.
// Upon successful creation, it returns a 201 Created status with the order details.
// In case of errors, it returns appropriate HTTP error codes:
// - 400 Bad Request: For invalid JSON or orders with no items
// - 500 Internal Server Error: For errors during order processing
//
// @param w http.ResponseWriter - The response writer to write the HTTP response
// @param r *http.Request - The HTTP request containing the order details in the body
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if len(req.Order.Items) == 0 {
		http.Error(w, "Order must contain at least one item", http.StatusBadRequest)
		return
	}

	// Process the order
	response, err := h.orderService.CreateOrder(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetOrder handles HTTP GET requests to retrieve order details by ID.
// It extracts the order ID from the URL path parameters, validates it,
// and calls the order service to fetch the requested order.
//
// If the order ID is invalid, it returns a 400 Bad Request response.
// If the order is not found or another error occurs during retrieval,
// it returns a 404 Not Found response with the error message.
// On success, it returns a 200 OK response with the order details as JSON.
//
// The response format is determined by the order service implementation.
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	// Parse order ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	// Get the order
	response, err := h.orderService.GetOrder(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
