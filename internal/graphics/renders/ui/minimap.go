package renders

import (
	MC "doom/internal/character/player"
	DM "doom/internal/models/global"
	MapModel "doom/internal/models/mapmodel"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

// RenderMinimap renders the minimap on the screen.
func RenderMinimap(renderer *sdl.Renderer, player *MC.Player, showMap bool) {
	if !showMap {
		return
	}
	mapSize := float64(DM.ScreenWidth) * 0.1
	tileSize := mapSize / float64(len(MapModel.GlobalMaps.Maps[DM.CurrentMap]))
	renderer.SetDrawColor(0, 0, 0, 200)
	renderer.FillRect(&sdl.Rect{
		X: 10,
		Y: 10,
		W: int32(mapSize),
		H: int32(mapSize),
	})
	for y := 0; y < len(MapModel.GlobalMaps.Maps[DM.CurrentMap]); y++ {
		for x := 0; x < len(MapModel.GlobalMaps.Maps[DM.CurrentMap][y]); x++ {
			if MapModel.GlobalMaps.Maps[DM.CurrentMap][y][x] > 0 {
				if MapModel.GlobalMaps.Maps[DM.CurrentMap][y][x] == 1 {
					renderer.SetDrawColor(128, 128, 128, 255)
				} else {
					renderer.SetDrawColor(139, 69, 19, 255)
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
	playerMapX := int32(10 + (player.X/100.0)*tileSize)
	playerMapY := int32(10 + (player.Y/100.0)*tileSize)
	x1 := playerMapX + int32(float64(6)*math.Cos(player.Angle))
	y1 := playerMapY + int32(float64(6)*math.Sin(player.Angle))
	x2 := playerMapX + int32(float64(6/2)*math.Cos(player.Angle+2.617))
	y2 := playerMapY + int32(float64(6/2)*math.Sin(player.Angle+2.617))
	x3 := playerMapX + int32(float64(6/2)*math.Cos(player.Angle-2.617))
	y3 := playerMapY + int32(float64(6/2)*math.Sin(player.Angle-2.617))
	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.DrawLine(x1, y1, x2, y2)
	renderer.DrawLine(x2, y2, x3, y3)
	renderer.DrawLine(x3, y3, x1, y1)
}
