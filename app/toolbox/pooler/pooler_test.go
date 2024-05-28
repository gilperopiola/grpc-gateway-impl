package pooler

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestStringsPool(t *testing.T) {

	// New pool of 20 strings, random values
	poolSize := 20
	newStringFn := func() string {
		return fmt.Sprintf("%d", rand.Intn(100))
	}

	pool := NewPool[string](poolSize, newStringFn)

	// Check pool size
	if pool.Len() != poolSize {
		t.Errorf("Expected pool size to be %d, got %d", poolSize, pool.Len())
	}

	// Take all strings from the pool
	releaseFns := make([]func(), poolSize)
	for i := 0; i < poolSize; i++ {
		taken, releaseFn := pool.TakeOne()
		if taken == "" {
			t.Errorf("Expected taken string to be non-empty")
		}
		releaseFns[i] = releaseFn
	}

	// Take one more string from the pool (should create a new one)
	newString, releaseFn := pool.TakeOne()
	if newString == "" {
		t.Errorf("Expected new taken string to be non-empty")
	}

	if pool.Len() != poolSize+1 {
		t.Errorf("Expected new pool size to be %d, got %d", poolSize+1, pool.Len())
	}

	// Take one more string from the pool (should create a new one)
	newString2, releaseFn2 := pool.TakeOne()
	if newString2 == "" {
		t.Errorf("Expected new taken string 2 to be non-empty")
	}

	if pool.Len() != poolSize+2 {
		t.Errorf("Expected new pool size 2 to be %d, got %d", poolSize+2, pool.Len())
	}

	// Release all strings
	for _, releaseFn := range releaseFns {
		releaseFn()
	}
	releaseFn()
	releaseFn2()

	newerString, releaseFn3 := pool.TakeOne()
	if newerString == "" {
		t.Errorf("Expected newer taken string to be non-empty")
	}

	if pool.Len() != poolSize+2 {
		t.Errorf("Expected newer pool size to be %d, got %d", poolSize+2, pool.Len())
	}

	releaseFn3()
}
