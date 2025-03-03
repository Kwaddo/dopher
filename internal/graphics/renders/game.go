package renders

import (
	NPC "doom/internal/char/npc"
	MC "doom/internal/char/player"
	DM "doom/internal/model"
	"github.com/veandco/go-sdl2/sdl"
)

// Global dialog renderer
var DialogRenderer *NPC.DialogRenderer

// RenderGame handles the rendering of the game scene.
func RenderGame(renderer *sdl.Renderer, player *MC.Player, npcManager *NPC.NPCManager) {
	RenderScene(
		renderer,
		DM.GlobalTextures,
		player,
		&DM.DynamicFOV,
		DM.RenderChan,
		&DM.ZBuffer,
		npcManager,
		DialogRenderer,
		&DM.ShowMiniMap,
		&DM.ShowMegaMap,
	)
}
