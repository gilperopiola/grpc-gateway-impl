package grpc

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/cfg"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/dependencies"
	"golang.org/x/time/rate"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/grpc"
)

/* ----------------------------------- */
/*        - Interceptor Tests -        */
/* ----------------------------------- */

func TestRateLimiterInterceptor(t *testing.T) {

	// Create rate limiter and logger.
	rateLimiterCfg := &cfg.RateLimiterConfig{MaxTokens: 2, TokensPerSecond: 1}
	limiter := rate.NewLimiter(rate.Limit(rateLimiterCfg.TokensPerSecond), rateLimiterCfg.MaxTokens)
	logger := zap.NewNop()

	// Config: allow 1 request per second with a max of 2.
	interceptor := getGRPCRateLimiterInterceptor(limiter, &dependencies.Logger{Logger: logger})

	// Mock handler to simulate gRPC method execution.
	mockHandler := func(ctx context.Context, req interface{}) (interface{}, error) { return nil, nil }

	// Simulate two quick successive requests, both should be allowed.
	for i := 0; i < 2; i++ {
		if _, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{}, mockHandler); err != nil {
			t.Errorf("Request %d was unexpectedly limited: %v", i+1, err)
		}
	}

	// Simulate a third request which should be limited.
	if _, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{}, mockHandler); err == nil {
		t.Error("Expected the third request to be rate limited, but it was not")
	}

	// Wait 1 second and retry, should be allowed.
	time.Sleep(1 * time.Second)
	if _, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{}, mockHandler); err != nil {
		t.Errorf("Expected the request to be allowed after waiting, but it was limited: %v", err)
	}
}

func TestGRPCLoggerInterceptor(t *testing.T) {
	core, recorded := observer.New(zap.InfoLevel)
	logger := zap.New(core)
	interceptor := getGRPCLoggerInterceptor(&dependencies.Logger{Logger: logger})

	// Mock handler to simulate gRPC method execution.
	mockHandler := func(ctx context.Context, req interface{}) (interface{}, error) { return "mockResponse", nil }

	// Simulate a gRPC request.
	if _, err := interceptor(context.Background(), "mockRequest", &grpc.UnaryServerInfo{FullMethod: "/test/method"}, mockHandler); err != nil {
		t.Fatalf("Interceptor returned an error: %v", err)
	}

	// Verify the log entry.
	if recorded.Len() != 1 {
		t.Errorf("Expected 1 log entry, got %d", recorded.Len())
	}
	entry := recorded.All()[0]
	if entry.Level != zap.InfoLevel {
		t.Errorf("Expected InfoLevel log, got %s", entry.Level)
	}
	endpoint := entry.ContextMap()["endpoint"].(string)
	if !bytes.Contains([]byte(entry.Message), []byte("gRPC Request")) || endpoint != "/test/method" {
		t.Errorf("Log entry does not contain expected message: %s", entry.Message)
	}
}
