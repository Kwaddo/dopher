package npc

import (
	DM "doom/internal/models/global"
	"math"
)

// UpdateEnemies updates all enemy NPCs in the game.
func (nm *NPCManager) UpdateEnemies(playerX, playerY float64) {
	if DM.CountdownFreeze {
		return
	}
	for i, npc := range nm.NPCs {
		if !npc.IsEnemy || !npc.IsAlive {
			continue
		}
		if DM.GlobalFrameCount%3 != i%3 {
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
				npc.LastMoveTime = 2
			} else {
				npc.LastMoveTime--
			}
			if distSquared <= (npc.Hitbox.Radius+30)*(npc.Hitbox.Radius+30) {
				npc.State = DM.EnemyStateAttack
				if npc.DialogTimer <= 0 {
					npc.DialogText = "ATTACK!"
					npc.ShowDialog = true
					npc.DialogTimer = 60
				}
			}
		} else {
			npc.State = DM.EnemyStateIdle
		}
	}
}

// DamageEnemy applies damage to an NPC and handles death logic
func (nm *NPCManager) DamageEnemy(npcIndex int, damage int) bool {
	if npcIndex < 0 || npcIndex >= len(nm.NPCs) {
		return false
	}
	npc := nm.NPCs[npcIndex]
	if !npc.IsEnemy || !npc.IsAlive {
		return false
	}
	npc.Health -= damage
	npc.ShowDialog = true
	npc.DialogText = "Ouch!"
	npc.DialogTimer = 60
	if npc.Health <= 0 {
		npc.Health = 0
		npc.IsAlive = false
		npc.DialogText = "Argh!"
		npc.DialogTimer = 120
		npc.State = DM.EnemyStateIdle
		return true
	}
	return false
}
