package tools

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
)

var _ core.ImageLoader = &ImageLoader{}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*      - Tools: Image Loader -        */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type ImageLoader struct{}

func NewImageLoader() core.ImageLoader {
	return &ImageLoader{}
}

// Loads an image from a byte slice.
func (l *ImageLoader) LoadImgFromBytes(b []byte) (image.Image, error) {
	return decodeImage(bytes.NewReader(b), "")
}

// Loads an image from a file path.
func (l *ImageLoader) LoadImgFromFile(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return decodeImage(file, filepath.Ext(path))
}

// Loads an image from a URL.
func (l *ImageLoader) LoadImgFromURL(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return decodeImage(resp.Body, filepath.Ext(url))
}

// Loads an image from a Base64-encoded string.
func (l *ImageLoader) LoadImgFromBase64(b64 string) (image.Image, error) {
	b, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}

	return l.LoadImgFromBytes(b)
}

// Decodes an image from a reader, using the file extension to determine the format.
func decodeImage(reader io.Reader, ext string) (image.Image, error) {
	switch strings.ToLower(ext) {
	case ".png":
		return png.Decode(reader)
	case ".jpg", ".jpeg":
		return jpeg.Decode(reader)
	default:
		img, _, err := image.Decode(reader)
		if err != nil {
			return nil, err
		}
		return img, nil
	}
}
