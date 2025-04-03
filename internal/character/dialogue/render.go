package dialogue

import (
	DM "doom/internal/global"
	Casts "doom/internal/graphics/casting"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
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

// SplitToWrappedLines splits text into lines that fit within maxWidth and returns at most maxLines.
func SplitToWrappedLines(text string, font *ttf.Font, maxWidth int, maxLines int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}
	var lines []string
	var currentLine string
	for _, word := range words {
		testLine := currentLine
		if len(testLine) > 0 {
			testLine += " "
		}
		testLine += word
		width, _, err := font.SizeUTF8(testLine)
		if err != nil || width > maxWidth {
			if len(currentLine) > 0 {
				lines = append(lines, currentLine)
				if len(lines) >= maxLines-1 {
					lastLine := word
					for i := 1; i < len(words)-len(lines); i++ {
						if i+len(lines) < len(words) {
							testLastLine := lastLine + " " + words[i+len(lines)]
							width, _, _ := font.SizeUTF8(testLastLine)
							if width > maxWidth-20 {
								break
							}
							lastLine = testLastLine
						}
					}
					if len(lines)+1 < len(words) {
						lastLine += "..."
					}

					lines = append(lines, lastLine)
					break
				}
				currentLine = word
			} else {
				lines = append(lines, word)
				if len(lines) >= maxLines {
					break
				}
				currentLine = ""
			}
		} else {
			currentLine = testLine
		}
	}
	if len(currentLine) > 0 && len(lines) < maxLines {
		lines = append(lines, currentLine)
	}
	return lines
}

// RenderDialogue renders a dialogue box with the given text.
func (dr *DialogueRenderer) RenderDialogue(renderer *sdl.Renderer, text string, charsToShow int) error {
	visibleText := text
	if charsToShow < len(text) {
		visibleText = text[:charsToShow]
	}
	boxWidth := int32(DM.ScreenWidth * 0.8)
	boxHeight := int32(150)
	boxX := int32(DM.ScreenWidth)/2 - boxWidth/2
	boxY := int32(DM.ScreenHeight) - boxHeight - 20
	boxRect := &sdl.Rect{
		X: boxX,
		Y: boxY,
		W: boxWidth,
		H: boxHeight,
	}
	renderer.SetDrawColor(0, 0, 0, 220)
	renderer.FillRect(boxRect)
	renderer.SetDrawColor(180, 180, 180, 255)
	renderer.DrawRect(boxRect)
	font, err := Casts.GlobalFontManager.GetFont(24)
	if err != nil {
		return err
	}
	maxWidth := int(boxWidth) - 40
	lines := SplitToWrappedLines(visibleText, font, maxWidth, 3)
	for i, line := range lines {
		surface, err := font.RenderUTF8Solid(line, sdl.Color{R: 255, G: 255, B: 255, A: 255})
		if err != nil {
			return err
		}
		defer surface.Free()
		texture, err := renderer.CreateTextureFromSurface(surface)
		if err != nil {
			return err
		}
		defer texture.Destroy()
		textY := boxY + 20 + int32(i*40)
		renderer.Copy(texture, nil, &sdl.Rect{
			X: boxX + 20,
			Y: textY,
			W: surface.W,
			H: surface.H,
		})
	}

	return nil
}

// RenderSimpleDialogue renders a simple floating text for NPCs without dialogue trees.
func (dr *DialogueRenderer) RenderSimpleDialogue(renderer *sdl.Renderer, text string) error {
	font, err := Casts.GlobalFontManager.GetFont(20)
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
	defer texture.Destroy()
	posX := int32(DM.ScreenWidth/2) - surface.W/2
	posY := int32(DM.ScreenHeight - 100)
	padding := int32(8)
	bgRect := &sdl.Rect{
		X: posX - padding,
		Y: posY - padding,
		W: surface.W + padding*2,
		H: surface.H + padding*2,
	}
	renderer.SetDrawColor(0, 0, 0, 180)
	renderer.FillRect(bgRect)
	renderer.Copy(texture, nil, &sdl.Rect{
		X: posX,
		Y: posY,
		W: surface.W,
		H: surface.H,
	})
	return nil
}

// RenderDialogueWithOptions renders a dialogue box with a continue prompt.
func (dr *DialogueRenderer) RenderDialogueWithOptions(renderer *sdl.Renderer, npc *DM.NPC) error {
	dialogueTree := npc.DialogueTree
	if dialogueTree == nil {
		return nil
	}
	err := dr.RenderDialogue(renderer, npc.DialogText, dialogueTree.CharsToShow)
	if err != nil {
		return err
	}
	if npc.DialogueTree.TextFullyShown {
		currentNode := npc.DialogueTree.Nodes[npc.DialogueTree.CurrentNodeID]
		if currentNode == nil {
			return nil
		}
		boxWidth := int32(DM.ScreenWidth * 0.8)
		boxHeight := int32(150)
		boxX := int32(DM.ScreenWidth)/2 - boxWidth/2
		boxY := int32(DM.ScreenHeight) - boxHeight - 20
		hintFont, err := Casts.GlobalFontManager.GetFont(16)
		if err == nil {
			hintText := "Press E to continue, Enter to exit"
			surface, err := hintFont.RenderUTF8Solid(hintText, sdl.Color{R: 180, G: 180, B: 180, A: 255})
			if err == nil {
				defer surface.Free()
				texture, err := renderer.CreateTextureFromSurface(surface)
				if err == nil {
					defer texture.Destroy()
					renderer.Copy(texture, nil, &sdl.Rect{
						X: boxX + 20,
						Y: boxY + boxHeight - surface.H - 10,
						W: surface.W,
						H: surface.H,
					})
				}
			}
		}
		blinkOn := (DM.GlobalFrameCount/30)%2 == 0
		if blinkOn {
			arrowFont, err := Casts.GlobalFontManager.GetFont(28)
			if err != nil {
				return err
			}
			surface, err := arrowFont.RenderUTF8Solid(">", sdl.Color{R: 255, G: 255, B: 0, A: 255})
			if err != nil {
				return err
			}
			defer surface.Free()

			texture, err := renderer.CreateTextureFromSurface(surface)
			if err != nil {
				return err
			}
			defer texture.Destroy()
			renderer.Copy(texture, nil, &sdl.Rect{
				X: boxX + boxWidth - surface.W - 20,
				Y: boxY + boxHeight - surface.H - 10,
				W: surface.W,
				H: surface.H,
			})
		}
	}
	return nil
}
