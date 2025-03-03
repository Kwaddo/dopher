package model

import (
	"github.com/veandco/go-sdl2/ttf"
	"math"
)

const (
	// FOV is the field of view of the player.
	FOV = math.Pi / 3.5
	// NumRays is the number of rays to cast.
	NumRays = 120
	// MaxDepth is the maximum depth of the ray.
	MaxDepth = 800.0
	// RotateSpeed is the speed of the player rotation.
	RotateSpeed = 0.1
	// CollisionBuffer is the buffer for collision detection.
	CollisionBuffer = 20.0
	// BaseAcceleration is the base acceleration of the player.
	BaseAcceleration = 0.5
	// BaseMaxSpeed is the base maximum speed of the player.
	BaseMaxSpeed = 5.0
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
	ScreenWidth float64
	// The height of the screen.
	ScreenHeight float64
	// The states of the game.
	GlobalGameState = GameState{
		IsPaused: false,
	}
	// The global variable of the pause menu.
	GlobalPauseMenu *PauseMenu
	// The font manager of the game.
	GlobalFontManager = &FontManager{
		Path:  "assets/font/dogicapixel.ttf",
		Cache: make(map[int]*ttf.Font),
	}
)
