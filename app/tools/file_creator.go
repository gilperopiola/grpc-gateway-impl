package tools

import (
	"os"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
)

var _ core.FileCreator = (*fileCreator)(nil)

type fileCreator struct {
}

func NewFileCreator() core.FileCreator {
	return &fileCreator{}
}

func (fc fileCreator) CreateFolders(paths ...string) (err error) {
	for _, path := range paths {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}
