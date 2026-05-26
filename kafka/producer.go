package kafka

import (
	log "github.com/dis70rt/TradeOrders/internals/logger"
	"os"

	"github.com/IBM/sarama"
)

type Producer struct {
	async sarama.AsyncProducer
}

func NewProducer() *Producer {
	brokers := []string{os.Getenv("KAFKA_BROKER")}
	if brokers[0] == "" {
		brokers = []string{"localhost:9092"}
	}

	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Partitioner = sarama.NewHashPartitioner
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	cfg.Producer.Idempotent = true
	cfg.Net.MaxOpenRequests = 1

	p, err := sarama.NewAsyncProducer(brokers, cfg)
	if err != nil {
		log.WithError(err).Fatal("kafka producer init failed")
	}

	// Drain successes and errors to prevent channel backpressure
	go func() {
		for range p.Successes() {
		}
	}()
	go func() {
		for err := range p.Errors() {
			log.WithError(err).Error("kafka async produce error")
		}
	}()

	return &Producer{async: p}
}

func (p *Producer) SendMessage(topic string, key string, msg []byte) {
	p.async.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(msg),
	}
}

func (p *Producer) Close() {
	if err := p.async.Close(); err != nil {
		log.WithError(err).Error("kafka close error")
	}
}
