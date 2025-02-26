package renders

import (
	NPC "doom/internal/char/npc"
	MC "doom/internal/char/player"
	Casts "doom/internal/graphics/casting"
	DM "doom/internal/model"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

// RenderNPCs renders all the NPCs in the scene.
func RenderNPCs(player *MC.Player, npcManager *NPC.NPCManager, DynamicFOV float64, zBuffer []float64) []*RenderSlice {
	sprites := make([]*RenderSlice, 0, len(npcManager.NPCs))
	npcManager.UpdateDistances(player.X, player.Y)
	var eyeOffset float64 = 0
	if player.Crouching {
		heightRatio := player.Height / player.DefaultHeight
		eyeOffset = DM.ScreenHeight * (1 - heightRatio) * 0.1
	}
	for _, npc := range npcManager.NPCs {
		dx := npc.X - player.X
		dy := npc.Y - player.Y
		distance := math.Sqrt(dx*dx + dy*dy)
		spriteAngle := math.Atan2(dy, dx) - player.Angle
		for spriteAngle < -math.Pi {
			spriteAngle += 2 * math.Pi
		}
		for spriteAngle > math.Pi {
			spriteAngle -= 2 * math.Pi
		}
		if math.Abs(spriteAngle) > DynamicFOV/2 {
			continue
		}
		spriteHeight := (DM.ScreenHeight / distance) * npc.Height
		spriteWidth := spriteHeight * (npc.Width / npc.Height)
		spriteScreenX := (DM.ScreenWidth / 2) + math.Tan(spriteAngle)*DM.ScreenWidth/DynamicFOV
		spriteTop := (DM.ScreenHeight-spriteHeight)/2 + player.BobOffset - eyeOffset
		startX := int32(spriteScreenX - spriteWidth/2)
		endX := int32(spriteScreenX + spriteWidth/2)
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
		darkness := Casts.CalculateDarkness(distance)
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
