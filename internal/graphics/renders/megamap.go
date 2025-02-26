package renders

import (
	MC "doom/internal/char/player"
	DM "doom/internal/model"
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

func RenderMegaMap(renderer *sdl.Renderer, player *MC.Player, showMegaMap bool) {
	if !showMegaMap {
		return
	}
	mapSize := float64(DM.ScreenHeight) * 0.9
	tileSize := mapSize / float64(len(DM.GlobalMap.WorldMap))
	mapX := (DM.ScreenWidth - mapSize) / 2
	mapY := (DM.ScreenHeight - mapSize) / 2
	renderer.SetDrawColor(0, 0, 0, 200)
	renderer.FillRect(&sdl.Rect{
		X: int32(mapX),
		Y: int32(mapY),
		W: int32(mapSize),
		H: int32(mapSize),
	})
	for y := 0; y < len(DM.GlobalMap.WorldMap); y++ {
		for x := 0; x < len(DM.GlobalMap.WorldMap[y]); x++ {
			if DM.GlobalMap.WorldMap[y][x] > 0 {
				if DM.GlobalMap.WorldMap[y][x] == 1 {
					renderer.SetDrawColor(128, 128, 128, 255)
				} else {
					renderer.SetDrawColor(139, 69, 19, 255)
				}
				renderer.FillRect(&sdl.Rect{
					X: int32(mapX + float64(x)*tileSize),
					Y: int32(mapY + float64(y)*tileSize),
					W: int32(tileSize),
					H: int32(tileSize),
				})
			}
		}
	}
	playerMapX := int32(mapX + (player.X/100.0)*tileSize)
	playerMapY := int32(mapY + (player.Y/100.0)*tileSize)
	x1 := playerMapX + int32((12)*math.Cos(player.Angle))
	y1 := playerMapY + int32((12)*math.Sin(player.Angle))
	x2 := playerMapX + int32((12/2)*math.Cos(player.Angle+2.617))
	y2 := playerMapY + int32((12/2)*math.Sin(player.Angle+2.617))
	x3 := playerMapX + int32((12/2)*math.Cos(player.Angle-2.617))
	y3 := playerMapY + int32((12/2)*math.Sin(player.Angle-2.617))
	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.DrawLine(x1, y1, x2, y2)
	renderer.DrawLine(x2, y2, x3, y3)
	renderer.DrawLine(x3, y3, x1, y1)
	renderer.SetDrawColor(255, 255, 0, 255)
	points := []sdl.Point{{X: x1, Y: y1}, {X: x2, Y: y2}, {X: x3, Y: y3}}
	renderer.DrawLines(points)
}
