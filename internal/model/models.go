package model

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// The player's values.
type Player struct {
	// X and Y are the player's position.
	X, Y float64
	// Angle is of where the player is looking.
	Angle float64
	// VelocityX and VelocityY are the player's movement speed based off of the direction.
	VelocityX float64
	VelocityY float64
	// Walking, Running, and Crouching are the player's movement states.
	Walking   bool
	Running   bool
	Crouching bool
	// The height is the player's height depending on the action, while the default height is the player's default static height.
	Height        float64
	DefaultHeight float64
	// The bob offset and cycle are for the player's bobbing animation.
	BobOffset float64
	BobCycle  float64
}

// The struct of the map itself.
type Map struct {
	// The map is a 2D array, with the numbers representing the walls.
	WorldMap [][]int
}

// The render slice is the slice of the screen that will be rendered.
type RenderSlice struct {
	// The distance rectangle is the according rectangle of which the wall will be rendered.
	DstRect *sdl.Rect
	// The darkness is the dark fog in the distance.
	Darkness uint8
	// The color is the color of the wall, in case it doesn't have a texture.
	Color sdl.Color
	// The wall type is which type of wall it is, for the texture to apply accordingly.
	WallType int
	// The texture coordinate is the texture's position on the wall.
	TexCoord int32
	// The distance is the pure distance from the player to the wall.
	Distance float64
}

// RayHit is the hit of the ray on the wall.
type RayHit struct {
	// The distance is the distance from the player to the wall.
	Distance float64
	// The wall type is which type of wall it is, for the texture to apply accordingly.
	WallType int
	// The hit point X and Y are the exact point where the ray hit the wall.
	HitPointX float64
	HitPointY float64
	// IsVertical is if the wall is vertical or not.
	IsVertical bool
}

// TextureMap is the map of textures that will be used in the game.
type TextureMap struct {
	// Textures is the map of textures, based off of numbers and local textures.
	Textures map[int]*sdl.Texture
}

// The struct of an NPC.
type NPC struct {
	// The positiion of the NPC.
	X, Y float64
	// The texture placed on the NPC.
	Texture int
	// The width and height of the NPC.
	Width  float64
	Height float64
	// The distance from the player to the NPC.
	Distance float64
	// The dialog text and if the dialog should be shown, and how long the dialog should be shown.
	DialogText  string
	ShowDialog  bool
	DialogTimer int
	// The hitbox radius of the NPC.
	Hitbox struct {
		Radius float64
	}
}

// The struct managing all NPCs.
type NPCManager struct {
	// The array of NPCs so that we can manage multiple in a game.
	NPCs []*NPC
}

// The dialog renderer is the renderer for the dialog box.
type DialogRenderer struct {
	// The font and if it's loaded or not.
	Font   *ttf.Font
	Loaded bool
}
