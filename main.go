package main

import (
	DM "doom/internal/constants"
	Game "doom/internal/funcs"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Doom in Go", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		DM.ScreenWidth, DM.ScreenHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	player := &DM.Player{
		X:     150,
		Y:     150,
		Angle: 0,
	}

	// Remove the goroutine and channel, run the game loop directly
	Game.GameLoop(renderer, player, nil)
}
