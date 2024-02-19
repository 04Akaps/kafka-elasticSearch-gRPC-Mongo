package repository

import (
	"fmt"
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

	elasticLog := make(chan interface{}, 100)

	if r.Kafka, err = kafka.NewKafka(config, elasticLog); err != nil {
		panic(err)
	} else if r.DB, err = db.NewDB(config); err != nil {
		panic(err)
	} else if r.Elastic, err = elastic.NewElastic(config, elasticLog); err != nil {
		panic(err)
	} else {

		go func() {
			for {
				test := <-elasticLog
				fmt.Println("testset", test)
			}
		}()
		return r, nil
	}
}
