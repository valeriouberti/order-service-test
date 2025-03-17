package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/valeriouberti/order-service-test/internal/domain"
)

type OrderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

// Create persists a new order record in the database along with its items.
// It uses a transaction to ensure atomicity - either all records are saved or none are.
//
// The method:
// 1. Inserts the order record and retrieves its generated ID and creation timestamp
// 2. Inserts all associated order items using the newly generated order ID
// 3. Commits the transaction if everything succeeds
//
// Parameters:
//   - ctx: The context for database operations, allows for cancellation and timeouts
//   - order: A pointer to the domain.Order to be persisted
//
// Returns:
//   - A pointer to the domain.Order with ID and CreatedAt populated from the database
//   - An error if any database operation fails
//
// The method will roll back the transaction on any error.
func (r *OrderRepo) Create(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insert the order
	query := `
        INSERT INTO orders (price, vat, created_at)
        VALUES ($1, $2, NOW())
        RETURNING id, created_at
    `

	err = tx.QueryRowContext(
		ctx,
		query,
		order.Price,
		order.VAT,
	).Scan(&order.ID, &order.CreatedAt)

	if err != nil {
		return nil, err
	}

	// Insert order items
	for i, item := range order.Items {
		query = `
            INSERT INTO order_items (order_id, product_id, quantity, price, vat)
            VALUES ($1, $2, $3, $4, $5)
        `

		_, err = tx.ExecContext(
			ctx,
			query,
			order.ID,
			item.ProductID,
			item.Quantity,
			item.Price,
			item.VAT,
		)

		if err != nil {
			return nil, err
		}

		// Update the item in our order with the database values
		order.Items[i] = item
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return order, nil
}

// GetByID retrieves an order by its ID from the database.
// It executes a SQL query that fetches the order details along with its associated order items.
// The order items are aggregated into a JSON array using PostgreSQL's json_agg function.
//
// Parameters:
//   - ctx: Context for database operations, allowing for cancellation and timeouts
//   - id: The unique identifier of the order to retrieve
//
// Returns:
//   - *domain.Order: A pointer to the retrieved order with its items populated
//   - error: An error if the order is not found or if there's a database or JSON unmarshaling error
//
// Errors:
//   - Returns "order not found" error if no order exists with the given ID
//   - Returns unmarshaling errors if the JSON data for items is malformed
func (r *OrderRepo) GetByID(ctx context.Context, id int64) (*domain.Order, error) {
	query := `
        SELECT o.id, o.price, o.vat, o.created_at, 
               COALESCE(json_agg(
                   json_build_object(
                       'product_id', oi.product_id,
                       'quantity', oi.quantity,
                       'price', oi.price,
                       'vat', oi.vat
                   )
               ) FILTER (WHERE oi.id IS NOT NULL), '[]') as items
        FROM orders o
        LEFT JOIN order_items oi ON o.id = oi.order_id
        WHERE o.id = $1
        GROUP BY o.id
    `

	var order domain.Order
	var itemsJSON string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&order.ID,
		&order.Price,
		&order.VAT,
		&order.CreatedAt,
		&itemsJSON,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	// Parse items JSON
	if err = json.Unmarshal([]byte(itemsJSON), &order.Items); err != nil {
		return nil, err
	}

	return &order, nil
}
