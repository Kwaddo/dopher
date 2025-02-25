package main

import (
	MC "doom/internal/char/player"
	Game "doom/internal/core"
	DM "doom/internal/model"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	DM.ScreenWidth = 1500
	DM.ScreenHeight = 900

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Dopher Engine", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		int32(DM.ScreenWidth), int32(DM.ScreenHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	player := &MC.Player{
		X:     150,
		Y:     150,
		Angle: 0,
		DefaultHeight: 64.0,
		Height: 64.0,
	}
	Game.GameLoop(renderer, player)
}
