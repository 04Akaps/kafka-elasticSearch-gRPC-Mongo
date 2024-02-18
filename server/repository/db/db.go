package db

import (
	"context"
	"github.com/04Akaps/kafka-go/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type DB struct {
	config *config.Config

	client *mongo.Client
	db     *mongo.Database

	block *mongo.Collection
	tx    *mongo.Collection
}

type DBImpl interface {
}

func NewDB(config *config.Config) (DBImpl, error) {
	d := &DB{config: config}

	var err error
	ctx := context.Background()

	if d.client, err = mongo.Connect(ctx, options.Client().ApplyURI(config.DB.URI)); err != nil {
		return nil, err
	} else if err = d.client.Ping(ctx, nil); err != nil {
		return nil, err
	} else {
		d.db = d.client.Database(config.DB.DB)
		// TODO Collection

		log.Println("Success To Connect DB")
		return d, nil
	}
}
