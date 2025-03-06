package player

import (
	DM "doom/internal/model"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

// AccelerationAndMaxSpeed returns the player's acceleration and max speed based on the player's movement state, as well as handling the speed of the crouching.
func AccelerationAndMaxSpeed(p *Player, state []uint8) (float64, float64, bool) {
	isMoving := (state[sdl.SCANCODE_W] == 1 || state[sdl.SCANCODE_A] == 1 ||
		state[sdl.SCANCODE_S] == 1 || state[sdl.SCANCODE_D] == 1)
	speedMultiplier := 1.0
	if state[sdl.SCANCODE_LSHIFT] == 1 && isMoving && !p.Crouching {
		p.Running = true
		speedMultiplier = 1.5
	} else {
		p.Running = false
	}
	if state[sdl.SCANCODE_LCTRL] == 1 {
		p.Crouching = true
		speedMultiplier = 0.5
	} else {
		p.Crouching = false
	}
	Acceleration := DM.BaseAcceleration * speedMultiplier
	MaxSpeed := DM.BaseMaxSpeed * speedMultiplier
	return Acceleration, MaxSpeed, isMoving
}

// Dash handles the player's dash mechanic.
func Dash(p *Player, state []uint8) {
	if p.DashCooldown > 0 {
		p.DashCooldown--
		return
	}
	if state[sdl.SCANCODE_SPACE] == 1 && !p.LastDashPressed {
		speed := math.Hypot(p.VelocityX, p.VelocityY)
		if speed > 0.1 {
			dirX := p.VelocityX / speed
			dirY := p.VelocityY / speed
			dashForce := 50.0
			p.VelocityX += dirX * dashForce
			p.VelocityY += dirY * dashForce
			p.DashCooldown = 30
		}
	}
	p.LastDashPressed = state[sdl.SCANCODE_SPACE] == 1
}

// FrictionAndLimitSpeed applies friction and limits the player's speed.
func FrictionAndLimitSpeed(p *Player, MaxSpeed float64) {
	p.VelocityX *= (1 - DM.Friction)
	p.VelocityY *= (1 - DM.Friction)
	speed := math.Hypot(p.VelocityX, p.VelocityY)
	if speed > MaxSpeed {
		scale := MaxSpeed / speed
		p.VelocityX *= scale
		p.VelocityY *= scale
	}
}
