package renders

import (
	MC "doom/internal/character/player"
	DM "doom/internal/global"
	Casts "doom/internal/graphics/casting"

	"github.com/veandco/go-sdl2/sdl"
)

// RenderFloor renders the floor of the game.
func RenderFloor(renderer *sdl.Renderer, player *MC.Player) {
	for y := int(DM.ScreenHeight / 2); y < int(DM.ScreenHeight); y++ {
		distance := float64(y) - DM.ScreenHeight/2
		darkness := Casts.CalculateVerticalDarkness(distance)
		r := uint8(float64(139) * darkness)
		g := uint8(float64(69) * darkness)
		b := uint8(float64(19) * darkness)
		renderer.SetDrawColor(r, g, b, 255)
		renderer.DrawLine(0, int32(y), int32(DM.ScreenWidth), int32(y))
	}
}
