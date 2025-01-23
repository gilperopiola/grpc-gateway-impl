package gui

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gilperopiola/god"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/webp"
)

func omitSecond[T any](first T, _ string) T {
	return first
}

type SupportedMedia interface {
	*ebiten.Image | *ebiten.Shader | *font.Face
}

type ContentMan struct {
	Images  ContentCycler[*ebiten.Image]
	Shaders ContentCycler[*ebiten.Shader]
	Fonts   ContentCycler[*font.Face]
}

type ContentCycler[T SupportedMedia] struct {
	Medias       map[string]T
	Names        []string
	CurrentIndex int
}

func NewContentMan() *ContentMan {
	images, err := LoadMediaMap[*ebiten.Image]("etc/gui/images")
	if err != nil {
		log.Fatalf("error loading images: %v", err)
	}

	shaders, err := LoadMediaMap[*ebiten.Shader]("etc/gui/shaders")
	if err != nil {
		log.Fatalf("error loading shaders: %v", err)
	}

	fonts, err := LoadMediaMap[*font.Face]("etc/gui/fonts")
	if err != nil {
		log.Fatalf("error loading fonts: %v", err)
	}

	return &ContentMan{
		Images: ContentCycler[*ebiten.Image]{
			Medias: images,
			Names:  god.GetMapKeys(images),
		},
		Shaders: ContentCycler[*ebiten.Shader]{
			Medias: shaders,
			Names:  god.GetMapKeys(shaders),
		},
		Fonts: ContentCycler[*font.Face]{
			Medias: fonts,
			Names:  god.GetMapKeys(fonts),
		},
	}
}

func (cc *ContentCycler[T]) GetCurrent() (T, string) {
	name := cc.Names[cc.CurrentIndex]
	media := cc.Medias[name]
	return media, name
}

func (cc *ContentCycler[T]) GetNext() (T, string) {
	cc.CurrentIndex++
	if cc.CurrentIndex >= len(cc.Names) {
		cc.CurrentIndex = 0
	}
	return cc.GetCurrent()
}

func (cc *ContentCycler[T]) GetPrev() (T, string) {
	cc.CurrentIndex--
	if cc.CurrentIndex < 0 {
		cc.CurrentIndex = len(cc.Names) - 1
	}
	return cc.GetCurrent()
}

/* - Load Media - */

func LoadMediaMap[T SupportedMedia](folder string) (map[string]T, error) {
	files, err := os.ReadDir(folder)
	if err != nil {
		return nil, fmt.Errorf("error reading folder [%s]: %w", folder, err)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no files found in folder [%s]", folder)
	}

	out := make(map[string]T)
	validExtensions := []string{}

	switch any(out).(type) {
	case map[string]*ebiten.Image:
		validExtensions = []string{".png", ".webp", ".jpg", ".jpeg"}
	case map[string]*ebiten.Shader:
		validExtensions = []string{".kage", ".kg"}
	case map[string]*font.Face:
		validExtensions = []string{".ttf", ".otf"}
	default:
		return nil, fmt.Errorf("unsupported media map type %T", out)
	}

	for c, file := range files {
		name := file.Name()
		ext := strings.ToLower(filepath.Ext(name))
		if !slices.Contains(validExtensions, ext) {
			continue
		}

		if c > 25 {
			break
		}

		content, err := loadMedia[T](filepath.Join(folder, name))
		if err != nil {
			return nil, fmt.Errorf("error loading file %s: %w", name, err)
		}

		out[name[:len(name)-len(ext)]] = content
	}
	return out, nil
}

func loadMedia[T SupportedMedia](path string) (T, error) {
	var out T
	var err error

	switch any(out).(type) {
	case *ebiten.Image:
		out, err = func() (T, error) {
			img, err := loadImage(path)
			return any(img).(T), err
		}()
	case *ebiten.Shader:
		out, err = func() (T, error) {
			shader, err := loadShader(path)
			return any(shader).(T), err
		}()
	case font.Face:
		out, err = func() (T, error) {
			face, err := loadFont(path)
			return face.(T), err
		}()
	}

	return out, err
}

func loadShader(fileName string) (*ebiten.Shader, error) {
	b, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("error reading shader file %s: %w", fileName, err)
	}
	return ebiten.NewShader(b)
}

func loadFont(fileName string) (font.Face, error) {
	b, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("error reading font file %s: %w", fileName, err)
	}

	ttf, err := opentype.Parse(b)
	if err != nil {
		return nil, fmt.Errorf("error parsing font file %s: %w", fileName, err)
	}

	return opentype.NewFace(ttf, &opentype.FaceOptions{
		Size: 24, DPI: 144,
		Hinting: font.HintingFull,
	})
}

func loadImage(fileName string) (*ebiten.Image, error) {
	b, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("error reading image file %s: %w", fileName, err)
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
		return nil, fmt.Errorf("unsupported image file type: %s", fileName)
	}
	if err != nil {
		return nil, fmt.Errorf("error decoding image file %s: %w", fileName, err)
	}
	return ebiten.NewImageFromImage(img), nil
}
