package tools

import (
	"reflect"

	"golang.org/x/time/rate"
)

type pooler struct {
	userRateLimitingPool Pool[poolResource[*rate.Limiter]]
}

type Pool[Resource any] interface {
	AddResources(int) []Resource
	GetResource() Resource
	ReleaseResource(Resource)
	SetNewResourceFn(func() Resource)
	Len() int
}

type pool[Resource any] struct {
	resources     []poolResource[Resource]
	newResourceFn func() Resource
}

type poolResource[Resource any] struct {
	Resource  Resource
	Available bool
}

func NewPool[Resource any](n int, newFn func() Resource) Pool[Resource] {
	pool := &pool[Resource]{
		resources:     []poolResource[Resource]{},
		newResourceFn: newFn,
	}
	pool.AddResources(n)
	return pool
}

func (p *pool[Resource]) AddResources(n int) []Resource {
	added := []Resource{}

	for i := 0; i < n; i++ {
		newRes := p.newResourceFn()
		p.resources = append(p.resources,
			poolResource[Resource]{Resource: newRes, Available: true},
		)
		added = append(added, newRes)
	}

	return added
}

func (p *pool[Resource]) GetResource() Resource {
	for i, resource := range p.resources {
		if resource.Available {
			p.resources[i].Available = false
			return resource.Resource
		}
	}

	return p.AddResources(1)[0] // We only add 1 instance so we get the first element.
}

func (p *pool[Resource]) ReleaseResource(resource Resource) {
	for i, poolResource := range p.resources {
		if reflect.DeepEqual(poolResource.Resource, resource) {
			p.resources[i].Available = true
		}
	}
}

func (p *pool[Resource]) SetNewResourceFn(fn func() Resource) {
	p.newResourceFn = fn
}

func (p *pool[Resource]) Len() int {
	return len(p.resources)
}
