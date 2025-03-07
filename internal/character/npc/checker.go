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
	// Find NPC that player is looking at
	var interactingNPC int = -1

	for i, npc := range nm.NPCs {
		// Skip if NPC is already showing dialogue or is dead enemy
		if npc.ShowDialog && npc.DialogueTree != nil && npc.DialogueTree.IsActive {
			interactingNPC = i

			// Check for dialogue option selection (keys 1-9)
			for num := 0; num < 9; num++ {
				if keyState[sdl.SCANCODE_1+num] == 1 {
					nm.SelectDialogueOption(i, num)
					return
				}
			}

			// Check for ESC to exit dialogue
			if keyState[sdl.SCANCODE_ESCAPE] == 1 {
				nm.EndDialogue(i)
			}

			return
		}

		// Skip dead enemies
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
				interactingNPC = i
				break
			}
		}
	}

	// Interact with the found NPC
	if interactingNPC >= 0 && keyState[sdl.SCANCODE_E] == 1 {
		nm.StartDialogue(interactingNPC)
	}
}

// CheckDialogueInput checks for dialogue option inputs when a dialogue is active
func (nm *NPCManager) CheckDialogueInput(keyState []uint8) {
	// Find active dialogue
	for i, npc := range nm.NPCs {
		if npc.ShowDialog && npc.DialogueTree != nil && npc.DialogueTree.IsActive {
			// Check if grace period has passed
			if !npc.DialogueTree.ReadyToAdvance {
				// If grace period has elapsed, mark as ready
				if int64(DM.GlobalFrameCount)-npc.DialogueTree.GraceStartTime >= npc.DialogueTree.GracePeriod {
					npc.DialogueTree.ReadyToAdvance = true
				} else {
					// Still in grace period, don't accept input yet
					return
				}
			}

			// Check for Enter or Space to advance dialogue
			if keyState[sdl.SCANCODE_RETURN] == 1 || keyState[sdl.SCANCODE_SPACE] == 1 {
				// If this is the first advance after grace period, show options
				// For subsequent advances, this will act as "continue"
				npc.DialogueTree.ReadyToAdvance = true

				// Check for number keys only after player has had time to read
				currentNode := npc.DialogueTree.Nodes[npc.DialogueTree.CurrentNodeID]
				if currentNode != nil && len(currentNode.Options) == 1 {
					// If there's only one option, automatically select it as a convenience
					nm.SelectDialogueOption(i, 0)
					return
				}
			}

			// Only process option selection if ready to advance
			if npc.DialogueTree.ReadyToAdvance {
				// Check for dialogue option selection (keys 1-9)
				for num := 0; num < 9; num++ {
					if keyState[sdl.SCANCODE_1+num] == 1 {
						nm.SelectDialogueOption(i, num)
						return
					}
				}
			}

			// Check for ESC to exit dialogue (always available)
			if keyState[sdl.SCANCODE_ESCAPE] == 1 {
				nm.EndDialogue(i)
			}

			return // Only process one active dialogue at a time
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
