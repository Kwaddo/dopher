package core

import (
	NPC "doom/internal/char/npc"
	MC "doom/internal/char/player"
	Graphics "doom/internal/graphics"
	DM "doom/internal/model"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

func GameLoop(renderer *sdl.Renderer, player *MC.Player) {
	// Load wall textures
	textures, err := Graphics.LoadTextures(renderer)
	if err != nil {
		panic(err)
	}
	// Clean up textures
	defer func() {
		for _, texture := range textures.Textures {
			texture.Destroy()
		}
	}()

	var DynamicFOV float64
	var currentFOV = DM.FOV
	var targetFOV = DM.FOV
	const lerpSpeed = 0.15

	// Create a channel for render slices
	renderChan := make(chan []*Graphics.RenderSlice, 1)

	// Initialize NPC manager globally
	npcManager := NPC.NewNPCManager()
	NPC.GlobalNPCManager = npcManager // Set the global reference

	// Create a z-buffer to store wall distances
	zBuffer := make([]float64, DM.ScreenWidth)

	// Initialize dialog renderer
	dialogRenderer, err := NPC.NewDialogRenderer()
	if err != nil {
		panic(err)
	}
	defer dialogRenderer.Close()

	for {
		// Handle events
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		// Player controls
		end := player.Movement(sdl.GetKeyboardState(), npcManager)
		if end {
			break
		}

		// Smooth FOV transition
		if player.Running {
			targetFOV = DM.FOV * 0.90
		} else {
			targetFOV = DM.FOV
		}

		currentFOV = MC.LERP(currentFOV, targetFOV, lerpSpeed)
		DynamicFOV = currentFOV

		// Update NPC distances
		npcManager.UpdateDistances(player.X, player.Y)
		npcManager.SortByDistance()

		// Clear z-buffer
		for i := range zBuffer {
			zBuffer[i] = math.MaxFloat64
		}

		// Update NPC dialogs
		npcManager.UpdateDialogs()

		// Render Scene
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		// Set blend mode for proper transparency
		renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)

		// Start goroutine to calculate render slices
		go Graphics.RenderSlices(player, DynamicFOV, renderChan)

		// Receive and render the slices
		renderSlices := <-renderChan
		for _, slice := range renderSlices {
			// Render with texture if available
			if texture, ok := textures.Textures[slice.WallType]; ok {
				// Apply darkness to the texture
				texture.SetColorMod(255-slice.Darkness, 255-slice.Darkness, 255-slice.Darkness)
				srcRect := &sdl.Rect{
					X: slice.TexCoord, // Use the calculated texture coordinate
					Y: 0,
					W: 1, // Use the entire texture width
					H: 64,
				}
				renderer.Copy(texture, srcRect, slice.DstRect)

				// Store wall distance in z-buffer
				screenX := int(slice.DstRect.X)
				for x := screenX; x < screenX+int(slice.DstRect.W) && x < DM.ScreenWidth; x++ {
					if x >= 0 {
						zBuffer[x] = slice.Distance
					}
				}
			}
		}

		// After rendering wall slices and updating z-buffer
		sprites := Graphics.RenderNPCs(player, npcManager, DynamicFOV, zBuffer)
		for _, sprite := range sprites {
			if texture, ok := textures.Textures[sprite.WallType]; ok {
				texture.SetColorMod(255-sprite.Darkness, 255-sprite.Darkness, 255-sprite.Darkness)
				texture.SetBlendMode(sdl.BLENDMODE_BLEND)

				// Only render visible parts of the sprite
				dstRect := sprite.DstRect
				for x := dstRect.X; x < dstRect.X+dstRect.W; x++ {
					if x >= 0 && x < int32(len(zBuffer)) {
						if sprite.Distance > zBuffer[x] {
							// Skip rendering this column if behind wall
							continue
						}

						// Render visible column
						columnRect := &sdl.Rect{
							X: x,
							Y: dstRect.Y,
							W: 1,
							H: dstRect.H,
						}

						// Calculate source X coordinate for proper sprite column
						srcX := int32(float64(x-dstRect.X) / float64(dstRect.W) * 64)
						srcColumnRect := &sdl.Rect{
							X: srcX,
							Y: 0,
							W: 1,
							H: 64,
						}

						renderer.Copy(texture, srcColumnRect, columnRect)
					}
				}
			}
		}

		for _, npc := range npcManager.NPCs {
			if npc.ShowDialog {
				err := dialogRenderer.RenderDialog(renderer, npc.DialogText)
				if err != nil {
					continue
				}
			}
		}

		renderer.Present()
		sdl.Delay(16) // ~60 FPS
	}
}
