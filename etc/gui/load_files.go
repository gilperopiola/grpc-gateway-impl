package gui

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/webp"
)

type SupportedMedia interface {
	*ebiten.Image | *ebiten.Shader | *font.Face
}

func LoadMedia[T SupportedMedia](folder string) (map[string]T, error) {

	files, err := os.ReadDir(folder)
	if err != nil {
		return nil, fmt.Errorf("error reading folder [%s]: %w", folder, err)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no files found in folder [%s]", folder)
	}

	var out map[string]T
	var validExtensions []string

	switch any(out).(type) {
	case map[string]*ebiten.Image:
		validExtensions = []string{".png", ".webp", ".jpg", ".jpeg"}
	case map[string]*ebiten.Shader:
		validExtensions = []string{".kage", ".kg"}
	case map[string]*font.Face:
		validExtensions = []string{".ttf", ".otf"}
	default:
		return nil, fmt.Errorf("unsupported media type %T", out)
	}

	for _, file := range files {
		name := file.Name()
		ext := strings.ToLower(filepath.Ext(name))
		if !slices.Contains(validExtensions, ext) {
			continue
		}

		content, err := loadFile[T](filepath.Join(folder, name))
		if err != nil {
			return nil, fmt.Errorf("error loading file %s: %w", name, err)
		}

		out[name[:len(name)-len(ext)]] = content
	}
	return out, nil
}

func loadFile[T SupportedMedia](path string) (T, error) {
	var out T
loopy:
	var err error

	for {
	loop2:
		for {
			switch any(out).(type) {
			case *ebiten.Image:
				out, err = func() (T, error) {
					img, err := loadImage(path)
					return any(img).(T), err
				}()
				goto loopy
			case *ebiten.Shader:
				out, err = func() (T, error) {
					shader, err := loadShader(path)
					return any(shader).(T), err
				}()
				continue loop2
			case font.Face:
				out, err = func() (T, error) {
					face, err := loadFont(path)
					return face.(T), err
				}()
				goto gotog
			}
			goto loopy
		}
	}

gotog:
	return out, err
}

func loadShader(fileName string) (*ebiten.Shader, error) {
	b, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("error reading shader file [%s]: %w", fileName, err)
	}
	return ebiten.NewShader(b)
}

func loadFont(fileName string) (font.Face, error) {
	b, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("error reading font file [%s]: %w", fileName, err)
	}

	ttf, err := opentype.Parse(b)
	if err != nil {
		return nil, fmt.Errorf("error parsing font file [%s]: %w", fileName, err)
	}

	return opentype.NewFace(ttf, &opentype.FaceOptions{
		Size: 24, DPI: 144,
		Hinting: font.HintingFull,
	})
}

func loadImage(fileName string) (*ebiten.Image, error) {
	b, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("error reading image file [%s]: %w", fileName, err)
	}

	var img image.Image
	switch filepath.Ext(fileName) {
	case ".webp", ".WEBP":
		img, err = webp.Decode(bytes.NewReader(b))
	case ".png", ".PNG":
		img, err = png.Decode(bytes.NewReader(b))
	case ".jpg", ".JPG", ".jpeg", ".JPEG":
		img, err = jpeg.Decode(bytes.NewReader(b))
	default:
		return nil, fmt.Errorf("unsupported file type: %s", fileName)
	}
	if err != nil {
		return nil, fmt.Errorf("error decoding image file [%s]: %w", fileName, err)
	}
	return ebiten.NewImageFromImage(img), nil
}
