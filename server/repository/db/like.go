package db

import (
	"github.com/04Akaps/kafka-go/server/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d *DB) Like(toUser string, point int64) error {
	filter := bson.M{"user": toUser}
	update := bson.M{"$inc": bson.M{"point": point}}

	opt := options.Update().SetUpsert(true)

	_, err := d.like.UpdateOne(util.Context(), filter, update, opt)
	return err
}

func (d *DB) UnLike(toUser string, point int64) error {
	filter := bson.M{"user": toUser}
	update := bson.M{"$inc": bson.M{"point": point}}

	opt := options.Update().SetUpsert(true)

	_, err := d.like.UpdateOne(util.Context(), filter, update, opt)
	return err
}
