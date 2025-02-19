package constants

type Player struct {
	X, Y      float64
	Angle     float64
	VelocityX float64
	VelocityY float64
	Walking   bool
	Running   bool
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
