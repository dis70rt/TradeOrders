package orders

import (
	"context"
	"database/sql"

	log "github.com/dis70rt/TradeOrders/internals/logger"

	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
	queries *PrepareQuery
}

type PrepareQuery struct {
	insertOrder *sql.Stmt
}

func NewRepository(db *sql.DB) *Repository {
	queries, err := NewPrepareQuery(db)
	if err != nil {
		log.WithError(err).Fatal("Failed to prepare queries")
	}
	return &Repository{db: db, queries: queries}
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

	return &PrepareQuery{
		insertOrder: insert,
	}, nil
}

func (repo *Repository) CreateOrder(ctx context.Context, order *OrderRequest) (string, error) {
	orderID := uuid.NewString()

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

	return orderID, err
}