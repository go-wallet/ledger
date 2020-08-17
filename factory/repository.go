package factory

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/vsmoraes/open-ledger/config"
	"github.com/vsmoraes/open-ledger/storage"
)

func NewDBRepository() (*storage.MongoClient, *mongo.Client) {
	mc, err := mongo.NewClient(options.Client().ApplyURI(config.Config().MongoDB.URI))
	if err != nil {
		panic(err.Error())
	}

	if err = mc.Connect(context.Background()); err != nil {
		panic(err.Error())
	}

	collection := mc.Database(config.Config().MongoDB.Database).
		Collection(config.Config().MongoDB.MovementsCollection)
	mongoClient := storage.NewMongoClient(collection)

	return mongoClient, mc
}
