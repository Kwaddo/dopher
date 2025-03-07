package npc

import (
	Casts "doom/internal/graphics/casting"
	DM "doom/internal/model"
	"fmt"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type DialogueRenderer DM.DialogueRenderer

// NewDialogueRenderer creates a new DialogueRenderer.
func NewDialogueRenderer() (*DialogueRenderer, error) {
	if err := Casts.InitFonts(); err != nil {
		return nil, err
	}
	return &DialogueRenderer{
		Loaded:    true,
		TextCache: make(map[string]*DM.TextureCacheEntry),
	}, nil
}

// Close closes the DialogRenderer.
func (dr *DialogueRenderer) Close() {
	for _, entry := range dr.TextCache {
		entry.Texture.Destroy()
	}
	dr.TextCache = nil
	dr.Loaded = false
}

// RenderDialogue renders a dialogue box with the given text.
func (dr *DialogueRenderer) RenderDialogue(renderer *sdl.Renderer, text string) error {
	if entry, exists := dr.TextCache[text]; exists {
		entry.LastUsed = int64(DM.GlobalFrameCount)
		return renderer.Copy(entry.Texture, nil, &sdl.Rect{
			X: int32(DM.ScreenWidth/2) - entry.Width/2,
			Y: int32(DM.ScreenHeight) - entry.Height - 20,
			W: entry.Width,
			H: entry.Height,
		})
	}
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
	dr.TextCache[text] = &DM.TextureCacheEntry{
		Texture:  texture,
		Width:    surface.W,
		Height:   surface.H,
		LastUsed: int64(DM.GlobalFrameCount),
	}
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

// RenderDialogueWithOptions renders a dialogue box with player response options.
func (dr *DialogueRenderer) RenderDialogueWithOptions(renderer *sdl.Renderer, npc *DM.NPC) error {
	err := dr.RenderDialogue(renderer, npc.DialogText)
	if err != nil {
		return err
	}
	if npc.DialogueTree == nil || !npc.DialogueTree.IsActive {
		return nil
	}
	currentNode := npc.DialogueTree.Nodes[npc.DialogueTree.CurrentNodeID]
	if currentNode == nil {
		return nil
	}
	font, err := Casts.GlobalFontManager.GetFont(18)
	if err != nil {
		return err
	}
	optionsBgRect := &sdl.Rect{
		X: int32(DM.ScreenWidth/2) - 300,
		Y: int32(DM.ScreenHeight) - int32(30*len(currentNode.Options)) - 80,
		W: 600,
		H: int32(30*len(currentNode.Options)) + 20,
	}
	renderer.SetDrawColor(0, 0, 0, 200)
	renderer.FillRect(optionsBgRect)
	renderer.SetDrawColor(100, 100, 100, 255)
	renderer.DrawRect(optionsBgRect)
	validOptionIndex := 0
	for i, option := range currentNode.Options {
		if option.Condition != nil && !option.Condition(npc) {
			continue
		}
		optionText := fmt.Sprintf("%d. %s", validOptionIndex+1, option.Text)
		validOptionIndex++
		surface, err := font.RenderUTF8Solid(optionText, sdl.Color{R: 200, G: 220, B: 255, A: 255})
		if err != nil {
			continue
		}
		defer surface.Free()
		texture, err := renderer.CreateTextureFromSurface(surface)
		if err != nil {
			continue
		}
		defer texture.Destroy()
		renderer.Copy(texture, nil, &sdl.Rect{
			X: int32(DM.ScreenWidth/2) - surface.W/2,
			Y: int32(DM.ScreenHeight) - int32(30*(len(currentNode.Options)-i)) - 70,
			W: surface.W,
			H: surface.H,
		})
	}
	return nil
}
