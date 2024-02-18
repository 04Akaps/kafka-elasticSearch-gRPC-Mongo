package repository

import (
	"github.com/04Akaps/kafka-go/config"
	"github.com/04Akaps/kafka-go/server/repository/db"
	"github.com/04Akaps/kafka-go/server/repository/elastic"
	"github.com/04Akaps/kafka-go/server/repository/kafka"
)

type Repository struct {
	config *config.Config

	Kafka   kafka.KafkaImpl
	Elastic elastic.ElasticImpl
	DB      db.DBImpl
}

func NewRepository(config *config.Config) (*Repository, error) {
	r := &Repository{config: config}
	var err error

	if r.Kafka, err = kafka.NewKafka(config); err != nil {
		panic(err)
	} else if r.DB, err = db.NewDB(config); err != nil {
		panic(err)
	} else if r.Elastic, err = elastic.NewElastic(config); err != nil {
		panic(err)
	} else {
		return r, nil
	}
}
