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

	like        *mongo.Collection
	likeHistory *mongo.Collection
}

type DBImpl interface {
	Like(toUser string, point int64) error
	UnLike(toUser string, point int64) error
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

		d.like = d.db.Collection("like")
		d.likeHistory = d.db.Collection("like-history")

		log.Println("Success To Connect DB")
		return d, nil
	}
}
