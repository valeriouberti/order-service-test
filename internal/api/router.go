package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/valeriouberti/order-service-test/internal/api/handlers"
)

func NewRouter(orderHandler *handlers.OrderHandler) *mux.Router {
	r := mux.NewRouter()

	// Define API routes
	r.HandleFunc("/api/orders", orderHandler.CreateOrder).Methods("POST")
	r.HandleFunc("/api/orders/{id}", orderHandler.GetOrder).Methods("GET")

	// Add health check endpoint
	r.HandleFunc("/health-check", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	return r
}
