package mapmodel

import (
	Loader "doom/internal/loader"
	DM "doom/internal/models/global"
	"fmt"
	"path/filepath"
)

// GlobalMaps holds the maps data loaded from JSON.
var GlobalMaps *DM.Map

// Initialize GlobalMaps during package initialization.
func init() {
	mapPath := filepath.Join("assets", "maps", "maps.json")
	mapArray, err := Loader.LoadMapsFromJSON(mapPath)
	if err != nil {
		fmt.Printf("Warning: Could not load maps from JSON (%v). Using fallback map data.\n", err)
		GlobalMaps = &DM.Map{
			Maps: [][][]int{
				{
					{1, 1, 1},
					{1, 0, 1},
					{1, 1, 1},
				},
				{
					{2, 2, 2},
					{2, 0, 2},
					{2, 2, 2},
				},
			},
		}
	} else {
		GlobalMaps = &DM.Map{Maps: mapArray}
	}
}
