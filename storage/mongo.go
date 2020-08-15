package storage

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/vsmoraes/open-ledger/ledger/account"
	"github.com/vsmoraes/open-ledger/ledger/movement"
)

const DefaultTimeout = 5 * time.Second

type MongoClient struct {
	collection *mongo.Collection
}

func New(collection *mongo.Collection) *MongoClient {
	return &MongoClient{
		collection: collection,
	}
}

func (cli *MongoClient) All(ctx context.Context, id account.ID) ([]*movement.Movement, error) {
	timeout, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"created_at", 1}})
	cursor, err := cli.collection.Find(timeout, bson.M{"id": id}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(timeout)

	var transactions []*movement.Movement
	for cursor.Next(timeout) {
		var t *movement.Movement
		err := cursor.Decode(t)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, t)
	}
	return transactions, nil
}

func (cli *MongoClient) Last(ctx context.Context, id account.ID) (*movement.Movement, error) {
	timeout, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	result := &movement.Movement{}
	findOptions := options.FindOne()
	findOptions.SetSort(bson.D{{"created_at", -1}})
	err := cli.collection.FindOne(timeout, bson.M{"id": id}, findOptions).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (cli *MongoClient) Create(ctx context.Context, t *movement.Movement) error {
	timeout, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	_, err := cli.collection.InsertOne(timeout, t)

	return err
}
