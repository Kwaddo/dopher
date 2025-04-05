package core

import (
	Casts "doom/internal/graphics/casting"
	DM "doom/internal/models/global"

	"github.com/veandco/go-sdl2/sdl"
)

// GameContext holds all initialized game resources.
type GameContext DM.GameContext

// The function to initalize everything.
func InitializeAll() *GameContext {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	if err := Casts.InitFonts(); err != nil {
		sdl.Quit()
		panic(err)
	}
	window, err := sdl.CreateWindow("Dopher Engine", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		int32(DM.ScreenWidth), int32(DM.ScreenHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		Casts.CleanupFonts()
		sdl.Quit()
		panic(err)
	}
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		window.Destroy()
		Casts.CleanupFonts()
		sdl.Quit()
		panic(err)
	}
	return &GameContext{
		Window:   window,
		Renderer: renderer,
	}
}

// Cleanup releases all resources.
func (ctx *GameContext) Cleanup() {
	ctx.Renderer.Destroy()
	ctx.Window.Destroy()
	Casts.CleanupFonts()
	sdl.Quit()
}
