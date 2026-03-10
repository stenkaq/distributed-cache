package serviceRegistry

import (
	"context"
	"iter"

	"github.com/maypok86/otter/v2"
)

type AddServiceInstanceParams struct {
	ServiceID string
	Host      string
	Port      int
	Hash      uint32
	Status    InstanceStatus
}

type ServiceInstanceRingRepository interface {
	GetRing(ctx context.Context) iter.Seq2[int, *ServiceInstance]
	AddServiceInstance(ctx context.Context, p AddServiceInstanceParams) (*ServiceInstance, bool)
}

type ringRepository struct {
	ring otter.Cache[int, *ServiceInstance]
}

func NewRingRepository(c otter.Cache[int, *ServiceInstance]) ServiceInstanceRingRepository {
	return &ringRepository{ring: c}
}

func (r *ringRepository) GetRing(ctx context.Context) iter.Seq2[int, *ServiceInstance] {
	return r.ring.All()
}

func (r *ringRepository) AddServiceInstance(ctx context.Context, p AddServiceInstanceParams) (*ServiceInstance, bool) {
	instance := &ServiceInstance{}

	val, created := r.ring.SetIfAbsent(int(p.Hash), instance)

	return val, created
}
