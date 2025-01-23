package gui

import (
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// Thing is anything that could be displayed on a screen
type Thing struct {
	Name       string
	Image      func() *ebiten.Image
	ImageName  string
	Shader     func() *ebiten.Shader
	ShaderName string

	X, Y          int
	Width, Height int
	Scale         float32

	Hovered   bool
	CreatedAt time.Time
	Nuance    float32
	Sin1      float32
}

// Draw said thing on the screen
func (t *Thing) Draw(scr *ebiten.Image, cam *Cam) {
	image := t.Image()
	shader := t.Shader()
	if image == nil || shader == nil {
		return
	}

	t.drawPolaroidFrame(scr, cam)

	opts := &ebiten.DrawRectShaderOptions{
		Images:   [4]*ebiten.Image{image, nil, nil, nil},
		Uniforms: t.makeUniforms(),
	}

	opts.GeoM.Scale(float64(cam.Scale), float64(cam.Scale))

	// Offset by the camera
	realX := float64(float32(t.X) - cam.X)
	realY := float64(float32(t.Y) - cam.Y)
	opts.GeoM.Translate(realX, realY)
	opts.GeoM.Translate(0, float64(t.Sin1*25))

	scr.DrawRectShader(t.Width, t.Height, shader, opts)
}

func (t *Thing) drawPolaroidFrame(scr *ebiten.Image, cam *Cam) {
	realX := float64(float32(t.X) - cam.X)
	realY := float64(float32(t.Y) - cam.Y)

	newimg := ebiten.NewImage(int((float32(t.Width)+24)*cam.Scale), int((float32(t.Height)+256)*cam.Scale))
	newimg.Fill(color.RGBA{R: 222, G: 66, B: 188, A: 255})

	opts := &ebiten.DrawRectShaderOptions{
		Images: [4]*ebiten.Image{newimg, nil, nil, nil},
	}

	opts.GeoM.Scale(float64(cam.Scale), float64(cam.Scale))
	opts.GeoM.Translate(float64(float32(realX)-(12*cam.Scale)), float64(float32(realY)-(12*cam.Scale)))
	opts.GeoM.Translate(0, float64(t.Sin1*25))

	scr.DrawRectShader(int((float32(t.Width)+24)*cam.Scale), int((float32(t.Height)+256)*cam.Scale), t.Shader(), opts)
}

func (this *Thing) makeUniforms() map[string]any {
	var Time float64 = time.Since(this.CreatedAt).Seconds()

	return map[string]any{
		"Time": Time,
	}
}

func (t *Thing) Update() {
	t.Sin1 = float32(math.Sin(float64(time.Since(t.CreatedAt.Add(time.Second * -5 * time.Duration(t.Nuance))).Seconds())))
}
