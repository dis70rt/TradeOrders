package kafka

import (
	"context"
	"os"
	log "github.com/dis70rt/TradeOrders/internals/logger"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

type Consumer struct {
	group   sarama.ConsumerGroup
	topic   string
	handler ConsumerHandler
}

type ConsumerHandler struct {
	Process func(message *sarama.ConsumerMessage) error
}

func (ConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (ConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if h.Process != nil {
			if err := h.Process(msg); err != nil {
				log.WithError(err).Errorf(
					"failed to process message topic=%s partition=%d offset=%d",
					msg.Topic, msg.Partition, msg.Offset,
				)
				continue
			}
		}
		session.MarkMessage(msg, "")
	}
	return nil
}

func NewConsumer(topic, groupID string, handler ConsumerHandler) *Consumer {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "localhost:9092"
	}

	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_5_0_0
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	cfg.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()

	group, err := sarama.NewConsumerGroup([]string{broker}, groupID, cfg)
	if err != nil {
		log.WithError(err).Fatal("kafka consumer init failed")
	}

	return &Consumer{
		group:   group,
		topic:   topic,
		handler: handler,
	}
}

func (c *Consumer) Start() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	for {
		if err := c.group.Consume(ctx, []string{c.topic}, c.handler); err != nil {
			log.WithError(err).Error("kafka consume error")
		}
		if ctx.Err() != nil {
			return
		}
	}
}

func (c *Consumer) Close() {
	if err := c.group.Close(); err != nil {
		log.WithError(err).Error("kafka consumer close error")
	}
}
