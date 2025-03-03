package renders

import (
	MC "doom/internal/char/player"
	Casts "doom/internal/graphics/casting"
	DM "doom/internal/model"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

// RenderSlices renders all the slices of the screen using a go routine.
func RenderSlices(player *MC.Player, dynamicFOV float64, renderChan chan []*DM.RenderSlice) {
	slices := make([]*DM.RenderSlice, DM.NumRays)
	rayAngleStep := dynamicFOV / float64(DM.NumRays-1)
	rayAngle := player.Angle - dynamicFOV/2
	rayWidth := DM.ScreenWidth / float64(DM.NumRays)

	eyeLevel := DM.ScreenHeight / 2
	if player.Crouching {
		heightRatio := player.Height / player.DefaultHeight
		eyeOffset := DM.ScreenHeight * (1 - heightRatio) * 0.1
		eyeLevel -= float64(eyeOffset)
	}
	for i := 0; i < DM.NumRays; i++ {
		rayResult := Casts.CastRay(player.X, player.Y, rayAngle)
		distance := rayResult.Distance * math.Cos(rayAngle-player.Angle)

		wallHeight := (DM.ScreenHeight / distance) * (player.DefaultHeight)
		if wallHeight > DM.ScreenHeight*2 {
			wallHeight = DM.ScreenHeight * 2
		}
		darkness := Casts.CalculateDarkness(distance)
		wallTop := (float64(eyeLevel) - wallHeight/2) + player.BobOffset
		var texCoord int32
		if rayResult.IsVertical {
			texCoord = int32(math.Mod(rayResult.HitPointY, 100) * 0.64)
		} else {
			texCoord = int32(math.Mod(rayResult.HitPointX, 100) * 0.64)
		}
		slices[i] = &DM.RenderSlice{
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
