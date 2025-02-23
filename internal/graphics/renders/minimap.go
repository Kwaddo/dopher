package renders

import (
	MC "doom/internal/char/player"
	DM "doom/internal/model"
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

// Add this function to render the minimap
func RenderMinimap(renderer *sdl.Renderer, player *MC.Player, showMap bool) {
	if !showMap {
		return
	}

	// Calculate minimap dimensions (10% of screen width)
	mapSize := float64(DM.ScreenWidth) * 0.1
	tileSize := mapSize / float64(len(DM.GlobalMap.WorldMap))

	// Draw map background
	renderer.SetDrawColor(0, 0, 0, 200)
	renderer.FillRect(&sdl.Rect{
		X: 10,
		Y: 10,
		W: int32(mapSize),
		H: int32(mapSize),
	})

	// Draw walls
	for y := 0; y < len(DM.GlobalMap.WorldMap); y++ {
		for x := 0; x < len(DM.GlobalMap.WorldMap[y]); x++ {
			if DM.GlobalMap.WorldMap[y][x] > 0 {
				// Different colors for different wall types
				if DM.GlobalMap.WorldMap[y][x] == 1 {
					renderer.SetDrawColor(128, 128, 128, 255) // Gray for type 1
				} else {
					renderer.SetDrawColor(139, 69, 19, 255) // Brown for type 2
				}

				renderer.FillRect(&sdl.Rect{
					X: int32(10 + float64(x)*tileSize),
					Y: int32(10 + float64(y)*tileSize),
					W: int32(tileSize),
					H: int32(tileSize),
				})
			}
		}
	}

	// Replace the player dot with a direction triangle
	playerMapX := int32(10 + (player.X/100.0)*tileSize)
	playerMapY := int32(10 + (player.Y/100.0)*tileSize)

	// Calculate triangle points for direction indicator
	size := int32(6)      // Size of the triangle
	angle := player.Angle // Player's current angle

	// Calculate three points of the triangle
	x1 := playerMapX + int32(float64(size)*math.Cos(angle))
	y1 := playerMapY + int32(float64(size)*math.Sin(angle))

	x2 := playerMapX + int32(float64(size/2)*math.Cos(angle+2.617)) // angle + 150 degrees
	y2 := playerMapY + int32(float64(size/2)*math.Sin(angle+2.617))

	x3 := playerMapX + int32(float64(size/2)*math.Cos(angle-2.617)) // angle - 150 degrees
	y3 := playerMapY + int32(float64(size/2)*math.Sin(angle-2.617))

	// Draw the triangle
	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.DrawLine(x1, y1, x2, y2)
	renderer.DrawLine(x2, y2, x3, y3)
	renderer.DrawLine(x3, y3, x1, y1)
}
