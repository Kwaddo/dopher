package casting

import (
	DM "doom/internal/model"
	"math"
)

// CastRay casts a ray from the player's position at a given angle.
func CastRay(startX, startY, angle float64) *DM.RayHit {
	for angle < 0 {
		angle += 2 * math.Pi
	}
	for angle >= 2*math.Pi {
		angle -= 2 * math.Pi
	}
	rayX := math.Cos(angle)
	rayY := math.Sin(angle)
	mapX := int(startX / 100)
	mapY := int(startY / 100)
	deltaDistX := math.Abs(1 / rayX)
	deltaDistY := math.Abs(1 / rayY)
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
	var isVertical bool
	WM := DM.GlobalMaps.Maps
	for {
		if sideDistX < sideDistY {
			sideDistX += deltaDistX * 100
			mapX += stepX
			isVertical = true
		} else {
			sideDistY += deltaDistY * 100
			mapY += stepY
			isVertical = false
		}
		if mapX < 0 || mapX >= len(WM[DM.CurrentMap][0]) || mapY < 0 || mapY >= len(WM[DM.CurrentMap]) {
			return &DM.RayHit{Distance: DM.MaxDepth, WallType: 0}
		}
		if WM[DM.CurrentMap][mapY][mapX] > 0 {
			var distance float64
			if isVertical {
				distance = (float64(mapX) - startX/100 + (1-float64(stepX))/2) / rayX * 100
			} else {
				distance = (float64(mapY) - startY/100 + (1-float64(stepY))/2) / rayY * 100
			}
			return &DM.RayHit{
				Distance:   distance,
				WallType:   WM[DM.CurrentMap][mapY][mapX],
				HitPointX:  startX + rayX*distance,
				HitPointY:  startY + rayY*distance,
				IsVertical: isVertical,
			}
		}
	}
}

// CalculateDarkness calculates the far darkness of the walls based on the distance.
func CalculateDarkness(distance float64) uint8 {
	return uint8(math.Min(DM.MaxDarkness, math.Max(0, distance/3)))
}

// CalculateVerticalDarkness calculates the darkness of the floor and ceiling/roof based on the distance.
func CalculateVerticalDarkness(distance float64) float64 {
	return distance / (float64(DM.ScreenHeight))
}
