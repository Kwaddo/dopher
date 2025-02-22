package renders

import (
	MC "doom/internal/char/player"
	Casts "doom/internal/graphics/casting"
	DM "doom/internal/model"

	"github.com/veandco/go-sdl2/sdl"
)

func RenderFloor(renderer *sdl.Renderer, player *MC.Player) {
	for y := int(DM.ScreenHeight / 2); y < int(DM.ScreenHeight); y++ {
		distance := float64(float64(y) - DM.ScreenHeight/2)
		darkness := Casts.CalculateFloorDarkness(distance)

		// Base brown color (139, 69, 19)
		r := uint8(float64(139) * darkness)
		g := uint8(float64(69) * darkness)
		b := uint8(float64(19) * darkness)

		renderer.SetDrawColor(r, g, b, 255)
		renderer.DrawLine(0, int32(y), int32(DM.ScreenWidth), int32(y))
	}
}
