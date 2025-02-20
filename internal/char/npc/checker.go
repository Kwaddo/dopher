package npc

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

func (nm *NPCManager) CheckInteraction(playerX, playerY float64) {
	for _, npc := range nm.NPCs {
		dx := playerX - npc.X
		dy := playerY - npc.Y
		distSquared := dx*dx + dy*dy

		// Check if player is within interaction range (slightly larger than hitbox)
		if distSquared < 100*100 && !npc.ShowDialog {
			npc.ShowDialog = true
			npc.DialogTimer = 180 // Show dialog for 3 seconds (60 fps * 3)
		}
	}
}
