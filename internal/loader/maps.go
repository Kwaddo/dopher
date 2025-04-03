package loader

import (
	DM "doom/internal/global"
	"encoding/json"
	"fmt"
	"os"
)

// LoadMapsFromJSON loads map data from a JSON file.
func LoadMapsFromJSON(path string) ([][][]int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read maps file: %v", err)
	}
	var mapData DM.Map
	if err := json.Unmarshal(data, &mapData); err != nil {
		return nil, fmt.Errorf("failed to parse maps JSON: %v", err)
	}
	return mapData.Maps, nil
}
