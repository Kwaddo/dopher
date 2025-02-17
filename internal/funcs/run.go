package funcs

import (
	DM "doom/internal/constants"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

func GameLoop(renderer *sdl.Renderer, player *Player) {
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
			distance := CastRay(player.X, player.Y, rayAngle)

			// Fix fisheye effect
			distance = distance * math.Cos(rayAngle-player.Angle)

			// Calculate wall height
			wallHeight := (DM.ScreenHeight / distance) * 50
			if wallHeight > DM.ScreenHeight {
				wallHeight = DM.ScreenHeight
			}

			// Draw wall slice
			wallTop := (DM.ScreenHeight - wallHeight) / 2
			wallRect := sdl.Rect{
				X: int32(i * (DM.ScreenWidth / DM.NumRays)),
				Y: int32(wallTop),
				W: int32(DM.ScreenWidth/DM.NumRays + 1),
				H: int32(wallHeight),
			}

			// Color based on distance
			intensity := uint8(255 - math.Min(255, distance/2))
			renderer.SetDrawColor(intensity, intensity/2, intensity/2, 255)
			renderer.FillRect(&wallRect)

			rayAngle += DM.FOV / float64(DM.NumRays)
		}

		renderer.Present()
		sdl.Delay(16) // ~60 FPS
	}
}
