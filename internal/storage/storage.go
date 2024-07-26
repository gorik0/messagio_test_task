package storage

import (
	"database/sql"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"messagio/internal/config"
	"messagio/internal/service"
	"time"
)

type Storage struct {
	db       *sql.DB
	producer *kafka.Producer
	consumer *kafka.Consumer
}

var attempts int

func NewConn() *Storage {

	connStr := config.PostgresUrl()
	kafkaBrokers := config.KafkaBrokers()
	topic := []string{config.ConsumerTopicID(), config.ProducerTopicID()}
	groupId := config.KafkaGroup()

	producer := getProducer(kafkaBrokers)
	consumer := getConsumer(kafkaBrokers, topic, groupId)
	db := getDb(connStr)

	return &Storage{
		db:       db,
		producer: producer,
		consumer: consumer,
	}

}

func getConsumer(brokers string, topic []string, id string) *kafka.Consumer {

	attempts = 5
	for i := range attempts {
		cfg := kafka.ConfigMap{
			"bootstrap.servers": brokers,
			"group.id":          id,
			"auto.offset.reset": "earliest",
		}
		consumer, err := kafka.NewConsumer(&cfg)
		if err == nil {
			err = consumer.SubscribeTopics(topic, nil)
			if err == nil {
				return consumer
			}
		}

		log.Println("Failed to subscribe topic  :", err)
		log.Printf("Attemt %d has been expired, left %d ", i, attempts-i)
		time.Sleep(time.Second * 2)
	}
	panic("Failed to create consumer")
}

func getProducer(brokers string) *kafka.Producer {

	cfg := kafka.ConfigMap{"bootstrap.servers": brokers}
	producer, err := kafka.NewProducer(&cfg)
	if err != nil {
		panic("Failed to create producer: " + err.Error())

	}
	return producer

}

func getDb(strConn string) *sql.DB {

	attempts = 5
	for i := range attempts {
		conn, err := makeConn(strConn)
		if err == nil && conn != nil {
			return conn
		}
		log.Println("Failed to connect to db  :", err)
		log.Printf("Attemt %d has been expired, left %d ", i, attempts-i)
		time.Sleep(time.Second * 2)

	}
	panic("Failed to connect to db")

}

func makeConn(connStr string) (*sql.DB, error) {

	conn, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}
	err = conn.Ping()
	if err != nil {
		return nil, err
	}
	return conn, nil

}

var _ service.Storager = &Storage{}
