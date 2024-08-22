package deribit

import (
	"context"
	"log"

	"github.com/IBM/sarama"
)

const (
	consumerGroupId = "orderbook-consumer-group"
)

type OrderbookConsumer struct {
	*Deribit
}

func NewOrderbookConsumer(d *Deribit) *OrderbookConsumer {
	return &OrderbookConsumer{
		Deribit: d,
	}
}

func (*OrderbookConsumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (*OrderbookConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (o OrderbookConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("key = %s, timestamp = %v orderbook = %s\n",
			msg.Key, string(msg.Value), msg.Timestamp)
		session.MarkMessage(msg, "")
	}
	return nil
}

func (o *OrderbookConsumer) StartConsuming(ctx context.Context, validCurrencies, validInstrumentKinds []string) {
	log.Println("Starting orderbook consumers")
	topics := []string{}
	for _, instrumentKind := range validInstrumentKinds {
		for _, currency := range validCurrencies {
			topics = append(topics, createOrderbookTopic(currency, instrumentKind))
		}
	}

	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	consumerConfig.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	consumerGroup, err := sarama.NewConsumerGroup([]string{o.config.KAFKA_SERVER_ADDRESS}, consumerGroupId, consumerConfig)
	if err != nil {
		log.Fatalf("Error creating consumer group: %v", err)
	}
	defer func() {
		if err := consumerGroup.Close(); err != nil {
			log.Fatalf("Error closing consumer group: %v", err)
		}
	}()
	for {
		select {
		case <-ctx.Done():
			log.Println("Context done exiting start consuming function")
			return
		default:
			err = consumerGroup.Consume(ctx, topics, o)
			if err != nil {
				log.Printf("Error from consumer: %v\n", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}
}
