package ui

import (
	Casts "doom/internal/graphics/casting"
	DM "doom/internal/model"
	"fmt"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type OptionsMenu DM.OptionsMenu

var GlobalOptionsMenu *OptionsMenu

// NewOptionsMenu creates a new options menu
func NewOptionsMenu() *OptionsMenu {
	menu := &OptionsMenu{
		CurrentOption: 0,
		Options: []string{
			"FOV",
			"Mouse Sensitivity",
			"Movement Speed",
			"Head Bobbing",
			"Back",
		},
		Settings: make(map[string]*DM.Setting),
	}
	menu.Settings["FOV"] = &DM.Setting{
		Value:       DM.FOV * 180 / math.Pi,
		Min:         60,
		Max:         110,
		Step:        5,
		DisplayFunc: func(v float64) string { return fmt.Sprintf("%.0f°", v) },
	}

	menu.Settings["Mouse Sensitivity"] = &DM.Setting{
		Value:       DM.RotateSpeed * 100,
		Min:         5,
		Max:         20,
		Step:        1,
		DisplayFunc: func(v float64) string { return fmt.Sprintf("%.0f", v) },
	}
	menu.Settings["Movement Speed"] = &DM.Setting{
		Value:       DM.BaseMaxSpeed,
		Min:         3,
		Max:         10,
		Step:        0.5,
		DisplayFunc: func(v float64) string { return fmt.Sprintf("%.1f", v) },
	}
	menu.Settings["Head Bobbing"] = &DM.Setting{
		Value: 1,
		Min:   0,
		Max:   1,
		Step:  1,
		DisplayFunc: func(v float64) string {
			if v > 0.5 {
				return "ON"
			}
			return "OFF"
		},
	}
	return menu
}

// MoveUp moves the selection up
func (om *OptionsMenu) MoveUp() {
	om.CurrentOption--
	if om.CurrentOption < 0 {
		om.CurrentOption = len(om.Options) - 1
	}
}

// MoveDown moves the selection down
func (om *OptionsMenu) MoveDown() {
	om.CurrentOption = (om.CurrentOption + 1) % len(om.Options)
}

// GetSelectedOption returns the currently selected option
func (om *OptionsMenu) GetSelectedOption() string {
	return om.Options[om.CurrentOption]
}

// IncreaseSetting increases the current setting value
func (om *OptionsMenu) IncreaseSetting() {
	option := om.GetSelectedOption()
	if option == "Back" {
		return
	}

	setting := om.Settings[option]
	setting.Value = math.Min(setting.Value+setting.Step, setting.Max)
	om.ApplySettings()
}

// DecreaseSetting decreases the current setting value
func (om *OptionsMenu) DecreaseSetting() {
	option := om.GetSelectedOption()
	if option == "Back" {
		return
	}
	setting := om.Settings[option]
	setting.Value = math.Max(setting.Value-setting.Step, setting.Min)
	om.ApplySettings()
}

// ApplySettings applies the current settings to the game
func (om *OptionsMenu) ApplySettings() {
	fovSetting := om.Settings["FOV"].Value
	DM.FOV = fovSetting * math.Pi / 180
	DM.TargetFOV = DM.FOV
	DM.CurrentFOV = DM.FOV
	DM.DynamicFOV = DM.FOV
	DM.RotateSpeed = om.Settings["Mouse Sensitivity"].Value / 100
	DM.BaseMaxSpeed = om.Settings["Movement Speed"].Value
	DM.HeadBobbingEnabled = om.Settings["Head Bobbing"].Value > 0.5
}

// RenderOptionsMenu renders the options menu.
func RenderOptionsMenu(renderer *sdl.Renderer) {
	if GlobalOptionsMenu == nil {
		GlobalOptionsMenu = NewOptionsMenu()
	}
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()
	for y := 0; y < int(DM.ScreenHeight); y++ {
		intensity := uint8(30 - (float64(y) / DM.ScreenHeight * 20))
		renderer.SetDrawColor(intensity, intensity+10, intensity+5, 255)
		renderer.DrawLine(0, int32(y), int32(DM.ScreenWidth), int32(y))
	}
	titleFont, err := Casts.GlobalFontManager.GetFont(36)
	if err == nil {
		RenderText(renderer, titleFont, "OPTIONS", int32(DM.ScreenWidth)/2, int32(DM.ScreenHeight)/6, true)
	}
	optionFont, err := Casts.GlobalFontManager.GetFont(24)
	if err != nil {
		return
	}
	baseY := int32(DM.ScreenHeight) / 3
	for i, option := range GlobalOptionsMenu.Options {
		color := sdl.Color{R: 180, G: 180, B: 180, A: 255}
		if i == GlobalOptionsMenu.CurrentOption {
			color = sdl.Color{R: 255, G: 215, B: 0, A: 255}
		}
		if option == "Back" {
			RenderColoredText(renderer, optionFont, option, int32(DM.ScreenWidth)/2,
				baseY+int32(i*50)+50, true, color)
		} else {
			setting := GlobalOptionsMenu.Settings[option]
			displayValue := setting.DisplayFunc(setting.Value)
			RenderColoredText(renderer, optionFont, option, int32(DM.ScreenWidth)/2-200,
				baseY+int32(i*50), false, color)
			valueText := fmt.Sprintf("< %s >", displayValue)
			RenderColoredText(renderer, optionFont, valueText, int32(DM.ScreenWidth)/2+200,
				baseY+int32(i*50), false, color)
		}
	}
	hintFont, err := Casts.GlobalFontManager.GetFont(16)
	if err == nil {
		RenderText(renderer, hintFont, "Up/Down: Select   Left/Right: Adjust   Enter: Confirm",
			int32(DM.ScreenWidth)/2, int32(DM.ScreenHeight)*5/6, true)
	}
}
