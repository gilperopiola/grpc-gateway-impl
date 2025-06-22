package tools

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/logs"

	"github.com/google/uuid"
)

type idGenerator[T core.IDType] struct {
	GenerateIDFn func() T
}

func NewIDGenerator[T core.IDType](fn func() T) core.IDGenerator[T] {
	return &idGenerator[T]{GenerateIDFn: fn}
}

func (g *idGenerator[T]) GenerateID() T {
	id := g.GenerateIDFn()
	logs.LogSimple("Generated ID", id)
	return id
}

/* -~-~-~-~- ID Generation Funcs -~-~-~-~- */

func GenerateUUID() uuid.UUID {
	return uuid.New()
}

func GenerateUUIDShort() string {
	id := uuid.New().String()[:8]
	return id
}

func GenerateCustomUUID() string {
	threeLetterAdjectives := []string{"red", "wet", "sad", "hot", "big", "new", "old", "bad", "fun", "mad", "raw", "low", "fat", "lil"}
	fourLetterNouns := []string{"frog", "wolf", "bear", "lion", "fish", "bird", "lava", "cave", "tuna", "crab", "land", "ruby"}
	months := []string{"JAN", "FEB", "MAR", "APR", "MAY", "JUN", "JUL", "AUG", "SEP", "OCT", "NOV", "DEC"}

	idFormat := "%s-%s-%s-%s-%s-%s"

	adj := threeLetterAdjectives[rand.Intn(len(threeLetterAdjectives))]
	noun := fourLetterNouns[rand.Intn(len(fourLetterNouns))]
	day := fmt.Sprintf("%02d", time.Now().Day())
	month := months[time.Now().Month()-1]
	year := fmt.Sprintf("%02d", time.Now().Year()-2000)
	shortUUID := uuid.New().String()[:8]

	return fmt.Sprintf(idFormat, adj, noun, day, month, year, shortUUID)
}
