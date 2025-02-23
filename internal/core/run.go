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

	// Convert FOV variables to pointers
	var currentFOV = DM.FOV
	var targetFOV = DM.FOV
	var DynamicFOV = currentFOV

	// Pass a pointer to these variables
	pCurrentFOV := &currentFOV
	pTargetFOV := &targetFOV
	pDynamicFOV := &DynamicFOV

	const lerpSpeed = 0.15

	renderChan := make(chan []*Graphics.RenderSlice, 1)

	npcManager := NPC.NewNPCManager()
	NPC.GlobalNPCManager = npcManager

	// Make zBuffer a pointer to a slice
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

	// Convert these into pointers
	showMiniMap := true
	showMegaMap := false
	pShowMiniMap := &showMiniMap
	pShowMegaMap := &showMegaMap

	for {
		// Updated signature uses pointers
		HandleEvents(window, renderer, pZBuf, pShowMiniMap, pShowMegaMap)

		end := player.Movement(sdl.GetKeyboardState(), npcManager)
		if end {
			break
		}

		// Smooth FOV transition
		if player.Running {
			*pTargetFOV = DM.FOV * 0.92
		} else {
			*pTargetFOV = DM.FOV
		}

		*pCurrentFOV = MC.LERP(*pCurrentFOV, *pTargetFOV, lerpSpeed)
		*pDynamicFOV = *pCurrentFOV

		npcManager.UpdateDistances(player.X, player.Y)
		npcManager.SortByDistance()

		// Clear z-buffer
		for i := range *pZBuf {
			(*pZBuf)[i] = math.MaxFloat64
		}

		npcManager.UpdateDialogs()

		// Pass pointers for the FOV and booleans
		Graphics.RenderScene(renderer, textures, player, pDynamicFOV, renderChan, pZBuf, npcManager, dialogRenderer, pShowMiniMap, pShowMegaMap)

		sdl.Delay(16)
	}
}

// Convert parameters to pointers
func HandleEvents(window *sdl.Window, renderer *sdl.Renderer, zBuffer *[]float64, showMap *bool, showMegaMap *bool) {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			break

		case *sdl.KeyboardEvent:
			if t.State == sdl.PRESSED {
				if t.Keysym.Sym == sdl.K_f {
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
				}

				if t.Keysym.Sym == sdl.K_TAB {
					*showMap = !*showMap
					*showMegaMap = false
				}
				if t.Keysym.Sym == sdl.K_m {
					*showMegaMap = !*showMegaMap
					*showMap = false
				}
			}
		}
	}
}
