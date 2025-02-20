package model

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Player struct {
	X, Y      float64
	Angle     float64
	VelocityX float64
	VelocityY float64
	Walking   bool
	Running   bool
	BobOffset float64
	BobCycle  float64
}

type RenderSlice struct {
	DstRect  *sdl.Rect
	Darkness uint8
	Color    sdl.Color
	WallType int
	TexCoord int32
	Distance float64
}

type RayHit struct {
	Distance   float64
	WallType   int
	HitPointX  float64
	HitPointY  float64
	IsVertical bool
}

type TextureMap struct {
	Textures map[int]*sdl.Texture
}

type NPC struct {
	X, Y        float64
	Texture     int
	Width       float64
	Height      float64
	Distance    float64
	DialogText  string
	ShowDialog  bool
	DialogTimer int
	Hitbox      struct {
		Radius float64
	}
}

type NPCManager struct {
	NPCs []*NPC
}

type DialogRenderer struct {
	Font   *ttf.Font
	Loaded bool
}
