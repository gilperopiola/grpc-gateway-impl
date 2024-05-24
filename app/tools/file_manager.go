package tools

import (
	"os"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
)

var _ core.FileManager = (*fileManager)(nil)

type fileManager struct {
	basePath string
}

func NewFileManager(basePath string) core.FileManager {
	return &fileManager{basePath}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (fc fileManager) CreateFolder(path string) error {
	return os.Mkdir(fc.basePath+path, os.ModePerm)
}

func (fc fileManager) CreateFolders(paths ...string) error {
	for _, path := range paths {
		if err := fc.CreateFolder(path); err != nil {
			return err
		}
	}
	return nil
}
