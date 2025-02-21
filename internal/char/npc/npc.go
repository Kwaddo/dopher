package npc

import (
	DM "doom/internal/model"
	"math"
)

type NPCManager DM.NPCManager

var GlobalNPCManager *NPCManager

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
				DialogText:  "Hello traveler!",
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
					Radius: 24,
				},
				DialogText:  "Beef.",
				ShowDialog:  false,
				DialogTimer: 0,
			},
		},
	}
}

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

func (nm *NPCManager) UpdateDistances(playerX, playerY float64) {
	for _, npc := range nm.NPCs {
		dx := npc.X - playerX
		dy := npc.Y - playerY
		npc.Distance = math.Sqrt(dx*dx + dy*dy)
	}
}

func (nm *NPCManager) SortByDistance() {
	// Sort NPCs by distance (furthest first for correct rendering)
	for i := 0; i < len(nm.NPCs)-1; i++ {
		for j := 0; j < len(nm.NPCs)-i-1; j++ {
			if nm.NPCs[j].Distance < nm.NPCs[j+1].Distance {
				nm.NPCs[j], nm.NPCs[j+1] = nm.NPCs[j+1], nm.NPCs[j]
			}
		}
	}
}

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
