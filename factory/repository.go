package factory

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/vsmoraes/open-ledger/storage"
)

func NewDBRepository() (*storage.MongoClient, *mongo.Client) {
	mc, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017/"))
	if err != nil {
		panic(err.Error())
	}

	if err = mc.Connect(context.Background()); err != nil {
		panic(err.Error())
	}

	collection := mc.Database("open-ledger").Collection("movements")
	mongoClient := storage.NewMongoClient(collection)

	return mongoClient, mc
}
