package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/04Akaps/kafka-go/config"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

type Kafka struct {
	producer    *kafka.Producer
	consumerMap map[string]*kafka.Consumer

	elasticLog chan<- interface{}

	topicMap map[string]bool
}

type KafkaImpl interface {
	SendEvent(topic string, value []byte, ch chan kafka.Event) (kafka.Event, error)
	SendEventWithKey(topic string, key string, ch chan kafka.Event, value []byte) (kafka.Event, error)
	AddNewConsumer(topic string, conf *kafka.ConfigMap) error
}

func NewKafka(config *config.Config, elasticLog chan interface{}) (KafkaImpl, error) {
	producerConf := &kafka.ConfigMap{
		"bootstrap.servers": config.Kafka.URI,
		"client.id":         config.Kafka.ClientId,
		"acks":              "all",
	}

	k := Kafka{
		topicMap:    make(map[string]bool),
		consumerMap: make(map[string]*kafka.Consumer),
		elasticLog:  elasticLog,
	}
	var err error

	if k.producer, err = kafka.NewProducer(producerConf); err != nil {
		return nil, err
	} else {

		for _, topic := range config.Kafka.Topics {
			k.topicMap[topic] = true

			consumerConf := &kafka.ConfigMap{
				"bootstrap.servers": config.Kafka.URI,
				"group.id":          "consumer_group_3",
				"auto.offset.reset": "latest",
				// "go.application.rebalance.enable": true,
				"enable.auto.commit": false,
				// "partition.assignment.strategy": "roundrobin",
			}

			k.AddNewConsumer(topic, consumerConf)
		}

		return &k, nil
	}
}

func (k *Kafka) SendEvent(topic string, value []byte, ch chan kafka.Event) (kafka.Event, error) {
	// 토픽이 생성이 되어 있지 않으면, 이벤트를 subscribe 할 수 없으니.
	// 토픽 체크
	if _, ok := k.topicMap[topic]; !ok {
		return nil, errors.New("Topic is Not Added")
	} else if err := k.producer.Produce(&kafka.Message{
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
					k.elasticLog <- consumeValue
					if partion, err := k.consumerMap[topic].CommitMessage(event); err != nil {
						log.Println(partion)
					}
				}

			case *kafka.Error:
				log.Println("Polling Event Err", event.Error())
				done()
			}
		}

	}
}
