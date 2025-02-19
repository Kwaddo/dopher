package funcs

import (
	DM "doom/internal/constants"

	"github.com/veandco/go-sdl2/sdl"
)

func GameLoop(renderer *sdl.Renderer, player *Player) {
	// Load wall textures
	textures, err := LoadTextures(renderer)
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
	renderChan := make(chan []RenderSlice, 1)

	for {
		// Handle events
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		// Player controls
		end := player.UpdateMovement(sdl.GetKeyboardState())
		if end {
			break
		}

		// Smooth FOV transition
		if player.Running {
			targetFOV = DM.FOV * 0.95
		} else {
			targetFOV = DM.FOV
		}

		currentFOV = LERP(currentFOV, targetFOV, lerpSpeed)
		DynamicFOV = currentFOV

		// Render Scene
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		// Start goroutine to calculate render slices
		go RenderSlices(player, DynamicFOV, renderChan)

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
					W: 64, // Use the entire texture width
					H: 64,
				}
				renderer.Copy(texture, srcRect, slice.DstRect)
			}
		}

		renderer.Present()
		sdl.Delay(16) // ~60 FPS
	}
}
