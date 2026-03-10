package serviceRegistry

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"slices"

	"github.com/spaolacci/murmur3"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	ErrServiceNotFound         = errors.New("service not found")
	ErrServiceInstanceNotFound = errors.New("service instance not found")
)

type RegisterServiceInstanceParams struct {
	ServiceID string
	Host      string
	Port      *int
	Hash      *uint32
	Status    InstanceStatus
}

type ServiceRegistryService interface {
	GetService(ctx context.Context, id string) (*Service, error)
	RegisterService(ctx context.Context, name string, host string, port *int) (*Service, error)

	GetRing(ctx context.Context) iter.Seq2[int, *ServiceInstance]
	RegisterServiceInstance(ctx context.Context, p RegisterServiceInstanceParams) (*ServiceInstance, error)
}

type serviceRegistryService struct {
	repo ServiceRepository
	ring ServiceInstanceRingRepository
}

func NewServiceRegistryService(repo ServiceRepository, ring ServiceInstanceRingRepository) ServiceRegistryService {
	return &serviceRegistryService{repo: repo, ring: ring}
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

func (s *serviceRegistryService) GetRing(ctx context.Context) iter.Seq2[int, *ServiceInstance] {
	return slices.All(s.ring.GetRing(ctx))
}

func (s *serviceRegistryService) RegisterServiceInstance(ctx context.Context, p RegisterServiceInstanceParams) (*ServiceInstance, error) {
	if p.ServiceID == "" {
		return nil, errors.New("instance service_id must not be empty")
	}
	if p.Host == "" {
		return nil, errors.New("instance host must not be empty")
	}
	if p.Port == nil {
		return nil, errors.New("instance port must not be empty")
	}

	var hash uint32
	if p.Hash != nil {
		hash = *p.Hash
	} else {
		hash = GetHash(fmt.Sprintf("%s:%d", p.Host, *p.Port))
	}

	instance, _ := s.ring.AddServiceInstance(ctx, AddServiceInstanceParams{
		ServiceID: p.ServiceID,
		Host:      p.Host,
		Port:      *p.Port,
		Hash:      hash,
		Status:    p.Status,
	})

	return instance, nil
}

func GetHash(key string) uint32 {
	return murmur3.Sum32([]byte(key))
}
