package service

import (
	"encoding/json"
	"github.com/04Akaps/kafka-go/config"
	"github.com/04Akaps/kafka-go/server/repository"
	"github.com/04Akaps/kafka-go/server/types"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"log"
	"strings"
	"time"
)

const likeTopic = "like-topic"

var topics map[string]string

type Service struct {
	config *config.Config

	repository *repository.Repository
}

type ServiceImpl interface {
	Like(fromUser, toUser string, point int64) error
	UnLike(fromUser, toUser string, point int64) error
}

func NewService(config *config.Config, repository *repository.Repository) ServiceImpl {
	s := &Service{
		config:     config,
		repository: repository,
	}

	return s
}

func (s *Service) Like(fromUser, toUser string, point int64) error {
	if err := s.repository.DB.Like(toUser, point); err != nil {
		log.Println("Failed To Like Request", err)
		return err
	} else {
		go s.sendLikeEventToKafka(fromUser, toUser, point)
		return nil
	}
}

func (s *Service) UnLike(fromUser, toUser string, point int64) error {
	if err := s.repository.DB.UnLike(toUser, point); err != nil {
		log.Println("Failed To UnLike Request", err)
		return err
	} else {
		go s.sendLikeEventToKafka(fromUser, toUser, point)
		return nil
	}
}

func (s *Service) sendLikeEventToKafka(fromUser, toUser string, point int64) {
	var action string

	if point == 0 {
		action = "zero"
	} else if point < 0 {
		action = "minus"
	} else {
		action = "plus"
	}

	event := types.KafkaEvent{
		Index:     likeTopic,
		ElasticId: uuid.New().String(),
		Data: types.LikeHistory{
			FromUser:   fromUser,
			ToUser:     toUser,
			Point:      point,
			Action:     action,
			SearchText: strings.Join([]string{fromUser, toUser, action}, " "),
			Time:       time.Now().Unix(),
		},
	}

	if value, err := json.Marshal(event); err != nil {
		log.Println("Failed To Marshal Like History")
	} else {
		ch := make(chan kafka.Event)

		if result, err := s.repository.Kafka.SendEvent(likeTopic, value, ch); err != nil {
			log.Println("Failed To Send Event")
		} else {
			log.Println("Success to send to kafka", result.String())
		}
	}
}
