package serviceRegistry

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DBInterface interface {
	ServiceCollection() *mongo.Collection
	ServiceInstanceCollection() *mongo.Collection
	Close() error
}

type DB struct {
	client              *mongo.Client
	serviceColl         *mongo.Collection
	serviceInstanceColl *mongo.Collection
}

func NewDB() (*DB, error) {
	uri := "http"

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("connecting to mongo: %w", err)
	}

	db := client.Database("service_registry")

	return &DB{
		client:              client,
		serviceColl:         db.Collection("services"),
		serviceInstanceColl: db.Collection("service_instances"),
	}, nil
}

func (db *DB) Close() error {
	return db.client.Disconnect(context.TODO())
}
