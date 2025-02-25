package renders

import (
	MC "doom/internal/char/player"
	Casts "doom/internal/graphics/casting"
	DM "doom/internal/model"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type RenderSlice DM.RenderSlice

func RenderSlices(player *MC.Player, dynamicFOV float64, renderChan chan []*RenderSlice) {
	slices := make([]*RenderSlice, DM.NumRays)

	rayAngleStep := dynamicFOV / float64(DM.NumRays-1)
	rayAngle := player.Angle - dynamicFOV/2
	rayWidth := DM.ScreenWidth / float64(DM.NumRays)

	// Calculate eye level based on crouching state
	eyeLevel := DM.ScreenHeight / 2
	if player.Crouching {
		// Lower the eye level when crouching (subtract instead of add)
		heightRatio := player.Height / player.DefaultHeight
		eyeOffset := DM.ScreenHeight * (1 - heightRatio) * 0.1 // Reduced multiplier for more natural crouch
		eyeLevel -= float64(eyeOffset)                         // Changed from += to -= to make player look down instead of up
	}

	for i := 0; i < DM.NumRays; i++ {
		rayResult := Casts.CastRay(player.X, player.Y, rayAngle)
		distance := rayResult.Distance * math.Cos(rayAngle-player.Angle)

		// Keep using DefaultHeight for consistent wall scaling
		wallHeight := (DM.ScreenHeight / distance) * (player.DefaultHeight)
		if wallHeight > DM.ScreenHeight*2 { // Allow for taller walls to extend beyond screen
			wallHeight = DM.ScreenHeight * 2
		}

		darkness := Casts.CalculateDarkness(distance)

		// Adjusted wall positioning with proper eyeLevel consideration
		wallTop := (float64(eyeLevel) - wallHeight/2) + player.BobOffset

		// Texture coord calculation
		var texCoord int32
		if rayResult.IsVertical {
			texCoord = int32(math.Mod(rayResult.HitPointY, 100) * 0.64)
		} else {
			texCoord = int32(math.Mod(rayResult.HitPointX, 100) * 0.64)
		}

		slices[i] = &RenderSlice{
			DstRect: &sdl.Rect{
				X: int32(float64(i) * rayWidth),
				Y: int32(wallTop),
				W: int32(math.Ceil(rayWidth)),
				H: int32(wallHeight),
			},
			Darkness: darkness,
			Color:    sdl.Color{R: 128, G: 128, B: 128, A: 255},
			WallType: rayResult.WallType,
			TexCoord: texCoord,
			Distance: distance,
		}

		rayAngle += rayAngleStep
	}

	renderChan <- slices
}
