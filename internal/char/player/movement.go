package player

import (
	NPC "doom/internal/char/npc"
	DM "doom/internal/model"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type Player DM.Player

func (p *Player) Movement(state []uint8, npcManager *NPC.NPCManager) bool {
	oldX := p.X
	oldY := p.Y

	Acceleration, MaxSpeed, isMoving := AccelerationAndMaxSpeed(p, state)
	Input(p, state, Acceleration)
	Rotation(p, state)

	if state[sdl.SCANCODE_ESCAPE] == 1 || state[sdl.SCANCODE_Q] == 1 {
		return true
	}
	if state[sdl.SCANCODE_E] == 1 {
		npcManager.CheckInteraction(p.X, p.Y, p.Angle)
	}

	// Update walking state
	p.Walking = isMoving
	if !isMoving {
		p.Running = false
	}

	FrictionAndLimitSpeed(p, MaxSpeed)

	// Compute new position
	newX := p.X + p.VelocityX
	newY := p.Y + p.VelocityY

	// Collision checks
	collidesX := CheckCollision(newX, p.Y)
	collidesY := CheckCollision(p.X, newY)
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

	// NPC collision
	if npcManager := NPC.GlobalNPCManager; npcManager != nil {
		if npcManager.CheckNPCCollision(p.X, p.Y) {
			p.X = oldX
			p.Y = oldY
			p.VelocityX = 0
			p.VelocityY = 0
		}
	}

	// Head bobbing logic
	if p.Walking {
		bobSpeed := 0.1
		if p.Running {
			bobSpeed = 0.15
		}
		bobAmplitude := 7.5

		p.BobCycle += bobSpeed
		p.BobOffset = math.Sin(p.BobCycle) * bobAmplitude
	} else {
		p.BobOffset *= 0.8
		p.BobCycle = 0
	}

	return state[sdl.SCANCODE_ESCAPE] == 1
}


func AccelerationAndMaxSpeed(p *Player, state []uint8) (float64, float64, bool) {
	isMoving := state[sdl.SCANCODE_W] == 1 ||
		state[sdl.SCANCODE_S] == 1 ||
		state[sdl.SCANCODE_A] == 1 ||
		state[sdl.SCANCODE_D] == 1

	speedMultiplier := 1.0
	if state[sdl.SCANCODE_LSHIFT] == 1 && isMoving {
		speedMultiplier = DM.SprintMultiplier
		p.Running = true
	} else {
		p.Running = false
	}
	Acceleration := DM.BaseAcceleration * speedMultiplier
	MaxSpeed := DM.BaseMaxSpeed * speedMultiplier

	return Acceleration, MaxSpeed, isMoving
}

func Input(p *Player, state []uint8, Acceleration float64) {
	forwardX := math.Cos(p.Angle)
	forwardY := math.Sin(p.Angle)
	strafeX := math.Cos(p.Angle + math.Pi/2)
	strafeY := math.Sin(p.Angle + math.Pi/2)

	if state[sdl.SCANCODE_W] == 1 {
		p.VelocityX += forwardX * Acceleration
		p.VelocityY += forwardY * Acceleration
	}
	if state[sdl.SCANCODE_S] == 1 {
		p.VelocityX -= forwardX * Acceleration
		p.VelocityY -= forwardY * Acceleration
	}
	if state[sdl.SCANCODE_A] == 1 {
		p.VelocityX -= strafeX * Acceleration
		p.VelocityY -= strafeY * Acceleration
	}
	if state[sdl.SCANCODE_D] == 1 {
		p.VelocityX += strafeX * Acceleration
		p.VelocityY += strafeY * Acceleration
	}
}

func Rotation(p *Player, state []uint8) {
	if state[sdl.SCANCODE_LEFT] == 1 {
		p.Angle -= DM.RotateSpeed
	}
	if state[sdl.SCANCODE_RIGHT] == 1 {
		p.Angle += DM.RotateSpeed
	}
}

func FrictionAndLimitSpeed(p *Player, MaxSpeed float64) {
	// Apply friction
	p.VelocityX *= (1 - DM.Friction)
	p.VelocityY *= (1 - DM.Friction)

	// Limit max speed
	speed := math.Hypot(p.VelocityX, p.VelocityY)
	if speed > MaxSpeed {
		scale := MaxSpeed / speed
		p.VelocityX *= scale
		p.VelocityY *= scale
	}
}

// Example linear interpolation
func LERP(start, end, t float64) float64 {
	return start + t*(end-start)
}
