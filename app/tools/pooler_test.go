package tools

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPool(t *testing.T) {
	nResources := 20

	pool := NewPool[string](nResources, func() string { return "" })

	for i := 0; i < nResources; i++ {
		resource := pool.GetResource()
		resource += strconv.Itoa(i)
		pool.ReleaseResource(resource)
	}

	for i := 0; i < nResources; i++ {
		resource := pool.GetResource()
		fmt.Printf("%s", resource)
	}

	assert.Equal(t, 13243, nResources)
}
