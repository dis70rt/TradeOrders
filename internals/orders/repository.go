package orders

import (
	"context"
	"database/sql"
	"encoding/json"

	log "github.com/dis70rt/TradeOrders/internals/logger"
	"github.com/dis70rt/TradeOrders/kafka"
	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
	queries *PrepareQuery
	producer *kafka.Producer
}

type PrepareQuery struct {
	insertOrder  *sql.Stmt
	getOrders  	 *sql.Stmt
	updateOrders *sql.Stmt
}

func NewRepository(db *sql.DB, producer *kafka.Producer) *Repository {
	queries, err := NewPrepareQuery(db)
	if err != nil {
		log.WithError(err).Fatal("Failed to prepare queries")
	}
	return &Repository{db: db, queries: queries, producer: producer}
}

func NewPrepareQuery(db *sql.DB) (*PrepareQuery, error) {
	insert, err := db.Prepare(`
		INSERT INTO orders (
			id, client_id, instrument, side, type, price, quantity, remaining
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`)
	if err != nil {
		return nil, err
	}

	getOrders, err := db.Prepare(`
		SELECT id, client_id, instrument, side, type, price, quantity, remaining, status, created_at, updated_at
		FROM orders
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2

	`)
	if err != nil {
		return nil, err
	}

	updateOrders, err := db.Prepare(`
		UPDATE orders
		SET 
			remaining = $2,
			status = $3,
			updated_at = NOW()
		WHERE id = $1
	`)

	return &PrepareQuery{
		insertOrder: insert,
		getOrders: 	getOrders,
		updateOrders: updateOrders,
	}, nil
}

func (repo *Repository) CreateOrder(ctx context.Context, order *OrderRequest) (string, error) {
	orderID := uuid.New()
	_, err := repo.queries.insertOrder.ExecContext(
		ctx,
		orderID,
		order.ClientID,
		order.Instrument,
		order.Side,
		order.Type,
		order.Price,
		order.Quantity,
		order.Remaining,
	)

	matchOrder := MatchOrder{
		ID: orderID,
		Instrument: order.Instrument,
		Side: order.Side,
		Type: order.Type,
		Price: order.Price,
		Quantity: order.Quantity,
	}

	orderJSON, _ := json.Marshal(matchOrder) 
	repo.producer.SendMessage("orders.inbound", order.Instrument, orderJSON)

	return orderID.String(), err
}

func (repo *Repository) GetOrders(ctx context.Context, limit, page int) ([]Order, error) {
	offset := (page - 1) * limit
	rows, err := repo.queries.getOrders.QueryContext(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]Order, 0, limit)
	for rows.Next() {
		var order Order
		if err := rows.Scan(
			&order.ID,
			&order.ClientID,
			&order.Instrument,
			&order.Side,
			&order.Type,
			&order.Price,
			&order.Quantity,
			&order.Remaining,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (repo *Repository) UpdateOrder(
	ctx context.Context,
	orderID uuid.UUID,
	remaining float64,
	status string,
) error {
	_, err := repo.queries.updateOrders.ExecContext(
		ctx,
		orderID,
		remaining,
		status,
	)
	return err
}
