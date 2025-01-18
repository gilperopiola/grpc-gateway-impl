package gui

import (
	"log"

	"github.com/gilperopiola/grpc-gateway-impl/app/service"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

type GUI struct {
	Screen  Thing
	Things  []Thing
	Shaders map[string]*ebiten.Shader
	Images  map[string]*ebiten.Image
	Fonts   map[string]*font.Face

	service *service.Service
}

type Thing struct {
	X, Y          int
	Width, Height int
	Image         *ebiten.Image
}

func (g *GUI) Draw(scr *ebiten.Image) {
	for _, thing := range g.Things {
		thing.Draw(scr)
	}
}

func (g *GUI) Update() error {
	for _, thing := range g.Things {
		thing.Update()
	}
	return nil
}

func (t *Thing) Draw(scr *ebiten.Image) {
	if t.Image != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(t.X), float64(t.Y))
		scr.DrawImage(t.Image, op)
	}
}

func (t *Thing) Update() {
}

func Start(service *service.Service) {

	gui := &GUI{
		Screen: Thing{Width: 1792, Height: 1024},
		Things: []Thing{},
		Shaders: func() map[string]*ebiten.Shader {
			shaders, _ := LoadMedia[*ebiten.Shader]("etc/gui/shaders/")
			return shaders
		}(),
		Images: func() map[string]*ebiten.Image {
			images, _ := LoadMedia[*ebiten.Image]("etc/gui/images/")
			return images
		}(),
		Fonts: func() map[string]*font.Face {
			fonts, _ := LoadMedia[*font.Face]("etc/gui/fonts/")
			return fonts
		}(),

		service: service,
	}

	ebiten.SetWindowSize(gui.Screen.Width, gui.Screen.Height)
	ebiten.SetWindowTitle("GGWI GUI")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(60)

	if err := ebiten.RunGame(gui); err != nil {
		log.Fatalf("Error in game loop ——> %s", err)
	}
}

func (g *GUI) Layout(_, _ int) (int, int) {
	return g.Screen.Width, g.Screen.Height
}
