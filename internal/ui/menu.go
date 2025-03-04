package ui

import (
	Casts "doom/internal/graphics/casting"
	DM "doom/internal/model"

	"github.com/veandco/go-sdl2/sdl"
)

type MainMenu DM.MainMenu

var GlobalMainMenu *MainMenu

// NewMainMenu creates a new main menu
func NewMainMenu() *MainMenu {
	return &MainMenu{
		CurrentOption: 0,
		Options:       []string{"Start Game", "Options", "Quit"},
	}
}

// MoveUp moves the selection up
func (mm *MainMenu) MoveUp() {
	mm.CurrentOption--
	if mm.CurrentOption < 0 {
		mm.CurrentOption = len(mm.Options) - 1
	}
}

// MoveDown moves the selection down
func (mm *MainMenu) MoveDown() {
	mm.CurrentOption = (mm.CurrentOption + 1) % len(mm.Options)
}

// GetSelectedOption returns the currently selected option
func (mm *MainMenu) GetSelectedOption() string {
	return mm.Options[mm.CurrentOption]
}

// RenderMainMenu renders the main menu
func RenderMainMenu(renderer *sdl.Renderer) {
	if GlobalMainMenu == nil {
		GlobalMainMenu = NewMainMenu()
	}
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()
	for y := 0; y < int(DM.ScreenHeight); y++ {
		intensity := uint8(40 - (float64(y) / DM.ScreenHeight * 30))
		renderer.SetDrawColor(intensity, intensity, intensity+10, 255)
		renderer.DrawLine(0, int32(y), int32(DM.ScreenWidth), int32(y))
	}
	titleFont, err := Casts.GlobalFontManager.GetFont(48)
	if err == nil {
		RenderText(renderer, titleFont, "DOPHER", int32(DM.ScreenWidth)/2, int32(DM.ScreenHeight)/4, true)
		subtitleFont, err := Casts.GlobalFontManager.GetFont(24)
		if err == nil {
			RenderText(renderer, subtitleFont, "A Doom-style Engine in Go",
				int32(DM.ScreenWidth)/2, int32(DM.ScreenHeight)/4+60, true)
		}
	}
	optionFont, err := Casts.GlobalFontManager.GetFont(30)
	if err != nil {
		return
	}
	baseY := int32(DM.ScreenHeight) / 2
	for i, option := range GlobalMainMenu.Options {
		color := sdl.Color{R: 180, G: 180, B: 180, A: 255}
		if i == GlobalMainMenu.CurrentOption {
			color = sdl.Color{R: 255, G: 215, B: 0, A: 255}
		}

		RenderColoredText(renderer, optionFont, option, int32(DM.ScreenWidth)/2,
			baseY+int32(i*60), true, color)
	}
	hintFont, err := Casts.GlobalFontManager.GetFont(16)
	if err == nil {
		RenderText(renderer, hintFont, "Up/Down: Navigate   Enter: Select",
			int32(DM.ScreenWidth)/2, int32(DM.ScreenHeight)*4/5, true)
	}
}
