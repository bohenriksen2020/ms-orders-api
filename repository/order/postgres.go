package order

import (
	"context"
	"database/sql"
	"fmt"

	// Import PostgreSQL driver for database/sql
	_ "github.com/lib/pq"

	"github.com/bohenriksen2020/ms-orders-api/model"
)

// PostgresRepo implements the Repo interface for PostgreSQL
type PostgresRepo struct {
	DB *sql.DB
}

// NewPostgresRepo initializes the PostgresRepo with a given database connection
func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{DB: db}
}

// Close closes the database connection
func (repo *PostgresRepo) Close() error {
	return repo.DB.Close()
}

// Insert inserts a new order into the PostgreSQL database
func (repo *PostgresRepo) Insert(ctx context.Context, ord model.Order) error {
	query := `
		INSERT INTO orders (order_id, customer_id, line_items, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := repo.DB.ExecContext(ctx, query, ord.OrderID, ord.CustomerID, ord.LineItems, ord.Created)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}
	return nil
}

// FindByID retrieves an order by its ID
func (repo *PostgresRepo) FindByID(ctx context.Context, id uint64) (model.Order, error) {
	query := `
		SELECT order_id, customer_id, line_items, created_at, shipped_at, completed_at
		FROM orders
		WHERE order_id = $1
	`

	var ord model.Order
	err := repo.DB.QueryRowContext(ctx, query, id).
		Scan(&ord.OrderID,
			&ord.CustomerID,
			&ord.LineItems,
			&ord.Created,
			&ord.ShippedAt,
			&ord.CompletedAt)
	if err == sql.ErrNoRows {
		return ord, ErrNotExist
	} else if err != nil {
		return ord, fmt.Errorf("failed to retrieve order: %w", err)
	}
	return ord, nil
}

// Update updates an existing order in the PostgreSQL database
func (repo *PostgresRepo) Update(ctx context.Context, ord model.Order) error {
	query := `
		UPDATE orders
		SET customer_id = $1, line_items = $2, shipped_at = $3, completed_at = $4
		WHERE order_id = $5
	`
	_, err := repo.DB.ExecContext(ctx, query, ord.CustomerID, ord.LineItems, ord.ShippedAt, ord.CompletedAt, ord.OrderID)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}
	return nil
}

// DeleteByID deletes an order by its ID
func (repo *PostgresRepo) DeleteByID(ctx context.Context, id uint64) error {
	query := `DELETE FROM orders WHERE order_id = $1`
	_, err := repo.DB.ExecContext(ctx, query, id)
	if err == sql.ErrNoRows {
		return ErrNotExist
	} else if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}
	return nil
}

// FindAll retrieves a paginated list of orders
func (repo *PostgresRepo) FindAll(ctx context.Context, page FindAllPage) (FindResult, error) {
	query := `
		SELECT order_id, customer_id, line_items, created_at, shipped_at, completed_at
		FROM orders
		ORDER BY order_id ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := repo.DB.QueryContext(ctx, query, page.Size, page.Offset)
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to find orders: %w", err)
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var ord model.Order
		if err := rows.Scan(&ord.OrderID, &ord.CustomerID, &ord.LineItems, &ord.Created, &ord.ShippedAt, &ord.CompletedAt); err != nil {
			return FindResult{}, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, ord)
	}

	if err := rows.Err(); err != nil {
		return FindResult{}, fmt.Errorf("row iteration error: %w", err)
	}

	return FindResult{
		Orders: orders,
		Cursor: page.Offset + uint64(len(orders)),
	}, nil
}
