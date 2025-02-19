package core

import (
	"math"
	DM "doom/internal/constants"
)

func CheckCollision(x, y float64) bool {
	// Check multiple points around the player
	WM := DM.GlobalMap.WorldMap
	for angle := 0.0; angle < 2*math.Pi; angle += math.Pi / 4 {
		checkX := x + math.Cos(angle)*DM.CollisionBuffer
		checkY := y + math.Sin(angle)*DM.CollisionBuffer

		mapX := int(checkX / 100)
		mapY := int(checkY / 100)

		if mapX >= 0 && mapX < len(WM[0]) && mapY >= 0 && mapY < len(WM) {
			if WM[mapY][mapX] == 1 {
				return true
			}
		}
	}
	return false
}

func (nm *NPCManager) CheckNPCCollision(x, y float64) bool {
	for _, npc := range nm.NPCs {
		dx := x - npc.X
		dy := y - npc.Y
		distSquared := dx*dx + dy*dy

		// If distance is less than hitbox radius, collision occurred
		if distSquared < npc.Hitbox.Radius*npc.Hitbox.Radius {
			return true
		}
	}
	return false
}
