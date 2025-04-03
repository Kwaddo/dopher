package main

import (
	MC "doom/internal/character/player"
	Game "doom/internal/core"
	DM "doom/internal/global"
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

	DM.GlobalGameState = DM.GameState{
		InMainMenu: true,
		IsPaused:   false,
	}

	// The gameloop!
	Game.RunGameLoop(gameCtx.Renderer, player)
}
