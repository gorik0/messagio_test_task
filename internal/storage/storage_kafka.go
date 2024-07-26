package storage

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func (s *Storage) Produce(topic string, msg []byte) error {
	err := s.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: msg,
	}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Consume() (topic string, messageBytes []byte, err error) {
	for {
		message, err := s.consumer.ReadMessage(-1)
		if err != nil {
			return "", nil, fmt.Errorf("Error reading message: %v", err)
		} else {
			return *message.TopicPartition.Topic, message.Value, nil
		}

	}

}
