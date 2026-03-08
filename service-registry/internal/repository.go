package serviceRegistry

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ServiceRepository interface {
	GetService(ctx context.Context, id string) (*Service, error)
	AddService(ctx context.Context, svc *Service) error

	GetServiceInstance(ctx context.Context, id string) (*ServiceInstance, error)
	AddServiceInstance(ctx context.Context, instance *ServiceInstance) error
}

type mongoServiceRepository struct {
	serviceColl         *mongo.Collection
	serviceInstanceColl *mongo.Collection
}

func NewServiceRepository(db *DB) ServiceRepository {
	return &mongoServiceRepository{
		serviceColl:         db.serviceColl,
		serviceInstanceColl: db.serviceInstanceColl,
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

func (r *mongoServiceRepository) AddService(ctx context.Context, svc *Service) error {
	_, err := r.serviceColl.InsertOne(ctx, svc)
	return err
}

func (r *mongoServiceRepository) GetServiceInstance(ctx context.Context, id string) (*ServiceInstance, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var instance ServiceInstance
	err = r.serviceInstanceColl.FindOne(ctx, bson.M{"_id": oid}).Decode(&instance)
	if err != nil {
		return nil, err
	}
	return &instance, nil
}

func (r *mongoServiceRepository) AddServiceInstance(ctx context.Context, svc *ServiceInstance) error {
	_, err := r.serviceInstanceColl.InsertOne(ctx, svc)
	return err
}
