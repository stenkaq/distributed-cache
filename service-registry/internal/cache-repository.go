package serviceRegistry

import (
	"context"
	"sort"
)

type AddServiceInstanceParams struct {
	ServiceID string
	Host      string
	Port      int
	Hash      uint32
	Status    InstanceStatus
}

type ServiceInstanceRingRepository interface {
	GetRing(ctx context.Context) []*ServiceInstance
	AddServiceInstance(ctx context.Context, p AddServiceInstanceParams) (*ServiceInstance, bool)
}

type ringRepository struct {
	ring []*ServiceInstance
}

func NewRingRepository(c []*ServiceInstance) ServiceInstanceRingRepository {
	return &ringRepository{ring: c}
}

func (r *ringRepository) GetRing(ctx context.Context) []*ServiceInstance {
	return r.ring
}

func (r *ringRepository) AddServiceInstance(ctx context.Context, p AddServiceInstanceParams) (*ServiceInstance, bool) {
	instance := &ServiceInstance{
		Hash: p.Hash,
		Host: p.Host,
		Port: p.Port,
	}

	idx := sort.Search(len(r.ring), func(i int) bool {
		return r.ring[i].Hash >= instance.Hash
	})

	if idx < len(r.ring) && r.ring[idx].Hash == instance.Hash {
		return r.ring[idx], false
	}

	r.ring = append(r.ring, nil)
	copy(r.ring[idx+1:], r.ring[idx:])
	r.ring[idx] = instance

	return instance, true
}
