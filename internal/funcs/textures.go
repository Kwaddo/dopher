package funcs

import (
	"github.com/veandco/go-sdl2/sdl"
)

func loadTexture(renderer *sdl.Renderer) (*sdl.Texture, error) {
	surface, err := sdl.LoadBMP("assets/wall.bmp")
	if err != nil {
		return nil, err
	}
	defer surface.Free()

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return nil, err
	}

	return texture, nil
}
