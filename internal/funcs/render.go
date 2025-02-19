package funcs

import (
	DM "doom/internal/constants"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type RenderSlice struct {
	DstRect  *sdl.Rect
	Darkness uint8
	Color    sdl.Color
	WallType int
	TexCoord int32
}

func RenderSlices(player *Player, DynamicFOV float64, renderChan chan<- []RenderSlice) {
	slices := make([]RenderSlice, DM.NumRays)
	rayAngle := player.Angle - DynamicFOV/2

	for i := 0; i < DM.NumRays; i++ {
		rayResult := CastRay(player.X, player.Y, rayAngle)

		// Fix fisheye effect
		distance := rayResult.Distance * math.Cos(rayAngle-player.Angle)

		// Adjust wall height calculation
		wallHeight := (DM.ScreenHeight / distance) * 75 // Increased scaling factor for better visibility
		if wallHeight > DM.ScreenHeight {
			wallHeight = DM.ScreenHeight
		}

		darkness := uint8(math.Min(255, math.Max(0, distance/3)))
		wallTop := (DM.ScreenHeight - wallHeight) / 2

		// Improved texture coordinate calculation with proper scaling
		var texCoord int32
		if rayResult.IsVertical {
			// For vertical walls, use Y coordinate
			texCoord = int32(math.Mod(rayResult.HitPointY, 100) * 0.01)
		} else {
			// For horizontal walls, use X coordinate
			texCoord = int32(math.Mod(rayResult.HitPointX, 100) * 0.01)
		}


		slices[i] = RenderSlice{
			DstRect: &sdl.Rect{
				X: int32(i * (DM.ScreenWidth / DM.NumRays)),
				Y: int32(wallTop),
				W: int32(DM.ScreenWidth/DM.NumRays + 1),
				H: int32(wallHeight),
			},
			Darkness: darkness,
			Color:    sdl.Color{R: 128, G: 128, B: 128, A: 255},
			WallType: rayResult.WallType,
			TexCoord: texCoord,
		}

		rayAngle += DynamicFOV / float64(DM.NumRays)
	}

	renderChan <- slices
}
