package db

import (
	"context"
	"time"

	"github.com/billettc/trades-tracker/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DATABASE = "trades-tracker"
const COLLECTION_PRICES = "prices"
const COLLECTION_OPERATIONS = "operations"

type MongoDB struct {
	client *mongo.Client
}

func NewMongoDB(address string) (*MongoDB, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(address))
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return &MongoDB{client: client}, nil

}

func (db *MongoDB) SavePrice(price *models.Price) error {
	collection := db.client.Database(DATABASE).Collection(COLLECTION_PRICES)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	filter := bson.M{"datetime": price.DateTime, "symbol": price.Symbol}
	update := bson.M{"$set": price}
	_, err := collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}

	return nil
}

func (db *MongoDB) SaveOperation(operation *models.Operation) error {
	collection := db.client.Database(DATABASE).Collection(COLLECTION_OPERATIONS)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	//filter := bson.M{"datetime": price.DateTime, "symbol": price.Symbol}
	//update := bson.M{"$set": price}
	_, err := collection.InsertOne(ctx, operation)
	if err != nil {
		return err
	}

	return nil
}

func (db MongoDB) GetPrice(symbol string, dateTime models.MinuteRoundedTime) (*models.Price, error) {
	collection := db.client.Database(DATABASE).Collection(COLLECTION_PRICES)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	filter := bson.M{"datetime": dateTime, "symbol": symbol}
	res, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if res.Next(ctx) {
		price := &models.Price{}
		err = res.Decode(&price)
		if err != nil {
			return nil, err
		}
		return price, nil
	}
	return nil, nil
}

func (db *MongoDB) ListOperations() ([]*models.Operation, error) {

	collection := db.client.Database(DATABASE).Collection(COLLECTION_OPERATIONS)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	filter := bson.M{}
	opts := options.Find().SetSort(bson.D{{"datetime", 1}})
	res, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	var operations []*models.Operation
	for res.Next(ctx) {
		o := &models.Operation{}
		err = res.Decode(&o)
		operations = append(operations, o)
	}
	return operations, nil
}
