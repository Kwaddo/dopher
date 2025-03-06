package player

import (
	NPC "doom/internal/character/npc"
	DM "doom/internal/model"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

// MovementInputs handles basic player input for movement.
func MovementInputs(p *Player, state []uint8, Acceleration float64) {
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

func FireWeaponInput(p *Player, npcManager *NPC.NPCManager, state []uint8) {
	if state[sdl.SCANCODE_UP] == 1 {
		if p.FireGun() {
			FireWeapon(p, npcManager)
		}
	}
}

// Rotation handles player rotation based on input.
func Rotation(p *Player, state []uint8) {
	if state[sdl.SCANCODE_LEFT] == 1 {
		p.Angle -= DM.RotateSpeed
	}
	if state[sdl.SCANCODE_RIGHT] == 1 {
		p.Angle += DM.RotateSpeed
	}
}
