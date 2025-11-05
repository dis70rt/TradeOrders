-- Active: 1762063885808@@127.0.0.1@5432@trade_db
CREATE TABLE orders (
    id UUID PRIMARY KEY,
    client_id VARCHAR(64) NOT NULL,
    instrument VARCHAR(16) NOT NULL,
    side VARCHAR(4) CHECK (side IN ('BUY', 'SELL')) NOT NULL,
    type VARCHAR(10) CHECK (type IN ('LIMIT', 'MARKET')) NOT NULL,
    price NUMERIC(18,8) NOT NULL CHECK (price > 0),
    quantity NUMERIC(18,8) NOT NULL CHECK (quantity > 0),
    remaining NUMERIC(18,8) NOT NULL CHECK (remaining >= 0),
    status VARCHAR(16) DEFAULT 'NEW' CHECK (status IN ('NEW', 'PARTIALLY_FILLED', 'FILLED', 'CANCELLED')),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE trades (
    id UUID PRIMARY KEY,
    buy_order_id UUID REFERENCES orders(id),
    sell_order_id UUID REFERENCES orders(id),
    instrument VARCHAR(16) NOT NULL,
    price NUMERIC(18,8) NOT NULL,
    quantity NUMERIC(18,8) NOT NULL,
    executed_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE order_queue (
    id UUID PRIMARY KEY,
    order_id UUID REFERENCES orders(id),
    status VARCHAR(16) CHECK (status IN ('QUEUED', 'PROCESSING', 'DONE')) DEFAULT 'QUEUED',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_orders_instrument_side_price ON orders (instrument, side, price DESC);
CREATE INDEX idx_orders_status ON orders (status);
CREATE INDEX idx_trades_buy_order ON trades (buy_order_id);
CREATE INDEX idx_trades_sell_order ON trades (sell_order_id);
CREATE INDEX idx_order_queue_status ON order_queue (status);
CREATE INDEX idx_order_queue_order_id ON order_queue (order_id);

-- DROP TABLE orders CASCADE;
-- DROP TABLE trades CASCADE;
-- DROP TABLE order_queue CASCADE;

CREATE OR REPLACE FUNCTION set_order_status()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.remaining = NEW.quantity THEN
        NEW.status := 'NEW';
    ELSIF NEW.remaining = 0 THEN
        NEW.status := 'FILLED';
    ELSE
        NEW.status := 'PARTIALLY_FILLED';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_order_status
BEFORE INSERT OR UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION set_order_status();
