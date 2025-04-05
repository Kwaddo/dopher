package loader

import (
	DM "doom/internal/models/global"
	"encoding/json"
	"fmt"
	"os"
)

func LoadNPCsFromJSON(path string) ([]*DM.NPC, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read npcs file: %v", err)
	}
	var npcData []DM.NPC
	if err := json.Unmarshal(data, &npcData); err != nil {
		return nil, fmt.Errorf("failed to parse npcs JSON: %v", err)
	}
	npcs := make([]*DM.NPC, len(npcData))
	for i := range npcData {
		npcs[i] = &npcData[i]
		if npcs[i].DialogueTree == nil {
			npcs[i].DialogueTree = CreateBasicDialogueTree()
		}
	}
	return npcs, nil
}
