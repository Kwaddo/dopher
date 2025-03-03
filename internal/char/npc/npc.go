package npc

import (
	DM "doom/internal/model"
	"math"
)

// NPCManager manages all NPCs in the game.
type NPCManager DM.NPCManager

var GlobalNPCManager *NPCManager

// NewNPCManager creates a new NPCManager with any according number of NPCs.
func NewNPCManager() *NPCManager {
	return &NPCManager{
		NPCs: []*DM.NPC{
			{
				X:       450,
				Y:       450,
				Texture: 3,
				Width:   48,
				Height:  64,
				Hitbox: struct{ Radius float64 }{
					Radius: 24,
				},
				DialogText:  "Hello traveler! How are you?",
				ShowDialog:  false,
				DialogTimer: 0,
			},
			{
				X:       550,
				Y:       550,
				Texture: 4,
				Width:   64,
				Height:  64,
				Hitbox: struct{ Radius float64 }{
					Radius: 12,
				},
				DialogText:  "Beef.",
				ShowDialog:  false,
				DialogTimer: 0,
			},
			{
				X:       1050,
				Y:       1050,
				Texture: 5,
				Width:   64,
				Height:  64,
				Hitbox: struct{ Radius float64 }{
					Radius: 12,
				},
				DialogText:  "Hello gang.",
				ShowDialog:  false,
				DialogTimer: 0,
			},
		},
	}
}

// AddNPC adds a new NPC to the NPCManager.
func (nm *NPCManager) AddNPC(x, y float64, texture int) {
	npc := &DM.NPC{
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
