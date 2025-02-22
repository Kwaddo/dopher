package renders

import (
	Casts "doom/internal/graphics/casting"
	DM "doom/internal/model"
	MC "doom/internal/char/player"
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

type RenderSlice DM.RenderSlice

func RenderSlices(player *MC.Player, DynamicFOV float64, renderChan chan<- []*RenderSlice) {
	slices := make([]*RenderSlice, DM.NumRays)

	// Adjust ray angle calculation to ensure full coverage
	rayAngleStep := DynamicFOV / float64(DM.NumRays-1) // Subtract 1 to include last ray
	rayAngle := player.Angle - DynamicFOV/2

	// Calculate ray width to ensure full screen coverage
	rayWidth := DM.ScreenWidth / float64(DM.NumRays)

	for i := 0; i < DM.NumRays; i++ {
		rayResult := Casts.CastRay(player.X, player.Y, rayAngle)

		// Fix fisheye effect
		distance := rayResult.Distance * math.Cos(rayAngle-player.Angle)

		// Adjust wall height calculation
		wallHeight := (DM.ScreenHeight / distance) * 75
		if wallHeight > DM.ScreenHeight {
			wallHeight = DM.ScreenHeight
		}

		darkness := Casts.CalculateDarkness(distance)
		// Apply head bobbing to wall position
		wallTop := (DM.ScreenHeight-wallHeight)/2 + player.BobOffset

		// Improved texture coordinate calculation
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
