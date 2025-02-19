package constants

import "github.com/veandco/go-sdl2/sdl"

type Player struct {
	X, Y      float64
	Angle     float64
	VelocityX float64
	VelocityY float64
	Walking   bool
	Running   bool
}

type RenderSlice struct {
	DstRect  *sdl.Rect
	Darkness uint8
	Color    sdl.Color
	WallType int
	TexCoord int32
	Distance float64
}

type NPC struct {
	X, Y     float64
	Texture  int
	Width    float64
	Height   float64
	Distance float64
	Hitbox   struct {
		Radius float64
	}
}

type NPCManager struct {
	NPCs []*NPC
}
