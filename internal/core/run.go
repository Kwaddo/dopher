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
	textures, err := Casts.LoadTextures(renderer)
	if err != nil {
		panic(err)
	}
	defer func() {
		for _, texture := range textures.Textures {
			texture.Destroy()
		}
	}()
	var currentFOV = DM.FOV
	var targetFOV = DM.FOV
	var DynamicFOV = currentFOV
	pCurrentFOV := &currentFOV
	pTargetFOV := &targetFOV
	pDynamicFOV := &DynamicFOV
	renderChan := make(chan []*Graphics.RenderSlice, 1)
	npcManager := NPC.NewNPCManager()
	NPC.GlobalNPCManager = npcManager
	zBuf := make([]float64, int(DM.ScreenWidth))
	pZBuf := &zBuf
	dialogRenderer, err := NPC.NewDialogRenderer()
	if err != nil {
		panic(err)
	}
	defer dialogRenderer.Close()
	window, err := sdl.GetWindowFromID(1)
	if err != nil {
		panic(err)
	}
	showMiniMap := true
	showMegaMap := false
	pShowMiniMap := &showMiniMap
	pShowMegaMap := &showMegaMap
	if err := DM.InitFonts(); err != nil {
		panic(err)
	}
	DM.GlobalPauseMenu = &DM.PauseMenu{
		CurrentOption: 0,
		Options:       []string{"Resume", "Quit"},
	}
	for {
		if HandleEvents(window, renderer, pZBuf, pShowMiniMap, pShowMegaMap) {
			break
		}
		if !DM.GlobalGameState.IsPaused {
			end := player.Movement(sdl.GetKeyboardState(), npcManager)
			if end {
				break
			}
			if player.Running {
				*pTargetFOV = DM.FOV * 0.92
			} else {
				*pTargetFOV = DM.FOV
			}
			*pCurrentFOV = MC.LERP(*pCurrentFOV, *pTargetFOV, DM.LerpSpeed)
			*pDynamicFOV = *pCurrentFOV
			npcManager.UpdateDistances(player.X, player.Y)
			npcManager.SortByDistance()
			for i := range *pZBuf {
				(*pZBuf)[i] = math.MaxFloat64
			}
			npcManager.UpdateDialogs()
			npcManager.UpdateEnemies(player.X, player.Y)
		}
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()
		if !DM.GlobalGameState.IsPaused {
			Graphics.RenderScene(renderer, textures, player, pDynamicFOV, renderChan, pZBuf, npcManager, dialogRenderer, pShowMiniMap, pShowMegaMap)
		} else {
			Graphics.RenderPauseMenu(renderer)
		}

		renderer.Present()
		sdl.Delay(16)
	}
}

// HandleEvents handles certain events for the game.
func HandleEvents(window *sdl.Window, renderer *sdl.Renderer, zBuffer *[]float64, showMap *bool, showMegaMap *bool) bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			return true
		case *sdl.KeyboardEvent:
			if t.State == sdl.PRESSED {
				if DM.GlobalGameState.IsPaused {
					switch t.Keysym.Sym {
					case sdl.K_UP:
						DM.GlobalPauseMenu.MoveUp()
					case sdl.K_DOWN:
						DM.GlobalPauseMenu.MoveDown()
					case sdl.K_RETURN, sdl.K_SPACE:
						switch DM.GlobalPauseMenu.GetSelectedOption() {
						case "Resume":
							DM.GlobalGameState.IsPaused = false
							return false
						case "Quit":
							return true
						}
					case sdl.K_ESCAPE:
						DM.GlobalGameState.IsPaused = false
					}
					continue
				}
				switch t.Keysym.Sym {
				case sdl.K_ESCAPE:
					DM.GlobalGameState.IsPaused = !DM.GlobalGameState.IsPaused
				case sdl.K_f:
					flags := window.GetFlags()
					if flags&sdl.WINDOW_FULLSCREEN_DESKTOP == sdl.WINDOW_FULLSCREEN_DESKTOP {
						window.SetFullscreen(0)
						DM.ScreenWidth = 1500
						DM.ScreenHeight = 900
						renderer.SetViewport(&sdl.Rect{X: 0, Y: 0, W: int32(DM.ScreenWidth), H: int32(DM.ScreenHeight)})
						*zBuffer = make([]float64, int(DM.ScreenWidth))
					} else {
						window.SetFullscreen(sdl.WINDOW_FULLSCREEN_DESKTOP)
						screenSurface, err := window.GetSurface()
						if err != nil {
							panic(err)
						}
						DM.ScreenWidth = float64(screenSurface.W)
						DM.ScreenHeight = float64(screenSurface.H)
						renderer.SetViewport(&sdl.Rect{X: 0, Y: 0, W: int32(DM.ScreenWidth), H: int32(DM.ScreenHeight)})
						*zBuffer = make([]float64, int(DM.ScreenWidth))
					}
				case sdl.K_TAB:
					*showMap = !*showMap
					*showMegaMap = false
				case sdl.K_m:
					*showMegaMap = !*showMegaMap
					*showMap = false
				}
			}
		}
	}
	return false
}
