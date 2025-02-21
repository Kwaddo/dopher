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
				W: int32(math.Ceil(rayWidth + 1)), // Ensure no gaps between slices
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
		spriteTop := (DM.ScreenHeight - spriteHeight) / 2

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
