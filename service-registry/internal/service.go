package serviceRegistry

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	ErrServiceNotFound         = errors.New("service not found")
	ErrServiceInstanceNotFound = errors.New("service instance not found")
)

type ServiceRegistryService interface {
	GetService(ctx context.Context, id string) (*Service, error)
	RegisterService(ctx context.Context, name string, host string, port *int) (*Service, error)

	GetServiceInstance(ctx context.Context, id string) (*ServiceInstance, error)
	RegisterServiceInstance(ctx context.Context, serviceID string, host string, port *int, status InstanceStatus) (*ServiceInstance, error)
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

func (s *serviceRegistryService) RegisterService(ctx context.Context, name string, host string, port *int) (*Service, error) {
	if name == "" || host == "" || port == nil {
		return nil, errors.New("Params must not be empty")
	}

	return s.repo.AddService(ctx, name)
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

func (s *serviceRegistryService) RegisterServiceInstance(ctx context.Context, serviceID string, host string, port *int, status InstanceStatus) (*ServiceInstance, error) {
	if serviceID == "" {
		return nil, errors.New("instance service_id must not be empty")
	}
	if host == "" {
		return nil, errors.New("instance host must not be empty")
	}
	if port == nil {
		return nil, errors.New("instance port must not be empty")
	}

	return s.repo.AddServiceInstance(ctx, serviceID, host, *port, status)
}
