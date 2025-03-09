package npc

import (
	DM "doom/internal/model"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

// CheckNPCCollision checks if the player is colliding with an NPC.
func (nm *NPCManager) CheckNPCCollision(x, y float64) bool {
	for _, npc := range nm.NPCs {
		dx := x - npc.X
		dy := y - npc.Y
		distSquared := dx*dx + dy*dy
		if distSquared < npc.Hitbox.Radius*npc.Hitbox.Radius {
			return true
		}
	}
	return false
}

// CheckInteraction checks if the player is interacting with an NPC.
func (nm *NPCManager) CheckInteraction(playerX, playerY, playerAngle float64, keyState []uint8) {
	for i, npc := range nm.NPCs {
		if npc.ShowDialog && npc.DialogueTree != nil && npc.DialogueTree.IsActive {
			DM.InteractingNPC = i
			if keyState[sdl.SCANCODE_ESCAPE] == 1 {
				nm.EndDialogue(i)
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
			eKeyPressed := keyState[sdl.SCANCODE_E] == 1
			keyJustPressed := eKeyPressed && !npc.DialogueTree.KeyWasDown
			npc.DialogueTree.KeyWasDown = eKeyPressed
			if keyState[sdl.SCANCODE_RETURN] == 1 {
				nm.EndDialogue(i)
				return
			}
			if keyJustPressed {
				if !npc.DialogueTree.TextFullyShown {
					npc.DialogueTree.CharsToShow = len(npc.DialogText)
					npc.DialogueTree.TextFullyShown = true
					npc.DialogueTree.ReadyToAdvance = true
					return
				} else if npc.DialogueTree.ReadyToAdvance {
					currentNode := npc.DialogueTree.Nodes[npc.DialogueTree.CurrentNodeID]
					if currentNode == nil {
						nm.EndDialogue(i)
						return
					}
					nextNodeExists := currentNode.NextID != "" && npc.DialogueTree.Nodes[currentNode.NextID] != nil
					if nextNodeExists {
						AdvanceToNextDialogue(npc, currentNode.NextID)
					} else {
						nm.EndDialogue(i)
					}
					return
				}
			}
			if keyState[sdl.SCANCODE_ESCAPE] == 1 {
				nm.EndDialogue(i)
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
			if checkY < 0 || checkY >= len(DM.GlobalMap.WorldMap) ||
				checkX < 0 || checkX >= len(DM.GlobalMap.WorldMap[0]) {
				continue
			}
			if DM.GlobalMap.WorldMap[checkY][checkX] > 0 {
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
