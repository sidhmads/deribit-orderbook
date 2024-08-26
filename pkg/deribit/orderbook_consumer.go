package deribit

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

type OrderbookConsumer struct {
	*Deribit
}

func NewOrderbookConsumer(d *Deribit) *OrderbookConsumer {
	return &OrderbookConsumer{
		Deribit: d,
	}
}

func (o *OrderbookConsumer) StartConsuming(ctx context.Context) error {
	log.Println("Starting orderbook consumers")

	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Return.Errors = true

	client, err := sarama.NewClient([]string{o.config.KAFKA_SERVER_ADDRESS}, consumerConfig)
	if err != nil {
		return fmt.Errorf("Failed to create Kafka client: %v", err)
	}
	defer client.Close()

	partitions, err := client.Partitions(orderbookTopic)
	if err != nil {
		return fmt.Errorf("Failed to get partitions for topic %s: %v", orderbookTopic, err)
	}

	consumerCtx, cancelFn := context.WithCancel(ctx)
	for _, partition := range partitions {
		go o.StartConsumingPerPartition(consumerCtx, partition, consumerConfig)
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("Context done exiting start consuming function")
			cancelFn()
			log.Println("Shutting down consumers")
			return nil
		}
	}
}

func (o *OrderbookConsumer) StartConsumingPerPartition(ctx context.Context, partition int32, config *sarama.Config) {
	consumer, err := sarama.NewConsumer([]string{o.config.KAFKA_SERVER_ADDRESS}, config)
	if err != nil {
		log.Fatalf("Failed to start consumer for partition %d: %v", partition, err)
	}
	defer consumer.Close()

	log.Printf("Starting to consume on partition: %d", partition)

	partitionConsumer, err := consumer.ConsumePartition(orderbookTopic, partition, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to consume partition %d: %v", partition, err)
	}
	defer partitionConsumer.Close()

	for {
		select {
		case message := <-partitionConsumer.Messages():
			if message != nil {
				log.Printf("key = %s, orderbook = %s,  timestamp = %v, partition: %d, offset: %d \n",
					message.Key, message.Value, message.Timestamp, message.Partition, message.Offset)
			}
		case err := <-partitionConsumer.Errors():
			if err != nil {
				log.Printf("Error consuming partition %d: %v", partition, err)
			}
		case <-ctx.Done():
			log.Printf("Context done, exiting consumer for partition %d", partition)
			return
		}
	}
}
