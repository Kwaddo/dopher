package npc

import (
	DM "doom/internal/global"
	MapModel "doom/internal/mapmodel"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

// Global cooldown for dialogue interaction.
var dialogueCooldown int = 0

// CheckNPCCollision checks if the player is colliding with an NPC and returns both collision status and NPC index.
func (nm *NPCManager) CheckNPCCollision(x, y float64) (bool, int) {
	for i, npc := range nm.NPCs {
		dx := x - npc.X
		dy := y - npc.Y
		distSquared := dx*dx + dy*dy
		if distSquared < npc.Hitbox.Radius*npc.Hitbox.Radius {
			return true, i
		}
	}
	return false, -1
}

// CheckInteraction checks if the player is interacting with an NPC.
func (nm *NPCManager) CheckInteraction(playerX, playerY, playerAngle float64, keyState []uint8) {
	if dialogueCooldown > 0 {
		dialogueCooldown--
		return
	}
	for i, npc := range nm.NPCs {
		if npc.ShowDialog && npc.DialogueTree != nil && npc.DialogueTree.IsActive {
			DM.InteractingNPC = i
			if keyState[sdl.SCANCODE_ESCAPE] == 1 {
				nm.EndDialogue(i)
				dialogueCooldown = 30
			}
			return
		}
		if npc.IsEnemy && !npc.IsAlive {
			continue
		}
		dx := playerX - npc.X
		dy := playerY - npc.Y
		distSquared := dx*dx + dy*dy
		if distSquared < 100*100 && !npc.ShowDialog {
			angleToNPC := math.Atan2(-dy, -dx)
			angleDiff := angleToNPC - playerAngle

			for angleDiff > math.Pi {
				angleDiff -= 2 * math.Pi
			}
			for angleDiff < -math.Pi {
				angleDiff += 2 * math.Pi
			}

			if math.Abs(angleDiff) < math.Pi/4 {
				DM.InteractingNPC = i
				break
			}
		}
	}
	if DM.InteractingNPC >= 0 && keyState[sdl.SCANCODE_E] == 1 {
		nm.StartDialogue(DM.InteractingNPC)
	}
}

// CheckDialogueInput checks for dialogue input when a dialogue is active
func (nm *NPCManager) CheckDialogueInput(keyState []uint8) {
	for i, npc := range nm.NPCs {
		if npc.ShowDialog && npc.DialogueTree != nil && npc.DialogueTree.IsActive {
			if DM.InteractingNPC != i {
				DM.InteractingNPC = i
			}
			eKeyCurrentlyPressed := keyState[sdl.SCANCODE_E] == 1
			keyJustPressed := eKeyCurrentlyPressed && !npc.DialogueTree.KeyWasDown
			npc.DialogueTree.KeyWasDown = eKeyCurrentlyPressed
			if keyState[sdl.SCANCODE_RETURN] == 1 {
				nm.EndDialogue(i)
				return
			}
			if keyJustPressed {
				if !npc.DialogueTree.TextFullyShown {
					npc.DialogueTree.CharsToShow = len(npc.DialogText)
					npc.DialogueTree.TextFullyShown = true
					npc.DialogueTree.ReadyToAdvance = true
				} else if npc.DialogueTree.ReadyToAdvance {
					currentNode := npc.DialogueTree.Nodes[npc.DialogueTree.CurrentNodeID]
					if currentNode == nil {
						npc.ShowDialog = false
						npc.DialogueTree.IsActive = false
						DM.InteractingNPC = -1
						return
					}
					if currentNode.ID == "end" {
						nm.EndDialogue(i)
						dialogueCooldown = 30
						return
					}
					isLastNode := currentNode.NextID == "" || currentNode.NextID == "end"
					if isLastNode && currentNode.NextID == "" {
						nm.EndDialogue(i)
						return
					} else if isLastNode && currentNode.NextID == "end" {
						AdvanceToNextDialogue(npc, currentNode.NextID)
					} else if npc.DialogueTree.Nodes[currentNode.NextID] != nil {
						AdvanceToNextDialogue(npc, currentNode.NextID)
					} else {
						npc.ShowDialog = false
						npc.DialogueTree.IsActive = false
						DM.InteractingNPC = -1
					}
				}
			}
			if keyState[sdl.SCANCODE_ESCAPE] == 1 {
				npc.ShowDialog = false
				npc.DialogueTree.IsActive = false
				DM.InteractingNPC = -1
			}
			return
		}
	}
}

// CheckWallCollision checks if an enemy would collide with a wall
func CheckWallCollision(x, y, radius float64) bool {
	mapX := int(x / 100)
	mapY := int(y / 100)
	for checkY := mapY - 1; checkY <= mapY+1; checkY++ {
		for checkX := mapX - 1; checkX <= mapX+1; checkX++ {
			if checkY < 0 || checkY >= len(MapModel.GlobalMaps.Maps[DM.CurrentMap]) ||
				checkX < 0 || checkX >= len(MapModel.GlobalMaps.Maps[DM.CurrentMap][0]) {
				continue
			}
			if MapModel.GlobalMaps.Maps[DM.CurrentMap][checkY][checkX] > 0 {
				wallMinX := float64(checkX) * 100
				wallMinY := float64(checkY) * 100
				wallMaxX := wallMinX + 100
				wallMaxY := wallMinY + 100
				closestX := math.Max(wallMinX, math.Min(x, wallMaxX))
				closestY := math.Max(wallMinY, math.Min(y, wallMaxY))
				dx := closestX - x
				dy := closestY - y
				if dx*dx+dy*dy < radius*radius {
					return true
				}
			}
		}
	}
	return false
}
