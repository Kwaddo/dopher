package funcs

import (
	DM "doom/internal/constants"
	"math"
)

// New struct to hold ray casting results
type RayHit struct {
	Distance float64
	TextureX int32
}

func CastRay(startX, startY, angle float64) RayHit {
	rayX := math.Cos(angle)
	rayY := math.Sin(angle)
	WM := DM.GlobalMap.WorldMap
	distance := 0.0

	// Use smaller step size for more precise hits
	step := 0.1 // More precise step size

	for distance < DM.MaxDepth {
		x := startX + rayX*distance
		y := startY + rayY*distance

		mapX := int(x / 100)
		mapY := int(y / 100)

		if mapX >= 0 && mapX < len(WM[0]) && mapY >= 0 && mapY < len(WM) {
			if WM[mapY][mapX] == 1 {
				// Calculate exact hit position
				textureX := int32(math.Floor(math.Mod(x, 100.0)) * 0.64) // Scale to texture width

				return RayHit{
					Distance: distance,
					TextureX: textureX,
				}
			}
		}

		distance += step
	}

	return RayHit{
		Distance: distance,
		TextureX: 0,
	}
}
