package model

import "math"

const (
	FOV              = math.Pi / 3.5
	NumRays          = 120
	MaxDepth         = 800.0
	RotateSpeed      = 0.1
	CollisionBuffer  = 20.0
	BaseAcceleration = 0.5
	BaseMaxSpeed     = 5.0
	Friction         = 0.15
	SprintMultiplier = 1.8
	MaxDarkness      = 255
)

var (
	ScreenWidth  float64
	ScreenHeight float64
)
