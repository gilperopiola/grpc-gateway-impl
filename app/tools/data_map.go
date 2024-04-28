package tools

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/grpc/metadata"
)

// DataMap is a really powerful interface, we use it to abstract access to headers and metadata.
// It can be used with any key-value pair storage.
// Can this be improved with generics?
type DataMap interface {
	Get(string) (string, error)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type httpHeadersMap struct {
	headers http.Header
}

func NewHTTPHeadersMap(headers any) DataMap {
	return &httpHeadersMap{headers.(http.Header)}
}

func (hhm *httpHeadersMap) Get(key string) (string, error) {
	if val := hhm.headers.Get(key); val != "" {
		return val, nil
	}
	return "", fmt.Errorf("header with key %s not found", key)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type grpcMetadataMap struct {
	ctx context.Context // we can access GRPC Metadata from the context
}

func NewGRPCMetadataMap(ctx any) DataMap {
	return &grpcMetadataMap{ctx.(context.Context)}
}

func (gmm *grpcMetadataMap) Get(key string) (string, error) {
	if val := metadata.ValueFromIncomingContext(gmm.ctx, key); len(val) > 0 {
		return val[0], nil
	}
	return "", fmt.Errorf("metadata with key %s not found", key)
}
