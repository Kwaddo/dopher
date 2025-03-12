package model

import (
	"math"
)

const (
	// MaxDepth is the maximum depth of the ray.
	MaxDepth = 800.0
	// CollisionBuffer is the buffer for collision detection.
	CollisionBuffer = 20.0
	// BaseAcceleration is the base acceleration of the player.
	BaseAcceleration = 0.5
	// BaseFriction is the base friction of the player to slow down.
	Friction = 0.15
	// SprintMultiplier is the multiplier for sprinting.
	SprintMultiplier = 1.8
	// MaxDarkness is the maximum constant of a far darkness.
	MaxDarkness = 255
	// Checkers for if an enemy is idle, chasing, or attacking.
	EnemyStateIdle    = 0
	EnemyStateChasing = 1
	EnemyStateAttack  = 2
	// The speed of the LERP.
	LerpSpeed = 0.15
)

var (
	// The width of the screen.
	ScreenWidth = 1500.0
	// The height of the screen.
	ScreenHeight = 900.0
	// Number of rays to cast, think of it as the graphics.
	NumRays = 150
	// The global frame count.
	GlobalFrameCount = 0
	// FOV variables for current, target, and dynamic.
	CurrentFOV = FOV
	TargetFOV  = FOV
	DynamicFOV = FOV
	// The pointer to the current FOV.
	ZBuffer     []float64
	ShowMiniMap = true
	ShowMegaMap = false
	// The states of the game.
	GlobalGameState = GameState{
		IsPaused: false,
	}
	// Global usage of textures
	GlobalTextures *TextureMap
	// The channel for rendering slices.
	RenderChan chan []*RenderSlice
	// If the head bobbing is enabled or not.
	HeadBobbingEnabled = true
	// FOV is the field of view of the player.
	FOV = math.Pi / 3.5
	// RotateSpeed is the speed of the player rotation.
	RotateSpeed = 0.1
	// BaseMaxSpeed is the base maximum speed of the player.
	BaseMaxSpeed = 5.0
	// Checks if buffers need to be recreated.
	NeedToRecreateBuffers = false
	// The last state of the key pressed.
	LastKeyState byte = 0
	// Shows if interacting an NPC.
	InteractingNPC int = -1
	// Which map the player is currently at.
	CurrentMap int = 0
)
