package gui

import (
	"log"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/service"
	"github.com/hajimehoshi/ebiten/v2"
)

type GUI struct {
	Width  int
	Height int

	// UI Stuff
	Cam    Cam
	Things []Thing

	// Managers
	*InputMan
	*ContentMan

	// Dependencies
	service *service.Service
}

type Cam struct {
	X float32
	Y float32

	Step  float32
	Scale float32
}

func Start(service *service.Service) {
	go func() {
		gui := GUI{
			Width:      1600,
			Height:     900,
			service:    service,
			Things:     []Thing{},
			InputMan:   NewInputMan(nil, nil, nil, nil, nil, nil, nil, nil),
			ContentMan: NewContentMan(),
			Cam:        Cam{X: 0, Y: 0, Step: 15, Scale: 0.5},
		}

		gui.InputMan.OnW = func() {
			gui.ContentMan.Shaders.GetNext()
		}
		gui.InputMan.OnS = func() {
			gui.ContentMan.Shaders.GetPrev()
		}
		gui.InputMan.OnD = func() {
			gui.ContentMan.Images.GetNext()
		}
		gui.InputMan.OnA = func() {
			gui.ContentMan.Images.GetPrev()
		}

		gui.InputMan.OnArrowUpHold = func() {
			gui.Cam.Y += gui.Cam.Step
		}
		gui.InputMan.OnArrowDownHold = func() {
			gui.Cam.Y -= gui.Cam.Step
		}
		gui.InputMan.OnArrowLeftHold = func() {
			gui.Cam.X += gui.Cam.Step
		}
		gui.InputMan.OnArrowRightHold = func() {
			gui.Cam.X -= gui.Cam.Step
		}

		gui.InputMan.OnLeftClick = func() {
			gui.Cam.Scale -= 0.04
		}
		gui.InputMan.OnRightClick = func() {
			gui.Cam.Scale += 0.04
		}

		gui.Things = append(gui.Things, Thing{
			Name:      "1",
			Image:     func() *ebiten.Image { return omitSecond(gui.ContentMan.Images.GetCurrent()) },
			Shader:    func() *ebiten.Shader { return omitSecond(gui.ContentMan.Shaders.GetCurrent()) },
			X:         0,
			Y:         0,
			Width:     1024,
			Height:    1024,
			CreatedAt: time.Now(),
			Nuance:    1.2,
		})

		gui.Things = append(gui.Things, Thing{
			Name:      "2",
			Image:     func() *ebiten.Image { return omitSecond(gui.ContentMan.Images.GetCurrent()) },
			Shader:    func() *ebiten.Shader { return omitSecond(gui.ContentMan.Shaders.GetCurrent()) },
			X:         int(float32(gui.Width)/2 - 512),
			Y:         0,
			Width:     1024,
			Height:    1024,
			CreatedAt: time.Now(),
			Nuance:    1.4,
		})
		gui.Things = append(gui.Things, Thing{
			Name:      "3",
			Image:     func() *ebiten.Image { return omitSecond(gui.ContentMan.Images.GetCurrent()) },
			Shader:    func() *ebiten.Shader { return omitSecond(gui.ContentMan.Shaders.GetCurrent()) },
			X:         int(float32(gui.Width) / 1.5),
			Y:         0,
			Width:     1024,
			Height:    1024,
			CreatedAt: time.Now(),
			Nuance:    1.6,
		})

		ebiten.SetWindowSize(gui.Width, gui.Height)
		ebiten.SetWindowTitle("GGWI GUI")
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
		ebiten.SetTPS(60)

		if err := ebiten.RunGame(&gui); err != nil {
			log.Fatalf("Error in game loop >>> %s", err)
		}
	}()
}

func (g *GUI) Draw(scr *ebiten.Image) {
	for _, thing := range g.Things {
		thing.Draw(scr, &g.Cam)
	}
}

func (g *GUI) Update() error {
	g.InputMan.Update()
	for _, thing := range g.Things {
		thing.Update()
	}
	return nil
}

func (gui *GUI) Layout(_, _ int) (int, int) {
	return gui.Width, gui.Height
}
