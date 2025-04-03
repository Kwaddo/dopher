package player

import (
	DM "doom/internal/global"
	MapModel "doom/internal/mapmodel"
	"math"
)

// CheckCollision checks if the player is colliding with a wall.
func CheckCollision(x, y float64) bool {
	WM := MapModel.GlobalMaps.Maps
	for angle := 0.0; angle < 2*math.Pi; angle += math.Pi / 4 {
		checkX := x + math.Cos(angle)*DM.CollisionBuffer
		checkY := y + math.Sin(angle)*DM.CollisionBuffer

		mapX := int(checkX / 100)
		mapY := int(checkY / 100)

		if mapX >= 0 && mapX < len(WM[DM.CurrentMap][DM.CurrentMap]) && mapY >= 0 && mapY < len(WM[DM.CurrentMap]) {
			if WM[DM.CurrentMap][mapY][mapX] != 0 {
				return true
			}
		}
	}
	return false
}
