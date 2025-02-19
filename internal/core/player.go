package core

import (
	DM "doom/internal/model"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type Player DM.Player

func (p *Player) Movement(state []uint8) bool {
	// Store current position for collision check
	oldX := p.X
	oldY := p.Y

	speedMultiplier := 1.0
	isMoving := state[sdl.SCANCODE_W] == 1 || state[sdl.SCANCODE_S] == 1 ||
		state[sdl.SCANCODE_A] == 1 || state[sdl.SCANCODE_D] == 1

	if state[sdl.SCANCODE_LSHIFT] == 1 && isMoving {
		speedMultiplier = DM.SprintMultiplier
		p.Running = true
	} else {
		p.Running = false
	}

	// Compute acceleration and max speed based on sprinting state
	Acceleration := DM.BaseAcceleration * speedMultiplier
	MaxSpeed := DM.BaseMaxSpeed * speedMultiplier
	// Compute directional vectors
	forwardX := math.Cos(p.Angle)
	forwardY := math.Sin(p.Angle)
	strafeX := math.Cos(p.Angle + math.Pi/2)
	strafeY := math.Sin(p.Angle + math.Pi/2)

	// Apply acceleration based on key presses
	if state[sdl.SCANCODE_W] == 1 { // Forward
		p.VelocityX += forwardX * Acceleration
		p.VelocityY += forwardY * Acceleration
	}
	if state[sdl.SCANCODE_S] == 1 { // Backward
		p.VelocityX -= forwardX * Acceleration
		p.VelocityY -= forwardY * Acceleration
	}
	if state[sdl.SCANCODE_A] == 1 { // Strafe left
		p.VelocityX -= strafeX * Acceleration
		p.VelocityY -= strafeY * Acceleration
	}
	if state[sdl.SCANCODE_D] == 1 { // Strafe right
		p.VelocityX += strafeX * Acceleration
		p.VelocityY += strafeY * Acceleration
	}
	// Rotate with left/right arrows
	if state[sdl.SCANCODE_LEFT] == 1 {
		p.Angle -= DM.RotateSpeed
	}
	if state[sdl.SCANCODE_RIGHT] == 1 {
		p.Angle += DM.RotateSpeed
	}
	if state[sdl.SCANCODE_ESCAPE] == 1 || state[sdl.SCANCODE_Q] == 1 {
		return true
	}
	if isMoving {
		p.Walking = true
	} else {
		p.Walking = false
		p.Running = false // Ensure running is false when not moving
	}

	// Apply friction when no keys are pressed
	p.VelocityX *= (1 - DM.Friction)
	p.VelocityY *= (1 - DM.Friction)

	// Limit max speed
	speed := math.Hypot(p.VelocityX, p.VelocityY)
	if speed > MaxSpeed {
		scale := MaxSpeed / speed
		p.VelocityX *= scale
		p.VelocityY *= scale
	}

	// Calculate new position
	newX := p.X + p.VelocityX
	newY := p.Y + p.VelocityY

	// Collision check before updating position
	collidesX := CheckCollision(newX, p.Y) // Only X movement
	collidesY := CheckCollision(p.X, newY) // Only Y movement

	if collidesX && collidesY {
		p.VelocityX = 0
		p.VelocityY = 0
	} else if collidesX {
		p.Y = newY - 0.5
		p.VelocityX = 0
	} else if collidesY {
		p.X = newX - 0.5
		p.VelocityY = 0
	} else {
		p.X = newX
		p.Y = newY
	}

	// After calculating new position, check for NPC collisions
	if npcManager := GlobalNPCManager; npcManager != nil {
		if npcManager.CheckNPCCollision(p.X, p.Y) {
			// Collision detected, revert position
			p.X = oldX
			p.Y = oldY
			// Also reset velocities for smoother collision response
			p.VelocityX = 0
			p.VelocityY = 0
		}
	}

	return state[sdl.SCANCODE_ESCAPE] == 1
}

func LERP(start, end, t float64) float64 {
	return start + t*(end-start)
}
