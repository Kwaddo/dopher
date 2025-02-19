package funcs

import (
	"github.com/veandco/go-sdl2/sdl"
)

type TextureMap struct {
	Textures map[int]*sdl.Texture
}

func LoadTextures(renderer *sdl.Renderer) (*TextureMap, error) {
	textures := &TextureMap{
		Textures: make(map[int]*sdl.Texture),
	}

	// Load different wall textures
	textureFiles := map[int]string{
		1: "assets/wall.bmp",
		2: "assets/wall.bmp",
		3: "assets/wall.bmp",
		// Add more textures as needed
	}

	for wallType, file := range textureFiles {
		surface, err := sdl.LoadBMP(file)
		if err != nil {
			return nil, err
		}
		defer surface.Free()

		texture, err := renderer.CreateTextureFromSurface(surface)
		if err != nil {
			return nil, err
		}
		textures.Textures[wallType] = texture
	}

	return textures, nil
}
