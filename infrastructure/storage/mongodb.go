package storage

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/lroman242/redirective/config"
	"github.com/lroman242/redirective/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	connectTimeout = 10
	pingTimeout    = 2
	queryTimeout   = 5
	ttlInSeconds   = 2592000
)

// MongoDB type describe MongoDB storage instance.
type MongoDB struct {
	collection *mongo.Collection
}

// NewMongoDB function will create new MongoDB (implements Storage interface)
// instance according to provided StorageConfig.
func NewMongoDB(conf *config.StorageConfig) (*MongoDB, error) {
	ctxConnect, cancelFuncConnect := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancelFuncConnect()

	client, err := mongo.Connect(ctxConnect, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s", conf.User, conf.Password, conf.Host, strconv.Itoa(conf.Port))))
	if err != nil {
		log.Printf("mongodb connection failed. error: %s", err)

		return nil, err
	}

	ctxPing, cancelFuncPing := context.WithTimeout(context.Background(), pingTimeout*time.Second)
	defer cancelFuncPing()

	err = client.Ping(ctxPing, readpref.Primary())
	if err != nil {
		log.Printf("mongodb ping failed. error: %s", err)

		return nil, err
	}

	collection := client.Database(conf.Database).Collection(conf.Table)

	// Set 30 days TTL for records
	ctxCreateIndex, cancelFuncCreateIndex := context.WithTimeout(context.Background(), pingTimeout*time.Second)
	defer cancelFuncCreateIndex()

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "createdAt", Value: 1},
		},
		Options: options.Index().SetExpireAfterSeconds(ttlInSeconds),
	}

	_, err = collection.Indexes().CreateOne(ctxCreateIndex, indexModel)
	if err != nil {
		log.Printf("mongodb set ttl failed. error: %s", err)

		return nil, err
	}

	return &MongoDB{collection: collection}, nil
}

// SaveTraceResults function used to save domain.TraceResults into storage.
func (m *MongoDB) SaveTraceResults(traceResults *domain.TraceResults) (interface{}, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), queryTimeout*time.Second)
	defer cancelFunc()

	record := struct {
		*domain.TraceResults
		CreatedAt time.Time `json:"created_at"`
	}{
		traceResults,
		time.Now(),
	}

	res, err := m.collection.InsertOne(ctx, record)
	if err != nil {
		return nil, err
	}

	return res.InsertedID, nil
}

// FindTraceResults function used to find domain.TraceResults into storage using ID.
func (m *MongoDB) FindTraceResults(id interface{}) (*domain.TraceResults, error) {
	var ID primitive.ObjectID

	var err error

	switch tp := id.(type) {
	case primitive.ObjectID:
		ID = tp
	case string:
		ID, err = primitive.ObjectIDFromHex(tp)
		if err != nil {
			return nil, err
		}
	default:
		return nil, &InvalidIDTypeError{}
	}

	results := &domain.TraceResults{}

	ctx, cancelFunc := context.WithTimeout(context.Background(), queryTimeout*time.Second)
	defer cancelFunc()

	err = m.collection.FindOne(ctx, bson.M{"_id": ID}).Decode(results)
	if err != nil {
		return nil, &CannotDecodeRecordError{err: err}
	}

	return results, nil
}
