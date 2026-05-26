package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	log "github.com/dis70rt/TradeOrders/internals/logger"
	"github.com/dis70rt/TradeOrders/internals/trades"
	"github.com/dis70rt/TradeOrders/kafka"
	"github.com/dis70rt/TradeOrders/streaming"
)

func main() {
    hub := streaming.NewHub()
    go hub.Run()

    handler := kafka.ConsumerHandler{
        Process: func(msg *sarama.ConsumerMessage) error {
            var t trades.Trade
            if err := json.Unmarshal(msg.Value, &t); err != nil {
                log.WithError(err).Error("stream: unmarshal trade")
                return err
            }
            hub.Broadcast(streaming.Broadcast{Instrument: t.Instrument, Data: msg.Value})
            return nil
        },
    }

    consumer := kafka.NewConsumer("TRADE_EXECUTED", "stream-trades", handler)
    go consumer.Start()
    defer consumer.Close()

    mux := http.NewServeMux()
    mux.HandleFunc("/ws/trades", func(w http.ResponseWriter, r *http.Request) {
        streaming.ServeWS(hub, w, r)
    })

    srv := &http.Server{
        Addr:              ":8081",
        Handler:           mux,
        ReadHeaderTimeout: 10 * time.Second,
    }

    go func() {
        log.Infof("WebSocket stream listening on %s (path: /ws/trades)", srv.Addr)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.WithError(err).Fatal("ws server error")
        }
    }()

    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
    <-sigc

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    _ = srv.Shutdown(ctx)
}