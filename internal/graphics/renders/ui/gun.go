package renders

import (
	MC "doom/internal/character/player"
	DM "doom/internal/models/global"

	"github.com/veandco/go-sdl2/sdl"
)

// RenderGun renders the gun on the screen.
func RenderGun(renderer *sdl.Renderer, player *MC.Player, textures *DM.TextureMap) {
	if DM.CurrentMap != 1 {
		return
	}
	gunTexture := textures.Textures[6]
	gunWidth := DM.ScreenWidth * 0.4
	gunHeight := gunWidth * 0.5
	offsetY := 0.0
	if player.Walking {
		offsetY = player.BobOffset * 0.5
	}
	recoilOffset := 0.0
	if player.Gun.IsFiring {
		recoilOffset = 10.0
	}
	dstRect := &sdl.Rect{
		X: int32((DM.ScreenWidth - gunWidth) / 2),
		Y: int32(DM.ScreenHeight - gunHeight + offsetY + recoilOffset),
		W: int32(gunWidth),
		H: int32(gunHeight),
	}
	renderer.Copy(gunTexture, nil, dstRect)
	if player.Gun.MuzzleFlash {
		flashTexture := textures.Textures[7]
		flashWidth := gunWidth * 0.3
		flashHeight := flashWidth
		flashRect := &sdl.Rect{
			X: int32(DM.ScreenWidth / 2),
			Y: int32(DM.ScreenHeight - gunHeight*0.8 + offsetY),
			W: int32(flashWidth),
			H: int32(flashHeight),
		}
		renderer.Copy(flashTexture, nil, flashRect)
	}
}
