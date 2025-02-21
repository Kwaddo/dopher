package graphics

import (
	NPC "doom/internal/char/npc"
	MC "doom/internal/char/player"
	DM "doom/internal/model"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type RenderSlice DM.RenderSlice

func RenderSlices(player *MC.Player, DynamicFOV float64, renderChan chan<- []*RenderSlice) {
	slices := make([]*RenderSlice, DM.NumRays)

	// Adjust ray angle calculation to ensure full coverage
	rayAngleStep := DynamicFOV / float64(DM.NumRays-1) // Subtract 1 to include last ray
	rayAngle := player.Angle - DynamicFOV/2

	// Calculate ray width to ensure full screen coverage
	rayWidth := DM.ScreenWidth / float64(DM.NumRays)

	for i := 0; i < DM.NumRays; i++ {
		rayResult := CastRay(player.X, player.Y, rayAngle)

		// Fix fisheye effect
		distance := rayResult.Distance * math.Cos(rayAngle-player.Angle)

		// Adjust wall height calculation
		wallHeight := (DM.ScreenHeight / distance) * 75
		if wallHeight > DM.ScreenHeight {
			wallHeight = DM.ScreenHeight
		}

		darkness := CalculateDarkness(distance)
		// Apply head bobbing to wall position
		wallTop := (DM.ScreenHeight-wallHeight)/2 + player.BobOffset

		// Improved texture coordinate calculation
		var texCoord int32
		if rayResult.IsVertical {
			texCoord = int32(math.Mod(rayResult.HitPointY, 100) * 0.64)
		} else {
			texCoord = int32(math.Mod(rayResult.HitPointX, 100) * 0.64)
		}

		slices[i] = &RenderSlice{
			DstRect: &sdl.Rect{
				X: int32(float64(i) * rayWidth),
				Y: int32(wallTop),
				W: int32(math.Ceil(rayWidth)),
				H: int32(wallHeight),
			},
			Darkness: darkness,
			Color:    sdl.Color{R: 128, G: 128, B: 128, A: 255},
			WallType: rayResult.WallType,
			TexCoord: texCoord,
			Distance: distance,
		}

		rayAngle += rayAngleStep
	}

	renderChan <- slices
}

func RenderNPCs(player *MC.Player, npcManager *NPC.NPCManager, DynamicFOV float64, zBuffer []float64) []*RenderSlice {
	sprites := make([]*RenderSlice, 0, len(npcManager.NPCs))
	npcManager.UpdateDistances(player.X, player.Y)

	for _, npc := range npcManager.NPCs {
		// Calculate angle and distance to NPC
		dx := npc.X - player.X
		dy := npc.Y - player.Y
		distance := math.Sqrt(dx*dx + dy*dy)

		// Calculate sprite angle relative to player's view
		spriteAngle := math.Atan2(dy, dx) - player.Angle

		// Normalize angle
		for spriteAngle < -math.Pi {
			spriteAngle += 2 * math.Pi
		}
		for spriteAngle > math.Pi {
			spriteAngle -= 2 * math.Pi
		}

		// Check if sprite is in view
		if math.Abs(spriteAngle) > DynamicFOV/2 {
			continue
		}

		// Calculate sprite size and position on screen
		spriteHeight := (DM.ScreenHeight / distance) * npc.Height
		spriteWidth := spriteHeight * (npc.Width / npc.Height)

		spriteScreenX := (DM.ScreenWidth / 2) + math.Tan(spriteAngle)*DM.ScreenWidth/DynamicFOV
		spriteTop := (DM.ScreenHeight-spriteHeight)/2 + player.BobOffset // Apply head bobbing offset

		// Check z-buffer for visibility
		startX := int32(spriteScreenX - spriteWidth/2)
		endX := int32(spriteScreenX + spriteWidth/2)

		// Skip if completely behind walls
		visible := false
		for x := startX; x < endX; x++ {
			if x >= 0 && x < int32(len(zBuffer)) && distance < zBuffer[x] {
				visible = true
				break
			}
		}

		if !visible {
			continue
		}

		darkness := CalculateDarkness(distance)

		sprites = append(sprites, &RenderSlice{
			DstRect: &sdl.Rect{
				X: startX,
				Y: int32(spriteTop),
				W: int32(spriteWidth),
				H: int32(spriteHeight),
			},
			Darkness: darkness,
			WallType: npc.Texture,
			TexCoord: 0,
			Distance: distance,
		})
	}

	return sprites
}

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
