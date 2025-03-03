package casting

import (
	DM "doom/internal/model"

	"github.com/veandco/go-sdl2/sdl"
)

// LoadTextures loads the textures for the game from the assets folder.
func LoadTextures(renderer *sdl.Renderer) (*DM.TextureMap, error) {
	textures := &DM.TextureMap{
		Textures: make(map[int]*sdl.Texture),
	}
	textureFiles := map[int]string{
		1: "assets/textures/wall.bmp",
		2: "assets/textures/wall2.bmp",
		3: "assets/textures/npc.bmp",
		4: "assets/textures/beef.bmp",
		5: "assets/textures/dictator.bmp",
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
