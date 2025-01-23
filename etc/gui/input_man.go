package gui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func NewInputMan(onW func(), onA func(), onS func(), onD func(),
	onArrowUp func(), onArrowRight func(), onArrowDown func(), onArrowLeft func()) *InputMan {
	return &InputMan{
		OnW:          onW,
		OnA:          onA,
		OnS:          onS,
		OnD:          onD,
		OnArrowUp:    onArrowUp,
		OnArrowRight: onArrowRight,
		OnArrowDown:  onArrowDown,
		OnArrowLeft:  onArrowLeft,
	}
}

type InputMan struct {
	OnW     func()
	OnWHold func()
	OnA     func()
	OnAHold func()
	OnS     func()
	OnSHold func()
	OnD     func()
	OnDHold func()

	OnArrowUp        func()
	OnArrowUpHold    func()
	OnArrowRight     func()
	OnArrowRightHold func()
	OnArrowDown      func()
	OnArrowDownHold  func()
	OnArrowLeft      func()
	OnArrowLeftHold  func()

	OnLeftClick       func()
	OnLeftClickHold   func()
	OnRightClick      func()
	OnRightClickHold  func()
	OnMiddleClick     func()
	OnMiddleClickHold func()
	OnScrollUp        func()
	OnScrollDown      func()

	OnSpace     func()
	OnSpaceHold func()
	OnEnter     func()
	OnEnterHold func()
	OnEsc       func()
	OnEscHold   func()

	gui *GUI
}

func (this *InputMan) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		if this.OnW != nil {
			this.OnW()
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		if this.OnWHold != nil {
			this.OnWHold()
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		if this.OnA != nil {
			this.OnA()
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		if this.OnAHold != nil {
			this.OnAHold()
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		if this.OnS != nil {
			this.OnS()
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		if this.OnSHold != nil {
			this.OnSHold()
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		if this.OnD != nil {
			this.OnD()
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		if this.OnDHold != nil {
			this.OnDHold()
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		if this.OnArrowUp != nil {
			this.OnArrowUp()
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		if this.OnArrowUpHold != nil {
			this.OnArrowUpHold()
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		if this.OnArrowRight != nil {
			this.OnArrowRight()
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		if this.OnArrowRightHold != nil {
			this.OnArrowRightHold()
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		if this.OnArrowDown != nil {
			this.OnArrowDown()
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		if this.OnArrowDownHold != nil {
			this.OnArrowDownHold()
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		if this.OnArrowLeft != nil {
			this.OnArrowLeft()
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		if this.OnArrowLeftHold != nil {
			this.OnArrowLeftHold()
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if this.OnLeftClick != nil {
			this.OnLeftClick()
		}
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if this.OnLeftClickHold != nil {
			this.OnLeftClickHold()
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		if this.OnRightClick != nil {
			this.OnRightClick()
		}
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		if this.OnRightClickHold != nil {
			this.OnRightClickHold()
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonMiddle) {
		if this.OnMiddleClick != nil {
			this.OnMiddleClick()
		}
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
		if this.OnMiddleClickHold != nil {
			this.OnMiddleClickHold()
		}
	}
}
