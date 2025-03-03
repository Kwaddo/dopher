package core

import (
	NPC "doom/internal/char/npc"
	MC "doom/internal/char/player"
	Casts "doom/internal/graphics/casting"
	Graphics "doom/internal/graphics/renders"
	DM "doom/internal/model"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

// GameLoop is the main game loop for the game.
func GameLoop(renderer *sdl.Renderer, player *MC.Player) {
	var err error
	DM.GlobalTextures, err = Casts.LoadTextures(renderer)
	if err != nil {
		panic(err)
	}
	defer func() {
		for _, texture := range DM.GlobalTextures.Textures {
			texture.Destroy()
		}
	}()
	DM.RenderChan = make(chan []*DM.RenderSlice, 1)
	DM.ZBuffer = make([]float64, int(DM.ScreenWidth))
	npcManager := NPC.NewNPCManager()
	NPC.GlobalNPCManager = npcManager
	Graphics.DialogRenderer, err = NPC.NewDialogRenderer()
	if err != nil {
		panic(err)
	}
	defer Graphics.DialogRenderer.Close()
	window, err := sdl.GetWindowFromID(1)
	if err != nil {
		panic(err)
	}
	if err := Casts.InitFonts(); err != nil {
		panic(err)
	}
	Graphics.GlobalPauseMenu = &Graphics.PauseMenu{
		CurrentOption: 0,
		Options:       []string{"Resume", "Quit"},
	}
	for {
		if HandleEvents(window, renderer) {
			break
		}
		if !DM.GlobalGameState.IsPaused {
			end := player.Movement(sdl.GetKeyboardState(), npcManager)
			if end {
				break
			}
			if player.Running {
				DM.TargetFOV = DM.FOV * 0.92
			} else {
				DM.TargetFOV = DM.FOV
			}
			DM.CurrentFOV = MC.LERP(DM.CurrentFOV, DM.TargetFOV, DM.LerpSpeed)
			DM.DynamicFOV = DM.CurrentFOV
			npcManager.UpdateDistances(player.X, player.Y)
			npcManager.SortByDistance()
			for i := range DM.ZBuffer {
				DM.ZBuffer[i] = math.MaxFloat64
			}
			npcManager.UpdateDialogs()
			npcManager.UpdateEnemies(player.X, player.Y)
		}
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()
		if !DM.GlobalGameState.IsPaused {
			Graphics.RenderGame(renderer, player, npcManager)
		} else {
			Graphics.RenderPauseMenu(renderer)
		}
		renderer.Present()
		sdl.Delay(16)
	}
}
