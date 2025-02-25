package casting

import (
	DM "doom/internal/model"
	"math"
)

func CastRay(startX, startY, angle float64) *DM.RayHit {
	// Normalize angle
	for angle < 0 {
		angle += 2 * math.Pi
	}
	for angle >= 2*math.Pi {
		angle -= 2 * math.Pi
	}

	// Precalculate ray direction
	rayX := math.Cos(angle)
	rayY := math.Sin(angle)

	// Current map position
	mapX := int(startX / 100)
	mapY := int(startY / 100)

	// Calculate delta distance
	deltaDistX := math.Abs(1 / rayX)
	deltaDistY := math.Abs(1 / rayY)

	// Calculate step and initial side distance
	var stepX, stepY int
	var sideDistX, sideDistY float64

	if rayX < 0 {
		stepX = -1
		sideDistX = (startX/100 - float64(mapX)) * deltaDistX * 100
	} else {
		stepX = 1
		sideDistX = (float64(mapX) + 1.0 - startX/100) * deltaDistX * 100
	}

	if rayY < 0 {
		stepY = -1
		sideDistY = (startY/100 - float64(mapY)) * deltaDistY * 100
	} else {
		stepY = 1
		sideDistY = (float64(mapY) + 1.0 - startY/100) * deltaDistY * 100
	}

	// DDA algorithm
	var isVertical bool
	WM := DM.GlobalMap.WorldMap

	for {
		// Jump to next map square
		if sideDistX < sideDistY {
			sideDistX += deltaDistX * 100
			mapX += stepX
			isVertical = true
		} else {
			sideDistY += deltaDistY * 100
			mapY += stepY
			isVertical = false
		}

		// Check if ray has hit a wall
		if mapX < 0 || mapX >= len(WM[0]) || mapY < 0 || mapY >= len(WM) {
			return &DM.RayHit{Distance: DM.MaxDepth, WallType: 0}
		}

		if WM[mapY][mapX] > 0 {
			var distance float64
			if isVertical {
				distance = (float64(mapX) - startX/100 + (1-float64(stepX))/2) / rayX * 100
			} else {
				distance = (float64(mapY) - startY/100 + (1-float64(stepY))/2) / rayY * 100
			}

			return &DM.RayHit{
				Distance:   distance,
				WallType:   WM[mapY][mapX],
				HitPointX:  startX + rayX*distance,
				HitPointY:  startY + rayY*distance,
				IsVertical: isVertical,
			}
		}
	}
}

func CalculateDarkness(distance float64) uint8 {
	return uint8(math.Min(DM.MaxDarkness, math.Max(0, distance/3)))
}

func CalculateVerticalDarkness(distance float64) float64 {
	return distance / (float64(DM.ScreenHeight))
}
