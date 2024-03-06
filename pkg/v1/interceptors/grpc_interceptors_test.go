package interceptors

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	v1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1"
	"github.com/gilperopiola/grpc-gateway-impl/server/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"github.com/bufbuild/protovalidate-go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/* ----------------------------------- */
/*        - Interceptor Tests -        */
/* ----------------------------------- */

func TestNewProtoValidator(t *testing.T) {
	validator := NewProtoValidator()
	assert.NotNil(t, validator)
}

func TestFromValidationErrToGRPCInvalidArgErr(t *testing.T) {
	tests := []struct {
		name             string
		validationErr    error
		invalidArgErrMsg string
	}{
		{
			name: "validation_error",
			validationErr: &protovalidate.ValidationError{
				Violations: []*validate.Violation{{FieldPath: "field", Message: "is required"}},
			},
			invalidArgErrMsg: fmt.Sprintf(v1.ErrMsgInValidation, "field is required"),
		},
		{
			name:             "unexpected_error",
			validationErr:    errors.New("unexpected error"),
			invalidArgErrMsg: fmt.Sprintf(v1.ErrMsgInValidationUnexpected, "unexpected error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validationErrToInvalidArgErr(tt.validationErr)
			grpcStatus, _ := status.FromError(err)

			assert.Equal(t, codes.InvalidArgument, grpcStatus.Code())
			assert.Equal(t, tt.invalidArgErrMsg, grpcStatus.Message())
		})
	}
}

func TestRateLimiterInterceptor(t *testing.T) {

	// Config: allow 1 request per second with a max of 2.
	interceptor := newGRPCRateLimiterInterceptor(config.RateLimiterConfig{
		MaxTokens:       2,
		TokensPerSecond: 1,
	})

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
	interceptor := newGRPCLoggerInterceptor(logger)

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
	method := entry.ContextMap()["method"].(string)
	if !bytes.Contains([]byte(entry.Message), []byte("gRPC Request")) || method != "/test/method" {
		t.Errorf("Log entry does not contain expected message: %s", entry.Message)
	}
}
