package tests

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/tools"

	"github.com/stretchr/testify/assert"
)

func TestSetupTools(t *testing.T) {
	cfg := core.LoadConfig()
	tools := tools.Setup(cfg)
	assert.NotNil(t, tools)
}

func BenchmarkTypeAssertion(b *testing.B) {
	var x interface{} = 10

	for i := 0; i < b.N; i++ {
		_ = x.(int)
	}
}

// ➤ 115 ns/op
func BenchmarkSprintf(b *testing.B) {
	var x int = 9

	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%d", x)
	}
}

// ➤ 5 ns/op 🏆
func BenchmarkStrconv(b *testing.B) {
	var x int = 9

	for i := 0; i < b.N; i++ {
		_ = strconv.Itoa(x)
	}
}

// ➤ 7 ns/op
func BenchmarkCast(b *testing.B) {
	var x int = 9

	for i := 0; i < b.N; i++ {
		_ = string(rune(x))
	}
}
