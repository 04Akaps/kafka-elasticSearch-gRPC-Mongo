package db

import "github.com/04Akaps/kafka-go/config"

type DB struct {
	config *config.Config
}

type DBImpl interface {
}

func NewDB(config *config.Config) (DBImpl, error) {
	d := &DB{
		config: config,
	}

	return d, nil
}
