package core

import (
	NPC "doom/internal/char/npc"
	MC "doom/internal/char/player"
	Graphics "doom/internal/graphics"
	DM "doom/internal/model"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

func GameLoop(renderer *sdl.Renderer, player *MC.Player) {
	// Load wall textures
	textures, err := Graphics.LoadTextures(renderer)
	if err != nil {
		panic(err)
	}
	// Clean up textures
	defer func() {
		for _, texture := range textures.Textures {
			texture.Destroy()
		}
	}()

	var DynamicFOV float64
	var currentFOV = DM.FOV
	var targetFOV = DM.FOV
	const lerpSpeed = 0.15

	// Create a channel for render slices
	renderChan := make(chan []*Graphics.RenderSlice, 1)

	// Initialize NPC manager globally
	npcManager := NPC.NewNPCManager()
	NPC.GlobalNPCManager = npcManager // Set the global reference

	// Create a z-buffer to store wall distances
	zBuffer := make([]float64, int(DM.ScreenWidth))

	// Initialize dialog renderer
	dialogRenderer, err := NPC.NewDialogRenderer()
	if err != nil {
		panic(err)
	}
	defer dialogRenderer.Close()

	window, err := sdl.GetWindowFromID(1)
	if err != nil {
		panic(err)
	}

	for {
		// Handle events
		HandleEvents(window, renderer, &zBuffer) // pass zBuffer as pointer so we can reassign

		// Player controls
		end := player.Movement(sdl.GetKeyboardState(), npcManager)
		if end {
			break
		}

		// Smooth FOV transition
		if player.Running {
			targetFOV = DM.FOV * 0.90
		} else {
			targetFOV = DM.FOV
		}

		currentFOV = MC.LERP(currentFOV, targetFOV, lerpSpeed)
		DynamicFOV = currentFOV

		// Update NPC distances
		npcManager.UpdateDistances(player.X, player.Y)
		npcManager.SortByDistance()

		// Clear z-buffer
		for i := range zBuffer {
			zBuffer[i] = math.MaxFloat64
		}

		// Update NPC dialogs
		npcManager.UpdateDialogs()

		// Render Scene
		Graphics.RenderScene(renderer, textures, player, DynamicFOV, renderChan, zBuffer, npcManager, dialogRenderer)

		sdl.Delay(16) // ~60 FPS
	}
}

func HandleEvents(window *sdl.Window, renderer *sdl.Renderer, zBuffer *[]float64) {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			break

		case *sdl.KeyboardEvent:
			if t.Keysym.Sym == sdl.K_f && t.State == sdl.PRESSED {
				// Toggle fullscreen mode
				flags := window.GetFlags()
				if flags&sdl.WINDOW_FULLSCREEN_DESKTOP == sdl.WINDOW_FULLSCREEN_DESKTOP {
					window.SetFullscreen(0)
					DM.ScreenWidth = 1500
					DM.ScreenHeight = 900
					renderer.SetViewport(&sdl.Rect{X: 0, Y: 0, W: int32(DM.ScreenWidth), H: int32(DM.ScreenHeight)})
					// Allocate a new slice for zBuffer so it can grow/shrink as needed
					*zBuffer = make([]float64, int(DM.ScreenWidth))
				} else {
					window.SetFullscreen(sdl.WINDOW_FULLSCREEN_DESKTOP)
					// Get the screen's surface
					screenSurface, err := window.GetSurface()
					if err != nil {
						panic(err)
					}
					DM.ScreenWidth = float64(screenSurface.W)
					DM.ScreenHeight = float64(screenSurface.H)
					renderer.SetViewport(&sdl.Rect{X: 0, Y: 0, W: int32(DM.ScreenWidth), H: int32(DM.ScreenHeight)})
					// Allocate a new slice for zBuffer so it can grow/shrink as needed
					*zBuffer = make([]float64, int(DM.ScreenWidth))
				}
			}
		}
	}
}
