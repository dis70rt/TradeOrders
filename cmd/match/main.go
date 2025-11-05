package main

import (
    log "github.com/dis70rt/TradeOrders/internals/logger"
    "github.com/dis70rt/TradeOrders/matching/engine"
)

func main() {
    log.Init()
    log.Info("Starting matching service...")
    engine.StartMatchingEngine()
}