package player

import (
	NPC "doom/internal/character/npc"
	DM "doom/internal/global"
	Visual "doom/internal/graphics/renders/visual"
	MapModel "doom/internal/mapmodel"
	"math"
)

// Store original positions and map information
var (
	originalPlayerX float64
	originalPlayerY float64
	originalEnemyX  float64
	originalEnemyY  float64
	originalMap     int = 0
	enemyInBattle   int = -1
)

// HandleBattle handles the battle between the player and an enemy NPC.
func (p *Player) HandleBattle(npcManager *NPC.NPCManager, oldX float64, oldY float64) {
	if enemyInBattle >= 0 && !npcManager.NPCs[enemyInBattle].IsAlive {
		Visual.StartTransition(func() {
			DM.CurrentMap = originalMap
			p.X = originalPlayerX + 32
			p.Y = originalPlayerY
			p.VelocityX = 0
			p.VelocityY = 0
			npcManager.NPCs[enemyInBattle].X = originalEnemyX
			npcManager.NPCs[enemyInBattle].Y = originalEnemyY
			npcManager.NPCs[enemyInBattle].MapIndex = originalMap
			npcManager.NPCs[enemyInBattle].DialogText = "I'll get you next time!"
			npcManager.NPCs[enemyInBattle].ShowDialog = true
			npcManager.NPCs[enemyInBattle].DialogTimer = 180
			enemyInBattle = -1
		})
		return
	}
	if DM.CountdownFreeze {
		p.VelocityX = 0
		p.VelocityY = 0
		return
	}
	if npcManager != nil {
		collides, npcIndex := npcManager.CheckNPCCollision(p.X, p.Y)
		if collides {
			if npcIndex >= 0 && npcManager.NPCs[npcIndex].IsEnemy && npcManager.NPCs[npcIndex].IsAlive {
				originalPlayerX = p.X
				originalPlayerY = p.Y
				originalEnemyX = npcManager.NPCs[npcIndex].X
				originalEnemyY = npcManager.NPCs[npcIndex].Y
				originalMap = DM.CurrentMap
				enemyInBattle = npcIndex
				Visual.StartTransition(func() {
					DM.CurrentMap = 1
					p.X = 300
					p.Y = 300
					p.VelocityX = 0
					p.VelocityY = 0
					npcManager.NPCs[npcIndex].X = 400
					npcManager.NPCs[npcIndex].Y = 400
					npcManager.NPCs[npcIndex].DialogText = "Welcome to my domain!"
					npcManager.NPCs[npcIndex].MapIndex = 1
					npcManager.NPCs[npcIndex].ShowDialog = true
					npcManager.NPCs[npcIndex].DialogTimer = 180
					if DM.CurrentMap < len(MapModel.GlobalMaps.Maps) {
						mapWidth := len(MapModel.GlobalMaps.Maps[DM.CurrentMap][0]) * 100
						mapHeight := len(MapModel.GlobalMaps.Maps[DM.CurrentMap]) * 100
						p.X = math.Min(math.Max(p.X, 150), float64(mapWidth-150))
						p.Y = math.Min(math.Max(p.Y, 150), float64(mapHeight-150))
						npcManager.NPCs[npcIndex].X = math.Min(math.Max(npcManager.NPCs[npcIndex].X, 150), float64(mapWidth-150))
						npcManager.NPCs[npcIndex].Y = math.Min(math.Max(npcManager.NPCs[npcIndex].Y, 150), float64(mapHeight-150))
					}
					Visual.StartCountdown()
				})
				return
			} else {
				p.X = oldX
				p.Y = oldY
				p.VelocityX = 0
				p.VelocityY = 0
			}
		}
	}
}
