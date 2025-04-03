package npcmodel

import (
	NPCs "doom/internal/character/npc"
	Loader "doom/internal/loader"
	"fmt"
	"path/filepath"
)

func init() {
	npcPath := filepath.Join("assets", "npcs", "npcs.json")
	npcArray, err := Loader.LoadNPCsFromJSON(npcPath)
	if err != nil {
		fmt.Printf("Warning: Could not load npcs from JSON (%v). Using fallback npcs data.\n", err)
	} else {
		NPCs.GlobalNPCManager = &NPCs.NPCManager{NPCs: npcArray}
	}
}
