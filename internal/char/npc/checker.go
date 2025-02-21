package npc

import "math"

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

func (nm *NPCManager) CheckInteraction(playerX, playerY, playerAngle float64) {
	for _, npc := range nm.NPCs {
		npc.ShowDialog = false // Reset dialog for all NPCs
	}

	for _, npc := range nm.NPCs {
		dx := playerX - npc.X
		dy := playerY - npc.Y
		distSquared := dx*dx + dy*dy

		// Check distance first
		if distSquared < 100*100 && !npc.ShowDialog {
			// Now check if player is facing NPC
			angleToNPC := math.Atan2(-dy, -dx) // negative if you want forward angle
			angleDiff := angleToNPC - playerAngle

			// Normalize angleDiff to [-π, π]
			for angleDiff > math.Pi {
				angleDiff -= 2 * math.Pi
			}
			for angleDiff < -math.Pi {
				angleDiff += 2 * math.Pi
			}

			// Only show if within ~45° to either side
			if math.Abs(angleDiff) < math.Pi/4 {
				npc.ShowDialog = true
				npc.DialogTimer = 180
			}
		}
	}
}
