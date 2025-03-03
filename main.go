package main

import (
	MC "doom/internal/char/player"
	Game "doom/internal/core"
	DM "doom/internal/model"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	// Original screen size, without going fullscreen.
	DM.ScreenWidth = 1500
	DM.ScreenHeight = 900

	// Initialize everything within SDL for the rendering to work.
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	if err := DM.InitFonts(); err != nil {
		panic(err)
	}
	defer DM.CleanupFonts()

	// Create a window with the title "Dopher Engine" and the screen size.
	window, err := sdl.CreateWindow("Dopher Engine", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		int32(DM.ScreenWidth), int32(DM.ScreenHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	// Create a renderer, and it has hardware acceleration.
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	// Initialize the player with the position, angle, and height.
	player := &MC.Player{
		X:               150,
		Y:               150,
		Angle:           0,
		DefaultHeight:   64.0,
		Height:          64.0,
		DashCooldown:    0,
		LastDashPressed: false,
	}

	// The gameloop!
	Game.GameLoop(renderer, player)
}
