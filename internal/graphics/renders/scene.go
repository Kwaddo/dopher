package renders

import (
	Dialogue "doom/internal/character/dialogue"
	NPC "doom/internal/character/npc"
	MC "doom/internal/character/player"
	DM "doom/internal/model"

	"github.com/veandco/go-sdl2/sdl"
)

// Render scene renders the entire scene, including walls, floors, ceilings, NPCs, and the player.
func RenderScene(
	renderer *sdl.Renderer,
	textures *DM.TextureMap,
	player *MC.Player,
	pDynamicFOV *float64,
	renderChan chan []*DM.RenderSlice,
	pZBuffer *[]float64,
	npcManager *NPC.NPCManager,
	dialogRenderer *Dialogue.DialogueRenderer,
	pShowMap *bool,
	pShowMegaMap *bool,
	npcRenderChan chan []*DM.RenderSlice,
) {
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()
	RenderFloor(renderer, player)
	RenderRoof(renderer, player)
	renderDone := make(chan struct{})
	npcRenderDone := make(chan struct{})
	go func() {
		RenderSlices(player, *pDynamicFOV, renderChan)
		renderDone <- struct{}{}
	}()
	wallSlices := <-renderChan
	<-renderDone
	batches := make(map[int]map[uint8]*DM.RenderBatch)
	for _, slice := range wallSlices {
		if texture, ok := textures.Textures[slice.WallType]; ok {
			if _, exists := batches[slice.WallType]; !exists {
				batches[slice.WallType] = make(map[uint8]*DM.RenderBatch)
			}
			if _, exists := batches[slice.WallType][slice.Darkness]; !exists {
				batches[slice.WallType][slice.Darkness] = &DM.RenderBatch{
					Texture:  texture,
					Darkness: slice.Darkness,
					Slices:   make([]*sdl.Rect, 0, 10),
					SrcRects: make([]*sdl.Rect, 0, 10),
				}
			}
			batch := batches[slice.WallType][slice.Darkness]
			batch.Slices = append(batch.Slices, slice.DstRect)

			srcRect := &sdl.Rect{
				X: slice.TexCoord,
				Y: 0,
				W: 1,
				H: 64,
			}
			batch.SrcRects = append(batch.SrcRects, srcRect)
			screenX := int(slice.DstRect.X)
			for x := screenX; x < screenX+int(slice.DstRect.W) && x < int(DM.ScreenWidth); x++ {
				if x >= 0 && x < len(*pZBuffer) {
					(*pZBuffer)[x] = slice.Distance
				}
			}
		}
	}
	var colorModCache = make(map[uint8][3]uint8)
	for _, textureBatches := range batches {
		for _, batch := range textureBatches {
			colorMod, exists := colorModCache[batch.Darkness]
			if !exists {
				value := 255 - batch.Darkness
				colorMod = [3]uint8{value, value, value}
				colorModCache[batch.Darkness] = colorMod
			}

			batch.Texture.SetColorMod(colorMod[0], colorMod[1], colorMod[2])
			for i, dstRect := range batch.Slices {
				renderer.Copy(batch.Texture, batch.SrcRects[i], dstRect)
			}
		}
	}
	go func() {
		npcSlices := RenderNPCs(player, npcManager, *pDynamicFOV, *pZBuffer)
		npcRenderChan <- npcSlices
		npcRenderDone <- struct{}{}
	}()
	npcSlices := <-npcRenderChan
	<-npcRenderDone
	for _, sprite := range npcSlices {
		if texture, ok := textures.Textures[sprite.WallType]; ok {
			texture.SetColorMod(255-sprite.Darkness, 255-sprite.Darkness, 255-sprite.Darkness)
			texture.SetBlendMode(sdl.BLENDMODE_BLEND)
			dstRect := sprite.DstRect
			for x := dstRect.X; x < dstRect.X+dstRect.W; x++ {
				if x >= 0 && x < int32(len(*pZBuffer)) {
					if sprite.Distance > (*pZBuffer)[x] {
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
	if *pShowMegaMap {
		RenderMegaMap(renderer, player, *pShowMegaMap)
	} else {
		RenderMinimap(renderer, player, *pShowMap)
	}
	RenderGun(renderer, player, textures)
	for _, npc := range npcManager.NPCs {
		if npc.ShowDialog {
			if npc.DialogueTree != nil {
				err := dialogRenderer.RenderDialogueWithOptions(renderer, npc)
				if err != nil {
					continue
				}
			} else {
				err := dialogRenderer.RenderSimpleDialogue(renderer, npc.DialogText)
				if err != nil {
					continue
				}
			}
		}
	}
	renderer.Present()
}
