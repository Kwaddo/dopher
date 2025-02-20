package graphics

import (
	DM "doom/internal/model"
	"github.com/veandco/go-sdl2/sdl"
)

func LoadTextures(renderer *sdl.Renderer) (*DM.TextureMap, error) {
	textures := &DM.TextureMap{
		Textures: make(map[int]*sdl.Texture),
	}

	// Load different textures
	textureFiles := map[int]string{
		1: "assets/wall.bmp",
		2: "assets/npc.bmp", 
	}

	for textureType, file := range textureFiles {
		surface, err := sdl.LoadBMP(file)
		if err != nil {
			return nil, err
		}
		defer surface.Free()

		texture, err := renderer.CreateTextureFromSurface(surface)
		if err != nil {
			return nil, err
		}
		textures.Textures[textureType] = texture
	}

	return textures, nil
}
