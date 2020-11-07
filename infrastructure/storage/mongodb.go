package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/lroman242/redirective/config"
	"github.com/lroman242/redirective/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"strconv"
	"time"
)

const connectTimeout = 10
const pingTimeout = 2
const queryTimeout = 5

// MongoDB type describe MongoDB storage instance
type MongoDB struct {
	collection *mongo.Collection
}

// NewMongoDB function will create new MongoDB (implements Storage interface)
// instance according to provided StorageConfig
func NewMongoDB(conf *config.StorageConfig) (*MongoDB, error) {
	ctx, _ := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s", conf.User, conf.Password, conf.Host, strconv.Itoa(conf.Port))))
	if err != nil {
		log.Printf("mongodb connection failed. error: %s", err)
		return nil, err
	}

	ctx, _ = context.WithTimeout(context.Background(), pingTimeout*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Printf("mongodb ping failed. error: %s", err)
		return nil, err
	}

	collection := client.Database(conf.Database).Collection(conf.Table)
	return &MongoDB{collection: collection}, nil
}

// SaveTraceResults function used to save domain.TraceResults into storage
func (m *MongoDB) SaveTraceResults(traceResults *domain.TraceResults) (interface{}, error) {
	ctx, _ := context.WithTimeout(context.Background(), queryTimeout*time.Second)
	res, err := m.collection.InsertOne(ctx, traceResults)
	if err != nil {
		return nil, err
	}

	return res.InsertedID, nil
}

// FindTraceResults function used to find domain.TraceResults into storage using ID
func (m *MongoDB) FindTraceResults(id interface{}) (*domain.TraceResults, error) {
	var ID primitive.ObjectID
	var err error

	switch id.(type) {
	case primitive.ObjectID:
		ID = id.(primitive.ObjectID)
	case string:
		ID, err = primitive.ObjectIDFromHex(id.(string))
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("invalid id type")
	}

	results := &domain.TraceResults{}
	ctx, _ := context.WithTimeout(context.Background(), queryTimeout*time.Second)
	err = m.collection.FindOne(ctx, bson.M{"_id": ID}).Decode(results)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("mongodb reulsts decode failed. error: %s \n", err))
	}

	return results, nil
}
