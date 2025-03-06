package player

import (
	NPC "doom/internal/character/npc"
	Casts "doom/internal/graphics/casting"
	DM "doom/internal/model"
	"math"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

// InitializeGun sets up the player's initial gun state.
func (p *Player) InitializeGun() {
	p.Gun = &DM.GunState{
		CurrentWeapon: 0,
		LastFired:     time.Now(),
		IsFiring:      false,
		MuzzleFlash:   false,
		FlashTimer:    0,
		Damage:        25,
		FireRate:      time.Millisecond * 100,
		Ammo:          50,
		MaxAmmo:       100,
	}
}

// FireGun attempts to fire the gun and returns whether it was successful or not.
func (p *Player) FireGun() bool {
	now := time.Now()
	if now.Sub(p.Gun.LastFired) < p.Gun.FireRate {
		return false
	}
	if p.Gun.Ammo <= 0 {
		return false
	}
	p.Gun.LastFired = now
	p.Gun.IsFiring = true
	p.Gun.MuzzleFlash = true
	p.Gun.FlashTimer = 5
	p.Gun.Ammo--
	return true
}

// UpdateGunState updates gun animation states.
func (p *Player) UpdateGunState() {
	if p.Gun.FlashTimer > 0 {
		p.Gun.FlashTimer--
		if p.Gun.FlashTimer <= 0 {
			p.Gun.MuzzleFlash = false
		}
	}
	if p.Gun.IsFiring && time.Since(p.Gun.LastFired) > time.Millisecond*100 {
		p.Gun.IsFiring = false
	}
}

// RenderGun renders the gun on the screen.
func RenderGun(renderer *sdl.Renderer, player *Player, textures *DM.TextureMap) {
	gunTexture := textures.Textures[6]
	gunWidth := DM.ScreenWidth * 0.4
	gunHeight := gunWidth * 0.5
	offsetY := 0.0
	if player.Walking {
		offsetY = player.BobOffset * 0.5
	}
	recoilOffset := 0.0
	if player.Gun.IsFiring {
		recoilOffset = 10.0
	}
	dstRect := &sdl.Rect{
		X: int32((DM.ScreenWidth - gunWidth) / 2),
		Y: int32(DM.ScreenHeight - gunHeight + offsetY + recoilOffset),
		W: int32(gunWidth),
		H: int32(gunHeight),
	}
	renderer.Copy(gunTexture, nil, dstRect)
	if player.Gun.MuzzleFlash {
		flashTexture := textures.Textures[7]
		flashWidth := gunWidth * 0.3
		flashHeight := flashWidth
		flashRect := &sdl.Rect{
			X: int32(DM.ScreenWidth / 2),
			Y: int32(DM.ScreenHeight - gunHeight*0.8 + offsetY),
			W: int32(flashWidth),
			H: int32(flashHeight),
		}
		renderer.Copy(flashTexture, nil, flashRect)
	}
}

// Weapon firing handler.
func FireWeapon(p *Player, npcManager *NPC.NPCManager) {
	rayHit := Casts.CastRay(p.X, p.Y, p.Angle)
	hitDistance := rayHit.Distance
	if hitDistance > 500 {
		hitDistance = 500
	}
	for i, npc := range npcManager.NPCs {
		if !npc.IsEnemy || !npc.IsAlive {
			continue
		}
		dx := npc.X - p.X
		dy := npc.Y - p.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist > hitDistance {
			continue
		}
		angleToNPC := math.Atan2(dy, dx)
		angleDiff := math.Abs(angleToNPC - p.Angle)
		if angleDiff > math.Pi {
			angleDiff = 2*math.Pi - angleDiff
		}
		if angleDiff < 0.17 {
			wallHit := Casts.CastRay(p.X, p.Y, angleToNPC)
			if wallHit.Distance > dist {
				npcManager.DamageEnemy(i, p.Gun.Damage)
				break
			}
		}
	}
}
