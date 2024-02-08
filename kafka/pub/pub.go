package pub

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

type Kafka struct {
	producer    *kafka.Producer
	consumerMap map[string]*kafka.Consumer

	topicMap map[string]string
}

type KafkaImpl interface {
	SendEvent(topic string, value []byte, ch chan kafka.Event) (kafka.Event, error)
	SendEventWithKey(topic string, key string, ch chan kafka.Event, value []byte) (kafka.Event, error)
	AddNewConsumer(topic string, conf *kafka.ConfigMap) error
}

func NewKafka(url, clientId string) (KafkaImpl, error) {
	producerConf := &kafka.ConfigMap{
		"bootstrap.servers": url,
		"client.id":         clientId,
		"acks":              "all",
	}

	k := Kafka{
		topicMap:    make(map[string]string),
		consumerMap: make(map[string]*kafka.Consumer),
	}
	var err error

	if k.producer, err = kafka.NewProducer(producerConf); err != nil {
		return nil, err
	} else {
		return &k, nil
	}
}

func (k *Kafka) SendEvent(topic string, value []byte, ch chan kafka.Event) (kafka.Event, error) {
	// TODO
	// 추가된 토픽 없이 생성이 된다면, 어떻게 처리되는지 테스트 필요
	// 테스트 이후 에러 케이스로 추가
	//if _, ok := k.topicMap[topic]; !ok {
	//	return nil, errors.New("Topic is Not Added")
	//}

	if err := k.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: value,
	}, ch); err != nil {
		return nil, err
	} else {
		return <-ch, nil
	}
}

func (k *Kafka) SendEventWithKey(topic string, key string, ch chan kafka.Event, value []byte) (kafka.Event, error) {
	if err := k.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: value,
		Key:   []byte(key),
	},
		ch,
	); err != nil {
		return nil, err
	} else {
		return <-ch, nil
	}
}

func (k *Kafka) AddNewConsumer(topic string, conf *kafka.ConfigMap) error {
	if err := k.checkTopicExisted(topic); err != nil {
		return err
	} else if consumer, err := kafka.NewConsumer(conf); err != nil {
		return err
	} else if err = consumer.Subscribe(topic, nil); err != nil {
		return err
	} else {
		k.consumerMap[topic] = consumer
		go k.subscribeEvent(topic, consumer)
		return nil
	}
}

func (k *Kafka) checkTopicExisted(topic string) error {
	if _, ok := k.topicMap[topic]; !ok {
		return errors.New("topic is not existed")
	} else {
		return nil
	}
}

func (k *Kafka) subscribeEvent(topic string, consumer *kafka.Consumer) error {
	ctx := context.Background()

	withCancel, done := context.WithCancel(ctx)

	for {
		select {
		case <-withCancel.Done():
			delete(k.topicMap, topic)
			delete(k.consumerMap, topic)
		default:
			ev := consumer.Poll(100)

			switch event := ev.(type) {
			case *kafka.Message:
				var consumeValue interface{}

				event.Timestamp.Unix()

				if err := json.Unmarshal(event.Value, &consumeValue); err != nil {
					log.Println("Failed To Decode", err.Error())
				} else {

					// TODO Commit Message
				}

			case *kafka.Error:
				log.Println("Polling Event Err", event.Error())
				done()
			}
		}

	}
}
