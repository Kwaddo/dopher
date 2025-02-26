package npc

import "math"

// CheckNPCCollision checks if the player is colliding with an NPC.
func (nm *NPCManager) CheckNPCCollision(x, y float64) bool {
	for _, npc := range nm.NPCs {
		dx := x - npc.X
		dy := y - npc.Y
		distSquared := dx*dx + dy*dy
		if distSquared < npc.Hitbox.Radius*npc.Hitbox.Radius {
			return true
		}
	}
	return false
}

// CheckInteraction checks if the player is interacting with an NPC.
func (nm *NPCManager) CheckInteraction(playerX, playerY, playerAngle float64) {
	for _, npc := range nm.NPCs {
		npc.ShowDialog = false
	}
	for _, npc := range nm.NPCs {
		dx := playerX - npc.X
		dy := playerY - npc.Y
		distSquared := dx*dx + dy*dy
		if distSquared < 100*100 && !npc.ShowDialog {
			angleToNPC := math.Atan2(-dy, -dx)
			angleDiff := angleToNPC - playerAngle
			for angleDiff > math.Pi {
				angleDiff -= 2 * math.Pi
			}
			for angleDiff < -math.Pi {
				angleDiff += 2 * math.Pi
			}
			if math.Abs(angleDiff) < math.Pi/4 {
				npc.ShowDialog = true
				npc.DialogTimer = 180
			}
		}
	}
}
