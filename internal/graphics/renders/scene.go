package renders

import (
	NPC "doom/internal/char/npc"
	MC "doom/internal/char/player"
	DM "doom/internal/model"

	"github.com/veandco/go-sdl2/sdl"
)

func RenderScene(
	renderer *sdl.Renderer,
	textures *DM.TextureMap,
	player *MC.Player,
	DynamicFOV float64,
	renderChan chan []*RenderSlice,
	zBuffer []float64,
	npcManager *NPC.NPCManager,
	dialogRenderer *NPC.DialogRenderer,
) {
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()

	// Set blend mode for proper transparency
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)

	// Start goroutine to calculate floor slices
	RenderFloor(renderer, player)

	// Start goroutine to calculate wall slices
	go RenderSlices(player, DynamicFOV, renderChan)

	// Receive and render the slices
	wallSlices := <-renderChan
	for _, slice := range wallSlices {
		if texture, ok := textures.Textures[slice.WallType]; ok {
			texture.SetColorMod(255-slice.Darkness, 255-slice.Darkness, 255-slice.Darkness)
			srcRect := &sdl.Rect{
				X: slice.TexCoord,
				Y: 0,
				W: 1,
				H: 64,
			}
			renderer.Copy(texture, srcRect, slice.DstRect)

			// Store wall distance in z-buffer
			screenX := int(slice.DstRect.X)
			for x := screenX; x < screenX+int(slice.DstRect.W) && x < int(DM.ScreenWidth); x++ {
				if x >= 0 && x < len(zBuffer) {
					zBuffer[x] = slice.Distance
				}
			}
		}
	}

	// Now start a goroutine to compute the NPC slices (no direct rendering!)
	npcRenderChan := make(chan []*RenderSlice, 1)
	go func() {
		npcSlices := RenderNPCs(player, npcManager, DynamicFOV, zBuffer)
		npcRenderChan <- npcSlices
	}()

	// Receive computed NPC slices and render them on the main goroutine
	npcSlices := <-npcRenderChan
	for _, sprite := range npcSlices {
		if texture, ok := textures.Textures[sprite.WallType]; ok {
			texture.SetColorMod(255-sprite.Darkness, 255-sprite.Darkness, 255-sprite.Darkness)
			texture.SetBlendMode(sdl.BLENDMODE_BLEND)

			dstRect := sprite.DstRect
			for x := dstRect.X; x < dstRect.X+dstRect.W; x++ {
				if x >= 0 && x < int32(len(zBuffer)) {
					if sprite.Distance > zBuffer[x] {
						continue
					}

					columnRect := &sdl.Rect{
						X: x,
						Y: dstRect.Y,
						W: 1,
						H: dstRect.H,
					}

					srcX := int32(float64(x-dstRect.X) / float64(dstRect.W) * 64)
					srcColumnRect := &sdl.Rect{
						X: srcX,
						Y: 0,
						W: 1,
						H: 64,
					}

					renderer.Copy(texture, srcColumnRect, columnRect)
				}
			}
		}
	}

	for _, npc := range npcManager.NPCs {
		if npc.ShowDialog {
			err := dialogRenderer.RenderDialog(renderer, npc.DialogText)
			if err != nil {
				continue
			}
		}
	}

	renderer.Present()
}
