package main

import (
	MC "doom/internal/char/player"
	Game "doom/internal/core"
)

func main() {
	// Initialize all required resources
	gameCtx := Game.InitializeAll()
	defer gameCtx.Cleanup()

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
	Game.GameLoop(gameCtx.Renderer, player)
}
