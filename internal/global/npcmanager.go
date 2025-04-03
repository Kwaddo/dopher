package global

import (
	"math"
)

// AddNPC adds a new NPC to the NPCManager.
func (nm *NPCManager) AddNPC(x, y float64, texture int) {
	npc := &NPC{
		X:       x,
		Y:       y,
		Texture: texture,
		Width:   50,
		Height:  100,
		Hitbox: struct{ Radius float64 }{
			Radius: 25,
		},
	}
	nm.NPCs = append(nm.NPCs, npc)
}

// UpdateDistances updates the distances of all NPCs from the player.
func (nm *NPCManager) UpdateDistances(playerX, playerY float64) {
	for _, npc := range nm.NPCs {
		dx := npc.X - playerX
		dy := npc.Y - playerY
		npc.Distance = math.Sqrt(dx*dx + dy*dy)
	}
}

// SortByDistance sorts the NPCs by distance from the player, furthest first for correct rendering.
func (nm *NPCManager) SortByDistance() {
	for i := 0; i < len(nm.NPCs)-1; i++ {
		for j := 0; j < len(nm.NPCs)-i-1; j++ {
			if nm.NPCs[j].Distance < nm.NPCs[j+1].Distance {
				nm.NPCs[j], nm.NPCs[j+1] = nm.NPCs[j+1], nm.NPCs[j]
			}
		}
	}
}

// UpdateDialogs updates the dialog timers for all NPCs.
func (nm *NPCManager) UpdateDialogs() {
	for _, npc := range nm.NPCs {
		if npc.ShowDialog {
			npc.DialogTimer--
			if npc.DialogTimer <= 0 {
				npc.ShowDialog = false
			}
		}
	}
}
