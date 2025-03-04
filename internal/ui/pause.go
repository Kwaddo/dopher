package ui

import (
	Casts "doom/internal/graphics/casting"
	DM "doom/internal/model"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type PauseMenu DM.PauseMenu

var GlobalPauseMenu *PauseMenu

// NewPauseMenu creates a new pause menu, as per global struct requirement.
func NewPauseMenu() *PauseMenu {
	return &PauseMenu{
		CurrentOption: 0,
		Options:       []string{"Resume", "Return to Menu", "Quit"},
	}
}

// MoveUp moves the selection up.
func (pm *PauseMenu) MoveUp() {
	pm.CurrentOption--
	if pm.CurrentOption < 0 {
		pm.CurrentOption = len(pm.Options) - 1
	}
}

// MoveDown moves the selection down.
func (pm *PauseMenu) MoveDown() {
	pm.CurrentOption = (pm.CurrentOption + 1) % len(pm.Options)
}

// GetSelectedOption returns the currently selected option.
func (pm *PauseMenu) GetSelectedOption() string {
	return pm.Options[pm.CurrentOption]
}

// RenderPauseMenu renders the pause menu overlay that you see and love!
func RenderPauseMenu(renderer *sdl.Renderer) {
	if GlobalPauseMenu == nil {
		GlobalPauseMenu = &PauseMenu{
			CurrentOption: 0,
			Options:       []string{"Resume", "Return to Menu", "Quit"},
		}
	}
	renderer.SetDrawColor(0, 0, 0, 200)
	renderer.FillRect(&sdl.Rect{
		X: 0,
		Y: 0,
		W: int32(DM.ScreenWidth),
		H: int32(DM.ScreenHeight),
	})
	titleFont, err := Casts.GlobalFontManager.GetFont(36)
	if err == nil {
		RenderText(renderer, titleFont, "PAUSED", int32(DM.ScreenWidth)/2, int32(DM.ScreenHeight)/3, true)
	}
	optionFont, err := Casts.GlobalFontManager.GetFont(24)
	if err != nil {
		return
	}
	baseY := int32(DM.ScreenHeight) / 2
	for i, option := range GlobalPauseMenu.Options {
		color := sdl.Color{R: 200, G: 200, B: 200, A: 255}
		if i == GlobalPauseMenu.CurrentOption {
			color = sdl.Color{R: 255, G: 255, B: 0, A: 255}
		}

		RenderColoredText(renderer, optionFont, option, int32(DM.ScreenWidth)/2,
			baseY+int32(i*50), true, color)
	}
	hintFont, err := Casts.GlobalFontManager.GetFont(16)
	if err == nil {
		RenderText(renderer, hintFont, "Up/Down: Select   Enter: Confirm",
			int32(DM.ScreenWidth)/2, int32(DM.ScreenHeight)*4/5, true)
	}
}

// Helper function to render text
func RenderText(renderer *sdl.Renderer, font *ttf.Font, text string, x, y int32, centered bool) {
	RenderColoredText(renderer, font, text, x, y, centered, sdl.Color{R: 255, G: 255, B: 255, A: 255})
}

// Helper function to render text with specific color
func RenderColoredText(renderer *sdl.Renderer, font *ttf.Font, text string, x, y int32, centered bool, color sdl.Color) {
	surface, err := font.RenderUTF8Solid(text, color)
	if err != nil {
		return
	}
	defer surface.Free()
	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return
	}
	defer texture.Destroy()
	rect := &sdl.Rect{
		W: surface.W,
		H: surface.H,
	}
	if centered {
		rect.X = x - surface.W/2
		rect.Y = y - surface.H/2
	} else {
		rect.X = x
		rect.Y = y
	}
	renderer.Copy(texture, nil, rect)
}
