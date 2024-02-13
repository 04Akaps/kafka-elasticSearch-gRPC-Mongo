package db

import (
	"github.com/04Akaps/kafka-go/config"
	"log"
)

type DB struct {
	config *config.Config
}

type DBImpl interface {
}

func NewDB(config *config.Config) (DBImpl, error) {
	d := &DB{
		config: config,
	}

	log.Println("Success To Connect DB")

	return d, nil
}
