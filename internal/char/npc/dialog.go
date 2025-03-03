package npc

import (
	DM "doom/internal/model"

	"github.com/veandco/go-sdl2/sdl"
)

// DialogRenderer is a renderer for dialog boxes.
type DialogRenderer DM.DialogRenderer

// NewDialogRenderer creates a new DialogRenderer.
func NewDialogRenderer() (*DialogRenderer, error) {
	if err := DM.InitFonts(); err != nil {
		return nil, err
	}
	return &DialogRenderer{
		Loaded: true,
	}, nil
}

// Close closes the DialogRenderer.
func (dr *DialogRenderer) Close() {
	dr.Loaded = false
}

// RenderDialog renders a dialog box with the given text.
func (dr *DialogRenderer) RenderDialog(renderer *sdl.Renderer, text string) error {
	font, err := DM.GlobalFontManager.GetFont(24)
	if err != nil {
		return err
	}
	surface, err := font.RenderUTF8Solid(
		text,
		sdl.Color{R: 255, G: 255, B: 255, A: 255},
	)
	if err != nil {
		return err
	}
	defer surface.Free()
	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return err
	}
	defer texture.Destroy()
	textW := surface.W
	textH := surface.H
	return renderer.Copy(texture, nil, &sdl.Rect{
		X: int32(int32(DM.ScreenWidth/2) - textW/2),
		Y: int32(int32(DM.ScreenHeight) - textH - 20),
		W: int32(textW),
		H: int32(textH),
	})
}
