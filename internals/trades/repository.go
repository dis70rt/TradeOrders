package trades

import "database/sql"

type Repository struct {
	db      *sql.DB
	queries *PrepareQuery
}

type PrepareQuery struct {
	insertTrade *sql.Stmt
}

func NewPrepareQuery(db *sql.DB) (*PrepareQuery, error) {
	insert, err := db.Prepare(`
		INSERT INTO trades (
			id, buy_order_id, sell_order_id, instrument, price, quantity, executed_at
		) VALUES ($1, $2, $3, $4, $5, $6, COALESCE($7, NOW()))
	`)
	if err != nil {
		return nil, err
	}

	return &PrepareQuery{
		insertTrade: insert,
	}, nil
}

func (r *Repository) InsertTrade(trade *Trade) error {
	_, err := r.queries.insertTrade.Exec(
		trade.ID, trade.BuyOrderID, trade.SellOrderID,
		trade.Instrument, trade.Price, trade.Quantity, trade.ExecutedAt,
	)
	return err
}
