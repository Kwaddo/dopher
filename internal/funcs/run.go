package funcs

import (
	DM "doom/internal/constants"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

func GameLoop(renderer *sdl.Renderer, player *Player) {
	// Load wall texture
	texture, err := loadTexture(renderer)
	if err != nil {
		panic(err)
	}
	defer texture.Destroy()

	running := true
	for running {
		// Handle events
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}

		// Handle keyboard input
		keys := sdl.GetKeyboardState()

		// Rotate with left/right arrows
		if keys[sdl.SCANCODE_LEFT] == 1 {
			player.Angle -= DM.RotateSpeed
		}
		if keys[sdl.SCANCODE_RIGHT] == 1 {
			player.Angle += DM.RotateSpeed
		}

		player.UpdateMovement(keys)

		// When quitting
		if keys[sdl.SCANCODE_ESCAPE] == 1 || keys[sdl.SCANCODE_Q] == 1 {
			running = false
		}

		// Render Scene
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		// Draw 3D view
		rayAngle := player.Angle - DM.FOV/2
		for i := 0; i < DM.NumRays; i++ {
			rayResult := CastRay(player.X, player.Y, rayAngle)

			// Fix fisheye effect
			distance := rayResult.Distance * math.Cos(rayAngle-player.Angle)

			// Calculate wall height
			wallHeight := (DM.ScreenHeight / distance) * 50
			if wallHeight > DM.ScreenHeight {
				wallHeight = DM.ScreenHeight
			}

			// Calculate darkness based on distance
			darkness := uint8(math.Min(255, math.Max(0, distance/2.5)))
			renderer.SetDrawColor(255-darkness, 255-darkness, 255-darkness, 255)

			// Draw wall slice
			wallTop := (DM.ScreenHeight - wallHeight) / 2

			// Calculate texture X coordinate based on exact hit position
			textureX := rayResult.TextureX

			srcRect := &sdl.Rect{
				X: textureX,
				Y: 0,
				W: 1,
				H: 64, // assuming texture height is 64 pixels
			}

			dstRect := &sdl.Rect{
				X: int32(i * (DM.ScreenWidth / DM.NumRays)),
				Y: int32(wallTop),
				W: int32(DM.ScreenWidth/DM.NumRays + 1),
				H: int32(wallHeight),
			}

			// Set the texture color modulation based on distance
			texture.SetColorMod(255-darkness, 255-darkness, 255-darkness)
			renderer.Copy(texture, srcRect, dstRect)

			rayAngle += DM.FOV / float64(DM.NumRays)
		}

		renderer.Present()
		sdl.Delay(16) // ~60 FPS
	}
}
