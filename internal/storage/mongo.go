package storage

import (
	"context"
	"errors"
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

type movementDocument struct {
	ID               string `bson:"_id"`
	AccountID        string `bson:"account_id"`
	IsDebit          bool   `bson:"is_debit"`
	Amount           int    `bson:"amount"`
	PreviousMovement string `bson:"previous_movement"`
	PreviousBalance  int    `bson:"previous_balance"`
	CreatedAt        string `bson:"created_at"`
}

func NewMongoClient(collection *mongo.Collection) *MongoClient {
	go initCollection(collection)

	return &MongoClient{
		collection: collection,
	}
}

func initCollection(collection *mongo.Collection) {
	collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.M{"account_id": 1, "previous_balance": 1},
		Options: options.Index().SetUnique(true),
	})

	collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.M{"account_id": 1, "created_at": 0},
	})
}

func (cli *MongoClient) All(ctx context.Context, aID account.ID) ([]*movement.Movement, error) {
	timeout, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"created_at", 1}})
	cursor, err := cli.collection.Find(timeout, bson.M{"account_id": aID}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(timeout)

	transactions := make([]*movement.Movement, 0)
	for cursor.Next(timeout) {
		t := &movementDocument{}
		err := cursor.Decode(t)
		if err != nil {
			return nil, err
		}

		cAt, _ := time.Parse(time.RFC3339, t.CreatedAt)
		transactions = append(transactions, &movement.Movement{
			ID:               movement.ID(t.ID),
			AccountID:        account.ID(t.AccountID),
			IsDebit:          t.IsDebit,
			Amount:           t.Amount,
			PreviousMovement: movement.ID(t.PreviousMovement),
			PreviousBalance:  t.PreviousBalance,
			CreatedAt:        cAt,
		})
	}
	return transactions, nil
}

func (cli *MongoClient) Last(ctx context.Context, id account.ID) (*movement.Movement, error) {
	timeout, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	result := &movementDocument{}
	findOptions := options.FindOne().SetSort(bson.D{{"created_at", -1}})
	err := cli.collection.FindOne(timeout, bson.M{"account_id": id}, findOptions).Decode(&result)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}

	cAt, _ := time.Parse(time.RFC3339, result.CreatedAt)
	return &movement.Movement{
		ID:               movement.ID(result.ID),
		AccountID:        account.ID(result.AccountID),
		IsDebit:          result.IsDebit,
		Amount:           result.Amount,
		PreviousMovement: movement.ID(result.PreviousMovement),
		PreviousBalance:  result.PreviousBalance,
		CreatedAt:        cAt,
	}, nil
}

func (cli *MongoClient) Create(ctx context.Context, m *movement.Movement) error {
	timeout, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	doc := movementDocument{
		ID:               m.ID.String(),
		AccountID:        string(m.AccountID),
		IsDebit:          m.IsDebit,
		Amount:           m.Amount,
		PreviousMovement: m.PreviousMovement.String(),
		PreviousBalance:  m.PreviousBalance,
		CreatedAt:        m.CreatedAt.Format(time.RFC3339),
	}

	_, err := cli.collection.InsertOne(timeout, doc)

	return err
}
