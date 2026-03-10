package serviceRegistry

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ServiceRepository interface {
	GetService(ctx context.Context, id string) (*Service, error)
	AddService(ctx context.Context, name string) (*Service, error)
}

type mongoServiceRepository struct {
	serviceColl *mongo.Collection
}

func NewServiceRepository(db *DB) ServiceRepository {
	return &mongoServiceRepository{
		serviceColl: db.serviceColl,
	}
}

func (r *mongoServiceRepository) GetService(ctx context.Context, id string) (*Service, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var svc Service
	err = r.serviceColl.FindOne(ctx, bson.M{"_id": oid}).Decode(&svc)
	if err != nil {
		return nil, err
	}
	return &svc, nil
}

func (r *mongoServiceRepository) AddService(ctx context.Context, name string) (*Service, error) {
	res := r.serviceColl.FindOne(ctx, bson.M{"name": name})
	svc := &Service{}

	if err := res.Decode(svc); err != nil {
		if err != mongo.ErrNoDocuments {
			fmt.Println("Error creating a service")
			return nil, err
		}

		now := time.Now().UTC()

		svc = &Service{
			ID:        bson.NewObjectID(),
			Name:      name,
			CreatedAt: now,
			UpdatedAt: now,
		}
		_, err = r.serviceColl.InsertOne(ctx, svc)

		if err != nil {
			fmt.Println("Error inserting service:", err)
			return nil, err
		}
	}

	return svc, nil
}
