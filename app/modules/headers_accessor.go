package modules

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/grpc/metadata"
)

type KeyValStoreAccessor interface {
	Get(string) (string, error)
}

type httpHeadersAccessor struct {
	headers http.Header
}

type grpcMetadataAccessor struct {
	ctx context.Context
}

func NewHTTPHeadersAccessor(headers any) KeyValStoreAccessor {
	return &httpHeadersAccessor{headers.(http.Header)}
}

func NewGRPCMetadataAccessor(ctx any) KeyValStoreAccessor {
	return &grpcMetadataAccessor{ctx.(context.Context)}
}

func (a *httpHeadersAccessor) Get(key string) (string, error) {
	if value := a.headers.Get(key); value != "" {
		return value, nil
	}
	return "", fmt.Errorf("header not found")
}

func (a *grpcMetadataAccessor) Get(key string) (string, error) {
	if md := metadata.ValueFromIncomingContext(a.ctx, key); len(md) > 0 {
		return md[0], nil
	}
	return "", fmt.Errorf("metadata not found")
}
