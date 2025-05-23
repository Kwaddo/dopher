package global

import (
	"sync"
	"time"

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
	// The dash cooldown, last dash pressed, and if the player is dashing are for the dash mechanic.
	DashCooldown    int
	LastDashPressed bool
	IsDashing       bool
	// The state of the gun held by the player.
	Gun *GunState
}

// The state of the gun held by the player.
type GunState struct {
	// The current weapon is the weapon the player is holding.
	CurrentWeapon int
	// The last fired time is the last time the player fired the weapon.
	LastFired time.Time
	// IsFiring is if the player is firing the weapon.
	IsFiring bool
	// MuzzleFlash is if the muzzle flash is shown.
	MuzzleFlash bool
	// How long the flash is being shown for.
	FlashTimer int
	// The damage of the weapon and the fire rate.
	Damage   int
	FireRate time.Duration
	// The amount ammo and max ammo of the weapon.
	Ammo    int
	MaxAmmo int
}

// The struct of the map itself.
type Map struct {
	// The maps is a 3D array, within are 2D arrays that has the numbers representing the walls.
	Maps [][][]int
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

// To render in batches.
type RenderBatch struct {
	// The texture that's to be rendered.
	Texture *sdl.Texture
	// The darkness that's to be rendered.
	Darkness uint8
	// The according slices.
	Slices []*sdl.Rect
	// The rectangle slices.
	SrcRects []*sdl.Rect
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
	// New dialogue tree field
	DialogueTree *DialogueTree
	// The hitbox radius of the NPC.
	Hitbox struct {
		Radius float64
	}
	// Enemy behavior, starting with if the NPC is an enemy.
	IsEnemy bool
	// State of the enemy, can be idle, chasing, or attacking.
	State int
	// The detection radius of the NPC.
	DetectionRadius float64
	// How fast the NPC moves.
	Speed float64
	// The last move time when the NPC moved.
	LastMoveTime int
	// The path blocked time when the NPC is blocked.
	PathBlockedTime int
	// The last direction the NPC moved.
	LastDirection struct {
		X, Y float64
	}
	// The remaining health of the NPC.
	Health int
	// The max health of the NPC.
	MaxHealth int
	// If the NPC is alive or not.
	IsAlive bool
	// The map that the NPC is currently in.
	MapIndex int
}

// DialogueNode represents a single node in a dialogue tree
type DialogueNode struct {
	// The ID of the node.
	ID string
	// The according dialogue text of the node.
	Text string
	// The ID of the next node, hence making it one.
	NextID string
	// To connect the nodes directly.
	OnEnter func(*NPC)
}

// DialogueTree contains all dialogue nodes for an NPC
type DialogueTree struct {
	// The nodes mapped.
	Nodes map[string]*DialogueNode
	// The ID of the current node.
	CurrentNodeID string
	// If the dialogue tree is active or not.
	IsActive bool
	// If ready to advance to the next node or not.
	ReadyToAdvance bool
	// The grace for the start period of the text box.
	GraceStartTime int64
	// The altogether grace period for the dialogue tree.
	GracePeriod int64
	// How many characters to show at a time.
	CharsToShow int
	// The speed of the text animation.
	TextSpeed int
	// The time of when the last character appears.
	LastCharTime int64
	// IF the text was fully shown or not.
	TextFullyShown bool
	// If the according action is held down or not.
	KeyWasDown bool
}

// The struct managing all NPCs.
type NPCManager struct {
	// The array of NPCs so that we can manage multiple in a game.
	NPCs []*NPC
}

// FontManager handles loading and caching fonts at different sizes
type FontManager struct {
	// The font path is the path to the font file.
	Path string
	// The font cache is the cache of fonts, based off of the size.
	Cache map[int]*ttf.Font
	// The mutex is the mutex for the font manager, for no issues.
	Mutex sync.RWMutex
	// If the font is initialized or not.
	IsInitialized bool
}

// The dialog renderer is the renderer for the dialog box.
type DialogueRenderer struct {
	// If loaded or not.
	Loaded bool
	// The text texture chache.
	TextCache map[string]*TextureCacheEntry
}

// The cache for the text box.
type TextureCacheEntry struct {
	// The texture itself.
	Texture *sdl.Texture
	// The width and height of it.
	Width, Height int32
	// When it was last used.
	LastUsed int64
}

// GameState tracks the current state of the game.
type GameState struct {
	// The game is paused or not.
	IsPaused bool
	// The game is in the main menu or not.
	InMainMenu    bool
	InOptionsMenu bool
}

// PauseMenu manages the pause menu state.
type PauseMenu struct {
	// The current option selected.
	CurrentOption int
	// The created options for the pause menu.
	Options []string
}

type MainMenu struct {
	CurrentOption int
	Options       []string
}

// OptionsMenu manages the options menu state and settings.
type OptionsMenu struct {
	CurrentOption int
	Options       []string
	Settings      map[string]*Setting
}

// Setting represents a configurable game setting.
type Setting struct {
	Value       float64
	Min         float64
	Max         float64
	Step        float64
	DisplayFunc func(float64) string
}

// GameContext holds all initialized game resources.
type GameContext struct {
	// The window of the game.
	Window *sdl.Window
	// The renderer of the game.
	Renderer *sdl.Renderer
}
