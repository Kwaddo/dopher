package renders

import (
	MC "doom/internal/character/player"
	Casts "doom/internal/graphics/casting"
	DM "doom/internal/models/global"

	"github.com/veandco/go-sdl2/sdl"
)

// RenderRoof renders the skybox/roof of the game.
func RenderRoof(renderer *sdl.Renderer, player *MC.Player) {
	for y := 0; y < int(DM.ScreenHeight/2); y++ {
		distance := DM.ScreenHeight/2 - float64(y)
		darkness := Casts.CalculateVerticalDarkness(distance)
		r := uint8(float64(112) * darkness)
		g := uint8(float64(112) * darkness)
		b := uint8(float64(112) * darkness)
		renderer.SetDrawColor(r, g, b, 255)
		renderer.DrawLine(0, int32(y), int32(DM.ScreenWidth), int32(y))
	}
}
