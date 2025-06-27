package tools

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/utils"

	"github.com/google/uuid"
)

type fileDownloader struct {
	client *http.Client
}

func NewFileDownloader(client *http.Client) core.FileDownloader {
	return &fileDownloader{client: client}
}

var _ core.FileDownloader = fileDownloader{}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Returns (filePath, fileSize, error).
// We avoid passing a context because we don't want to cancel the download if the context is cancelled.
func (fd fileDownloader) DownloadFileToDisk(url, fileExt string, extraHeaders map[string]string) (string, int, error) {

	status, content, err := utils.GET(context.Background(), url, extraHeaders, "", fd.client)
	if err != nil {
		return "", 0, fmt.Errorf("error downloading file from %s: %w", url, err)
	}
	if status != http.StatusOK {
		return "", 0, fmt.Errorf("error downloading file from %s: received status %d with body %s", url, status, string(content))
	}

	filePath := fmt.Sprintf("./etc/downloads/%s.%s", uuid.New().String()[:12], fileExt)
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return "", 0, fmt.Errorf("error saving file from %s to disk: %w", url, err)
	}

	return filePath, len(content), nil
}
