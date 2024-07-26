package service

import (
	"encoding/json"
	"fmt"
	"log"
	"messagio/internal/config"
	"messagio/internal/models"
)

func (s *Service) HandleMessage(msg models.JsonMsg) error {
	queueName := config.ProducerTopicID()

	if msg.Data == nil {
		return fmt.Errorf("Data can not be nil")
	}

	dataString, ok := msg.Data.(string)
	if !ok {
		return fmt.Errorf("Data is not a string")
	}
	id, err := s.St.AddRecord(queueName, dataString)
	if err != nil {
		return fmt.Errorf("Error adding record: %s", err)
	}

	kafkaMsg := models.KafkaMesg{
		ID:      int(id),
		Message: dataString,
	}

	kafkaMsgBytes, _ := json.Marshal(kafkaMsg)
	err = s.St.Produce(queueName, kafkaMsgBytes)
	if err != nil {
		return fmt.Errorf("Error producing message: %s", err)
	}
	return nil

}

func (s *Service) GiveMeStats() (int, int, error) {
	total, processed, err := s.St.NumberOfProcessedMessages()
	if err != nil {
		return -1, -1, err
	}
	return total, processed, nil
}

func (s *Service) StartConsumeMessages() {
	go s.consumeMessages()
}
func (s *Service) CloseDB() error {
	err := s.St.CloseDB()
	if err != nil {
		return err

	}
	return nil
}

func (s *Service) consumeMessages() {

	for {
		topic, mesgBytes, err := s.St.Consume()
		if err != nil {
			log.Println("Error consuming messages:", err)
			continue
		}
		switch topic {
		case config.ConsumerTopicID():
			{
				var mesJSON models.KafkaMesg
				err := json.Unmarshal(mesgBytes, &mesJSON)
				if err != nil {
					log.Printf("Error unmarshaling JSON: %s", err)
					continue

				}
				err = s.St.MessageConsumed(int64(mesJSON.ID))
				if err != nil {
					log.Printf("Error consuming message: %s", err)
					continue

				}
			}
		case config.ProducerTopicID():
			{

				err := s.St.Produce(config.ConsumerTopicID(), mesgBytes)
				if err != nil {
					log.Printf("Error producing message: %s", err)
					continue
				}

			}

		}

		//	:::: PRODUCE mesaage or mark as proccessed depends on  TOPIC name
		//::: TODO(){for further understsanding}

	}
}
