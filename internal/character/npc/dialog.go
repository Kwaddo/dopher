package npc

import (
	Casts "doom/internal/graphics/casting"
	DM "doom/internal/model"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

// DialogRenderer is a renderer for dialog boxes.
type DialogRenderer struct {
	Loaded    bool
	TextCache map[string]*TextureCacheEntry
}

type TextureCacheEntry struct {
	Texture       *sdl.Texture
	Width, Height int32
	LastUsed      int64 // Timestamp when last used
}

// NewDialogRenderer creates a new DialogRenderer.
func NewDialogRenderer() (*DialogRenderer, error) {
	if err := Casts.InitFonts(); err != nil {
		return nil, err
	}
	return &DialogRenderer{
		Loaded:    true,
		TextCache: make(map[string]*TextureCacheEntry),
	}, nil
}

// Close closes the DialogRenderer.
func (dr *DialogRenderer) Close() {
	for _, entry := range dr.TextCache {
		entry.Texture.Destroy()
	}
	dr.TextCache = nil
	dr.Loaded = false
}

// RenderDialog renders a dialog box with the given text.
func (dr *DialogRenderer) RenderDialog(renderer *sdl.Renderer, text string) error {
	// Check cache for this text
	if entry, exists := dr.TextCache[text]; exists {
		entry.LastUsed = int64(DM.GlobalFrameCount)
		return renderer.Copy(entry.Texture, nil, &sdl.Rect{
			X: int32(DM.ScreenWidth/2) - entry.Width/2,
			Y: int32(DM.ScreenHeight) - entry.Height - 20,
			W: entry.Width,
			H: entry.Height,
		})
	}

	// Create new texture
	font, err := Casts.GlobalFontManager.GetFont(24)
	if err != nil {
		return err
	}

	surface, err := font.RenderUTF8Solid(text, sdl.Color{R: 255, G: 255, B: 255, A: 255})
	if err != nil {
		return err
	}
	defer surface.Free()

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return err
	}

	// Store in cache
	dr.TextCache[text] = &TextureCacheEntry{
		Texture:  texture,
		Width:    surface.W,
		Height:   surface.H,
		LastUsed: int64(DM.GlobalFrameCount),
	}

	// Clean cache if too large (keep only last 20 entries)
	if len(dr.TextCache) > 20 {
		oldest := int64(math.MaxInt64)
		oldestKey := ""
		for k, v := range dr.TextCache {
			if v.LastUsed < oldest {
				oldest = v.LastUsed
				oldestKey = k
			}
		}
		if oldestKey != "" {
			dr.TextCache[oldestKey].Texture.Destroy()
			delete(dr.TextCache, oldestKey)
		}
	}

	return renderer.Copy(texture, nil, &sdl.Rect{
		X: int32(DM.ScreenWidth/2) - surface.W/2,
		Y: int32(DM.ScreenHeight) - surface.H - 20,
		W: surface.W,
		H: surface.H,
	})
}
