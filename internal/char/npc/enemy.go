package npc

import (
	DM "doom/internal/model"
	"math"
)

// UpdateEnemies updates all enemy NPCs in the game.
func (nm *NPCManager) UpdateEnemies(playerX, playerY float64) {
	for _, npc := range nm.NPCs {
		if !npc.IsEnemy {
			continue
		}
		dx := playerX - npc.X
		dy := playerY - npc.Y
		distSquared := dx*dx + dy*dy
		if distSquared <= npc.DetectionRadius*npc.DetectionRadius {
			npc.State = DM.EnemyStateChasing
			if npc.LastMoveTime <= 0 {
				dist := math.Sqrt(distSquared)
				dirX := dx / dist
				dirY := dy / dist
				npc.LastDirection.X = dirX
				npc.LastDirection.Y = dirY
				newX := npc.X + dirX*npc.Speed
				newY := npc.Y + dirY*npc.Speed
				if !CheckWallCollision(newX, newY, npc.Hitbox.Radius) {
					npc.X = newX
					npc.Y = newY
					npc.PathBlockedTime = 0
				} else {
					npc.PathBlockedTime++
					leftDirX := -dirY
					leftDirY := dirX
					rightDirX := dirY
					rightDirY := -dirX
					newLeftX := npc.X + leftDirX*npc.Speed
					newLeftY := npc.Y + leftDirY*npc.Speed
					newRightX := npc.X + rightDirX*npc.Speed
					newRightY := npc.Y + rightDirY*npc.Speed
					if !CheckWallCollision(newLeftX, newLeftY, npc.Hitbox.Radius) {
						npc.X = newLeftX
						npc.Y = newLeftY
					} else if !CheckWallCollision(newRightX, newRightY, npc.Hitbox.Radius) {
						npc.X = newRightX
						npc.Y = newRightY
					}
				}

				// Reset move timer - determines movement speed
				npc.LastMoveTime = 2
			} else {
				npc.LastMoveTime--
			}

			// Check if close enough to attack
			if distSquared <= (npc.Hitbox.Radius+30)*(npc.Hitbox.Radius+30) {
				npc.State = DM.EnemyStateAttack
				if npc.DialogTimer <= 0 {
					npc.DialogText = "ATTACK!"
					npc.ShowDialog = true
					npc.DialogTimer = 60
				}
			}
		} else {
			// Player not detected - return to idle
			npc.State = DM.EnemyStateIdle
		}
	}
}
