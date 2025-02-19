package main

import (
	Game "doom/internal/core"
	DM "doom/internal/model"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Dopher Engine", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
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

	player := &Game.Player{
		X:     150,
		Y:     150,
		Angle: 0,
	}
	Game.GameLoop(renderer, player)
}
