package trades

import (
	"context"
	"database/sql"

	log "github.com/dis70rt/TradeOrders/internals/logger"
)

type Repository struct {
	db      *sql.DB
	queries *PrepareQuery
}

type PrepareQuery struct {
	insertTrade *sql.Stmt
    getTrades   *sql.Stmt
}

func NewRepository(db *sql.DB) *Repository {
	queries, err := newPrepareQuery(db)
	if err != nil {
		log.WithError(err).Fatal("Failed to prepare queries")
	}
	return &Repository{
		db: db,
		queries: queries,
	}
}

func newPrepareQuery(db *sql.DB) (*PrepareQuery, error) {
	insert, err := db.Prepare(`
		INSERT INTO trades (
			id, buy_order_id, sell_order_id, instrument, price, quantity, executed_at
		) VALUES ($1, $2, $3, $4, $5, $6, COALESCE($7, NOW()))
	`)
	if err != nil {
		return nil, err
	}

    get, err := db.Prepare(`
        SELECT id, buy_order_id, sell_order_id, instrument, price, quantity, executed_at
        FROM trades
        ORDER BY executed_at DESC
        LIMIT $1 OFFSET $2
    `)
    if err != nil {
        return nil, err
    }

	return &PrepareQuery{
		insertTrade: insert,
        getTrades:   get,
	}, nil
}

func (r *Repository) InsertTrade(trade *Trade) error {
	_, err := r.queries.insertTrade.Exec(
		trade.ID, trade.BuyOrderID, trade.SellOrderID,
		trade.Instrument, trade.Price, trade.Quantity, trade.ExecutedAt,
	)
	return err
}

func (r *Repository) ApplyTrade(
    ctx context.Context,
    trade *Trade,
) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    if _, err := tx.ExecContext(
        ctx,
        `UPDATE orders SET remaining = remaining - $1, updated_at=NOW()
         WHERE id=$2`,
        trade.Quantity, trade.BuyOrderID,
    ); err != nil {
        return err
    }

    if _, err := tx.ExecContext(
        ctx,
        `UPDATE orders SET remaining = remaining - $1, updated_at=NOW()
         WHERE id=$2`,
        trade.Quantity, trade.SellOrderID,
    ); err != nil {
        return err
    }

    if _, err := tx.ExecContext(
        ctx,
        `INSERT INTO trades (id, buy_order_id, sell_order_id, instrument, price, quantity, executed_at)
         VALUES ($1,$2,$3,$4,$5,$6,COALESCE($7,NOW()))`,
        trade.ID,
        trade.BuyOrderID,
        trade.SellOrderID,
        trade.Instrument,
        trade.Price,
        trade.Quantity,
        trade.ExecutedAt,
    ); err != nil {
        return err
    }

    return tx.Commit()
}

func (r *Repository) GetTrades(ctx context.Context, limit, offset int) ([]Trade, error) {
    rows, err := r.queries.getTrades.QueryContext(ctx, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    trades := make([]Trade, 0, limit)
    for rows.Next() {
        var trade Trade
        if err := rows.Scan(
            &trade.ID,
            &trade.BuyOrderID,
            &trade.SellOrderID,
            &trade.Instrument,
            &trade.Price,
            &trade.Quantity,
            &trade.ExecutedAt,
        ); err != nil {
            return nil, err
        }
        trades = append(trades, trade)
    }
    return trades, rows.Err()
}
