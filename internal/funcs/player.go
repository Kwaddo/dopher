package funcs

import (
	"math"
	"github.com/veandco/go-sdl2/sdl"
	DM "doom/internal/constants"
)

type Player struct {
	X, Y      float64
	Angle     float64
	VelocityX float64
	VelocityY float64
}

func (player *Player) UpdateMovement(keys []uint8) {
	speedMultiplier := 1.0
    if keys[sdl.SCANCODE_LSHIFT] == 1 {
        speedMultiplier = DM.SprintMultiplier
    }

    // Compute acceleration and max speed based on sprinting state
    Acceleration := DM.BaseAcceleration * speedMultiplier
    MaxSpeed := DM.BaseMaxSpeed * speedMultiplier
	// Compute directional vectors
	forwardX := math.Cos(player.Angle)
	forwardY := math.Sin(player.Angle)
	strafeX := math.Cos(player.Angle + math.Pi/2)
	strafeY := math.Sin(player.Angle + math.Pi/2)

	// Apply acceleration based on key presses
	if keys[sdl.SCANCODE_W] == 1 { // Forward
		player.VelocityX += forwardX * Acceleration
		player.VelocityY += forwardY * Acceleration
	}
	if keys[sdl.SCANCODE_S] == 1 { // Backward
		player.VelocityX -= forwardX * Acceleration
		player.VelocityY -= forwardY * Acceleration
	}
	if keys[sdl.SCANCODE_A] == 1 { // Strafe left
		player.VelocityX -= strafeX * Acceleration
		player.VelocityY -= strafeY * Acceleration
	}
	if keys[sdl.SCANCODE_D] == 1 { // Strafe right
		player.VelocityX += strafeX * Acceleration
		player.VelocityY += strafeY * Acceleration
	}

	// Apply friction when no keys are pressed
	player.VelocityX *= (1 - DM.Friction)
	player.VelocityY *= (1 - DM.Friction)

	// Limit max speed
	speed := math.Hypot(player.VelocityX, player.VelocityY)
	if speed > MaxSpeed {
		scale := MaxSpeed / speed
		player.VelocityX *= scale
		player.VelocityY *= scale
	}

	// Calculate new position
	newX := player.X + player.VelocityX
	newY := player.Y + player.VelocityY

	// Collision check before updating position
	collidesX := CheckCollision(newX, player.Y) // Only X movement
	collidesY := CheckCollision(player.X, newY) // Only Y movement

	if collidesX && collidesY {
		player.VelocityX = 0
		player.VelocityY = 0
	} else if collidesX {
		player.Y = newY-0.5
		player.VelocityX = 0
	} else if collidesY {
		player.X = newX-0.5
		player.VelocityY = 0
	} else {
		player.X = newX
		player.Y = newY
	}
}
