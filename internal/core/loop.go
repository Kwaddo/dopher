package core

import (
	NPC "doom/internal/character/npc"
	MC "doom/internal/character/player"
	Casts "doom/internal/graphics/casting"
	Graphics "doom/internal/graphics/renders"
	DM "doom/internal/model"
	Menu "doom/internal/ui"
	"math"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
)

// RunGameLoop is the main game loop for the game.
func RunGameLoop(renderer *sdl.Renderer, player *MC.Player) {
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
	npcRenderChan := make(chan []*DM.RenderSlice, 1)
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
	Menu.GlobalPauseMenu = Menu.NewPauseMenu()
	Menu.GlobalMainMenu = Menu.NewMainMenu()
	Menu.GlobalOptionsMenu = Menu.NewOptionsMenu()
	DM.GlobalGameState = DM.GameState{
		InMainMenu:    true,
		InOptionsMenu: false,
		IsPaused:      false,
	}
	frameBuffer1, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888,
		sdl.TEXTUREACCESS_TARGET,
		int32(DM.ScreenWidth),
		int32(DM.ScreenHeight))
	if err != nil {
		panic(err)
	}
	defer frameBuffer1.Destroy()
	frameBuffer2, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888,
		sdl.TEXTUREACCESS_TARGET,
		int32(DM.ScreenWidth),
		int32(DM.ScreenHeight))
	if err != nil {
		panic(err)
	}
	defer frameBuffer2.Destroy()
	currentBuffer := frameBuffer1
	targetFPS := 60
	frameDelay := uint32(1000 / targetFPS)
	var bufferMutex sync.Mutex

	for {
		frameStart := sdl.GetTicks()
		if HandleEvents(window, renderer) {
			break
		}
		DM.GlobalFrameCount++
		renderer.SetRenderTarget(currentBuffer)
		renderer.Clear()

		if DM.GlobalGameState.InMainMenu {
			Menu.RenderMainMenu(renderer)
		} else if DM.GlobalGameState.InOptionsMenu {
			Menu.RenderOptionsMenu(renderer)
		} else if DM.GlobalGameState.IsPaused {
			Menu.RenderPauseMenu(renderer)
		} else {
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
			if len(DM.ZBuffer) != int(DM.ScreenWidth) {
				DM.ZBuffer = make([]float64, int(DM.ScreenWidth))
			} else {
				for i := range DM.ZBuffer {
					DM.ZBuffer[i] = math.MaxFloat64
				}
			}
			npcManager.UpdateDialogs()
			npcManager.UpdateEnemies(player.X, player.Y)
			Graphics.RenderGame(renderer, player, npcManager, npcRenderChan)
		}
		renderer.SetRenderTarget(nil)
		renderer.Clear()
		renderer.Copy(currentBuffer, nil, nil)
		renderer.Present()
		if currentBuffer == frameBuffer1 {
			currentBuffer = frameBuffer2
		} else {
			currentBuffer = frameBuffer1
		}
		sdl.Delay(16)
		frameTime := sdl.GetTicks() - frameStart
		if frameDelay > frameTime {
			sdl.Delay(frameDelay - frameTime)
		}
		if DM.NeedToRecreateBuffers {
			bufferMutex.Lock()
			if frameBuffer1 != nil {
				frameBuffer1.Destroy()
			}
			if frameBuffer2 != nil {
				frameBuffer2.Destroy()
			}
			frameBuffer1, err = renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888,
				sdl.TEXTUREACCESS_TARGET,
				int32(DM.ScreenWidth),
				int32(DM.ScreenHeight))
			if err != nil {
				panic(err)
			}

			frameBuffer2, err = renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888,
				sdl.TEXTUREACCESS_TARGET,
				int32(DM.ScreenWidth),
				int32(DM.ScreenHeight))
			if err != nil {
				panic(err)
			}

			currentBuffer = frameBuffer1
			DM.NeedToRecreateBuffers = false
			bufferMutex.Unlock()
		}
	}
}
