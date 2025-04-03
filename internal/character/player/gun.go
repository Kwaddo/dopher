package player

import (
	NPC "doom/internal/character/npc"
	DM "doom/internal/global"
	Casts "doom/internal/graphics/casting"
	"math"
	"time"
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
	if DM.CountdownFreeze {
		return false
	}
	if DM.CurrentMap != 1 {
		return false
	}
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
	if DM.CurrentMap != 1 {
		p.Gun.IsFiring = false
		p.Gun.MuzzleFlash = false
		p.Gun.FlashTimer = 0
		return
	}
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

// Weapon firing handler.
func FireWeapon(p *Player, npcManager *NPC.NPCManager) {
	if DM.CurrentMap != 1 {
		return
	}
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
