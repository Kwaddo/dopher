package player

import (
	NPC "doom/internal/character/npc"
	DM "doom/internal/model"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type Player DM.Player

// Actions handles player movement, action, weapons, and collision detection.
func (p *Player) Actions(state []uint8, npcManager *NPC.NPCManager) bool {
	if DM.GlobalGameState.IsPaused {
		return false
	}
	if state[sdl.SCANCODE_Q] == 1 {
		return true
	}
	oldX := p.X
	oldY := p.Y
	Acceleration, MaxSpeed, isMoving := AccelerationAndMaxSpeed(p, state)
	MovementInputs(p, state, Acceleration)
	Rotation(p, state)
	Dash(p, state)
	if state[sdl.SCANCODE_ESCAPE] == 1 || state[sdl.SCANCODE_Q] == 1 {
		return true
	}

	// Check for both direct interaction via E and
	// continuous dialogue option selection
	if state[sdl.SCANCODE_E] == 1 {
		npcManager.CheckInteraction(p.X, p.Y, p.Angle, state)
	} else {
		// Always check for dialogue option inputs when an NPC dialogue is active
		npcManager.CheckDialogueInput(state)
	}

	p.Walking = isMoving
	if !isMoving {
		p.Running = false
	}
	FrictionAndLimitSpeed(p, MaxSpeed)
	newX := p.X + p.VelocityX
	newY := p.Y + p.VelocityY
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
	if npcManager := NPC.GlobalNPCManager; npcManager != nil {
		if npcManager.CheckNPCCollision(p.X, p.Y) {
			p.X = oldX
			p.Y = oldY
			p.VelocityX = 0
			p.VelocityY = 0
		}
	}
	if p.Walking && DM.HeadBobbingEnabled {
		bobSpeed := 0.1
		if p.Running {
			bobSpeed = 0.15
		}
		bobAmplitude := 7.5
		p.BobCycle += bobSpeed
		p.BobOffset = math.Sin(p.BobCycle) * bobAmplitude
	} else {
		p.BobOffset *= 0.8
		if !DM.HeadBobbingEnabled {
			p.BobOffset = 0
		}
		p.BobCycle = 0
	}
	UpdatePlayerHeight(p)
	FireWeaponInput(p, npcManager, state)
	p.UpdateGunState()
	return false
}

// LERP stands for Linear Interpolation, which is a way to smoothly transition between two values.
func LERP(start, end, t float64) float64 {
	return start + t*(end-start)
}

// UpdatePlayerHeight updates the player's height based on the player's crouching state.
func UpdatePlayerHeight(p *Player) {
	if p.Crouching {
		targetHeight := p.DefaultHeight * 0.5
		p.Height = LERP(p.Height, targetHeight, 0.2)
	} else {
		p.Height = LERP(p.Height, p.DefaultHeight, 1)
	}
}
