package tools

import (
	"fmt"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"

	"google.golang.org/grpc/metadata"
)

var _ core.MetadataGetter = (*grpcMetadataGetter)(nil)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type grpcMetadataGetter struct{}

func NewMetadataGetter() core.MetadataGetter {
	return &grpcMetadataGetter{}
}

func (gmm *grpcMetadataGetter) GetMD(ctx core.Ctx, key string) (string, error) {
	if val := metadata.ValueFromIncomingContext(ctx, key); len(val) > 0 {
		return val[0], nil
	}
	return "", fmt.Errorf("metadata with key %s not found", key)
}
