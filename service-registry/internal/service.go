package serviceRegistry

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	ErrServiceNotFound         = errors.New("service not found")
	ErrServiceInstanceNotFound = errors.New("service instance not found")
	ErrServiceAlreadyExists    = errors.New("service already exists")
)

type ServiceRegistryService interface {
	GetService(ctx context.Context, id string) (*Service, error)
	RegisterService(ctx context.Context, svc *Service) error

	GetServiceInstance(ctx context.Context, id string) (*ServiceInstance, error)
	RegisterServiceInstance(ctx context.Context, instance *ServiceInstance) error
}

type serviceRegistryService struct {
	repo ServiceRepository
}

func NewServiceRegistryService(repo ServiceRepository) ServiceRegistryService {
	return &serviceRegistryService{repo: repo}
}

func (s *serviceRegistryService) GetService(ctx context.Context, id string) (*Service, error) {
	if id == "" {
		return nil, errors.New("service id must not be empty")
	}

	svc, err := s.repo.GetService(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrServiceNotFound
		}
		return nil, err
	}

	return svc, nil
}

func (s *serviceRegistryService) RegisterService(ctx context.Context, svc *Service) error {
	if svc == nil {
		return errors.New("service must not be nil")
	}
	if svc.ID == bson.NilObjectID {
		return errors.New("service id must not be empty")
	}
	if svc.Name == "" {
		return errors.New("service name must not be empty")
	}

	existing, err := s.GetService(ctx, svc.ID.Hex())
	if err != nil && !errors.Is(err, ErrServiceNotFound) {
		return err
	}
	if existing != nil {
		return ErrServiceAlreadyExists
	}

	if svc.CreatedAt.IsZero() {
		svc.CreatedAt = time.Now().UTC()
	}

	return s.repo.AddService(ctx, svc)
}

func (s *serviceRegistryService) GetServiceInstance(ctx context.Context, id string) (*ServiceInstance, error) {
	if id == "" {
		return nil, errors.New("instance id must not be empty")
	}

	instance, err := s.repo.GetServiceInstance(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrServiceInstanceNotFound
		}
		return nil, err
	}

	return instance, nil
}

func (s *serviceRegistryService) RegisterServiceInstance(ctx context.Context, instance *ServiceInstance) error {
	if instance == nil {
		return errors.New("instance must not be nil")
	}
	if instance.ID == bson.NilObjectID {
		return errors.New("instance id must not be empty")
	}
	if instance.ServiceID == bson.NilObjectID {
		return errors.New("instance service_id must not be empty")
	}

	_, err := s.GetService(ctx, instance.ServiceID.Hex())
	if err != nil {
		return err
	}

	if instance.LastHeartbeat.IsZero() {
		instance.LastHeartbeat = time.Now().UTC()
	}

	return s.repo.AddServiceInstance(ctx, instance)
}
