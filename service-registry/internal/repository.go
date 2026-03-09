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

	GetServiceInstance(ctx context.Context, id string) (*ServiceInstance, error)
	AddServiceInstance(ctx context.Context, serviceID string, host string, port int, status InstanceStatus) (*ServiceInstance, error)
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
	res := r.serviceColl.FindOne(ctx, bson.M{"name": name})
	svc := &Service{}

	if err := res.Decode(svc); err != nil {
		if err != mongo.ErrNoDocuments {
			fmt.Println("Error creating a service")
			return nil, err
		}

		svc = &Service{
			ID:        bson.NewObjectID(),
			Name:      name,
			CreatedAt: time.Now().UTC(),
		}
		_, err = r.serviceColl.InsertOne(ctx, svc)

		if err != nil {
			fmt.Println("Error inserting service:", err)
			return nil, err
		}
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

func (r *mongoServiceRepository) AddServiceInstance(ctx context.Context, serviceID string, host string, port int, status InstanceStatus) (*ServiceInstance, error) {
	res := r.serviceInstanceColl.FindOne(ctx, bson.M{
		"serviceID": serviceID,
		"host":      host,
		"port":      port,
	})
	instance := &ServiceInstance{}

	if err := res.Decode(*instance); err != nil {
		if err != mongo.ErrNoDocuments {
			fmt.Printf("Error creating an instance")
			return nil, err
		}

		oid, err := bson.ObjectIDFromHex(serviceID)

		if err != nil {
			fmt.Printf("Error converting to ObjectID\n")
			return nil, err
		}

		instance = &ServiceInstance{
			ID:            bson.NewObjectID(),
			ServiceID:     oid,
			Host:          host,
			Port:          port,
			Status:        status,
			LastHeartbeat: time.Now().UTC(),
		}
		_, err = r.serviceInstanceColl.InsertOne(ctx, instance)
		if err != nil {
			return nil, err
		}
	}

	return instance, nil
}
