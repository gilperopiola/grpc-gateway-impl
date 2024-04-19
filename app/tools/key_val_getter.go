package tools

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/grpc/metadata"
)

// Used to abstract the way we access headers and metadata, or any other key-value pair storage.
type KeyValueMapper interface {
	Get(string) (string, error)
}

type httpHeadersKVMapper struct {
	headers http.Header
}

func NewHTTPHeadersMapper(headers any) KeyValueMapper {
	return &httpHeadersKVMapper{headers.(http.Header)}
}

func (a *httpHeadersKVMapper) Get(headerKey string) (string, error) {
	if headerValue := a.headers.Get(headerKey); headerValue != "" {
		return headerValue, nil
	}
	return "", fmt.Errorf("header with key %s not found", headerKey)
}

type grpcMetadataKVMapper struct {
	ctx context.Context
}

func NewGRPCMetadataMapper(ctx any) KeyValueMapper {
	return &grpcMetadataKVMapper{ctx.(context.Context)}
}

func (a *grpcMetadataKVMapper) Get(mdKey string) (string, error) {
	if mdValue := metadata.ValueFromIncomingContext(a.ctx, mdKey); len(mdValue) > 0 {
		return mdValue[0], nil
	}
	return "", fmt.Errorf("metadata with key %s not found", mdKey)
}
