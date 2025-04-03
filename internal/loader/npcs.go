package loader

import (
	DM "doom/internal/global"
	"encoding/json"
	"fmt"
	"os"
)

type NPCData struct {
	NPCs []DM.NPC
}

func LoadNPCsFromJSON(path string) ([]*DM.NPC, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read npcs file: %v", err)
	}
	var npcData []DM.NPC
	if err := json.Unmarshal(data, &npcData); err != nil {
		return nil, fmt.Errorf("failed to parse npcs JSON: %v", err)
	}
	npcPointers := make([]*DM.NPC, len(npcData))
	for i := range npcData {
		npcPointers[i] = &npcData[i]
	}
	return npcPointers, nil
}
