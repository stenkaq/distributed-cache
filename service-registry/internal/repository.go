package serviceRegistry

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ServiceRepository interface {
	GetService(ctx context.Context, id string) (*Service, error)
	AddService(ctx context.Context, name string) (*Service, error)

	GetServiceInstance(ctx context.Context, id string) (*ServiceInstance, error)
	AddServiceInstance(ctx context.Context, serviceID bson.ObjectID, host string, port int, status InstanceStatus) (*ServiceInstance, error)
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

func (r *mongoServiceRepository) AddService(ctx context.Context, name string) (*Service, error) {
	svc := &Service{
		ID:        bson.NewObjectID(),
		Name:      name,
		CreatedAt: time.Now().UTC(),
	}
	_, err := r.serviceColl.InsertOne(ctx, svc)
	if err != nil {
		return nil, err
	}
	return svc, nil
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

func (r *mongoServiceRepository) AddServiceInstance(ctx context.Context, serviceID bson.ObjectID, host string, port int, status InstanceStatus) (*ServiceInstance, error) {
	instance := &ServiceInstance{
		ID:            bson.NewObjectID(),
		ServiceID:     serviceID,
		Host:          host,
		Port:          port,
		Status:        status,
		LastHeartbeat: time.Now().UTC(),
	}
	_, err := r.serviceInstanceColl.InsertOne(ctx, instance)
	if err != nil {
		return nil, err
	}
	return instance, nil
}
