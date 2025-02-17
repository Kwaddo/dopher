package funcs

import (
	"math"
	DM "doom/internal/constants"
)

func CastRay(startX, startY, angle float64) float64 {
	rayX := math.Cos(angle)
	rayY := math.Sin(angle)
	WM := DM.GlobalMap.WorldMap
	distance := 0.0
	for distance < DM.MaxDepth {
		x := startX + rayX*distance
		y := startY + rayY*distance

		mapX := int(x / 100)
		mapY := int(y / 100)

		if mapX >= 0 && mapX < len(WM[0]) && mapY >= 0 && mapY < len(WM) {
			if WM[mapY][mapX] == 1 {
				return distance
			}
		}

		distance += 1.0
	}
	return DM.MaxDepth
}
