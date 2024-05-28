package toolbox

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

func (fm fileManager) CreateFolder(path string) error {
	return os.Mkdir(fm.basePath+path, os.ModePerm)
}

func (fm fileManager) CreateFolders(paths ...string) error {
	for _, path := range paths {
		if err := fm.CreateFolder(path); err != nil {
			return err
		}
	}
	return nil
}
