package tools

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/logs"

	"github.com/google/uuid"
)

func NewIDGenerator[T core.IDType](generateFn func() T) IDGenerator[T] {
	return &idGenerator[T]{GenerateIDFn: generateFn}
}

type IDGenerator[T core.IDType] interface {
	GenerateID() T
}

type idGenerator[T core.IDType] struct {
	GenerateIDFn func() T
}

func (g *idGenerator[T]) GenerateID() T {
	return g.GenerateIDFn()
}

func GenerateUUID() uuid.UUID {
	return uuid.New()
}

func GenerateUUIDShort() string {
	id := uuid.New().String()[:8]
	logs.LogSimple("Generated short UUID", id)
	return id
}
