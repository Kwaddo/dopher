package npc

import (
	Dialogue "doom/internal/character/dialogue"
	DM "doom/internal/model"
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

// StartDialogue initiates dialogue with an NPC.
func (nm *NPCManager) StartDialogue(npcIndex int) bool {
	if npcIndex < 0 || npcIndex >= len(nm.NPCs) {
		return false
	}
	activeDialogue := nm.GetActiveDialogueNPC()
	if activeDialogue >= 0 && activeDialogue != npcIndex {
		nm.EndDialogue(activeDialogue)
	}
	npc := nm.NPCs[npcIndex]
	if npc.DialogueTree != nil {
		npc.DialogueTree.IsActive = true
		npc.DialogueTree.ReadyToAdvance = false
		npc.DialogueTree.GraceStartTime = int64(DM.GlobalFrameCount)
		npc.DialogueTree.GracePeriod = 60
		npc.DialogueTree.CharsToShow = 0
		npc.DialogueTree.TextFullyShown = false
		npc.DialogueTree.LastCharTime = int64(DM.GlobalFrameCount)
		npc.DialogueTree.TextSpeed = 2
		if npc.DialogueTree.CurrentNodeID == "" {
			npc.DialogueTree.CurrentNodeID = "start"
		}
		currentNode := npc.DialogueTree.Nodes[npc.DialogueTree.CurrentNodeID]
		if currentNode == nil {
			return false
		}
		npc.DialogText = currentNode.Text
		npc.ShowDialog = true
		npc.DialogTimer = 600
		if currentNode.OnEnter != nil {
			currentNode.OnEnter(npc)
		}
		return true
	} else {
		npc.ShowDialog = true
		npc.DialogTimer = 180
		return true
	}
}

// EndDialogue ends the current dialogue.
func (nm *NPCManager) EndDialogue(npcIndex int) bool {
	if npcIndex < 0 || npcIndex >= len(nm.NPCs) {
		return false
	}
	npc := nm.NPCs[npcIndex]
	npc.ShowDialog = false
	if npc.DialogueTree != nil {
		npc.DialogueTree.IsActive = false
		npc.DialogueTree.TextFullyShown = false
		npc.DialogueTree.ReadyToAdvance = false
		npc.DialogueTree.CharsToShow = 0
		npc.DialogueTree.KeyWasDown = false
		npc.DialogueTree.CurrentNodeID = "start"
	}
	DM.InteractingNPC = -1
	return true
}

// UpdateTextAnimations updates the text animations for all active dialogues.
func (nm *NPCManager) UpdateTextAnimations() {
	keyState := sdl.GetKeyboardState()
	speedBoost := 1
	if keyState[sdl.SCANCODE_E] == 1 {
		speedBoost = 5
	}
	for _, npc := range nm.NPCs {
		if npc.ShowDialog && npc.DialogueTree != nil && npc.DialogueTree.IsActive {
			if !npc.DialogueTree.TextFullyShown {
				currentFrame := int64(DM.GlobalFrameCount)
				baseSpeed := npc.DialogueTree.TextSpeed
				effectiveSpeed := baseSpeed + speedBoost
				if effectiveSpeed > 8 {
					effectiveSpeed = 8
				}
				updateInterval := int64(4 - effectiveSpeed)
				if updateInterval < 1 {
					updateInterval = 1
				}
				if currentFrame > npc.DialogueTree.LastCharTime+updateInterval {
					charsToAdd := 1
					if speedBoost > 1 {
						charsToAdd = speedBoost
					}
					for i := 0; i < charsToAdd; i++ {
						if npc.DialogueTree.CharsToShow < len(npc.DialogText) {
							npc.DialogueTree.CharsToShow++
						}
					}
					npc.DialogueTree.LastCharTime = currentFrame
					if npc.DialogueTree.CharsToShow >= len(npc.DialogText) {
						npc.DialogueTree.CharsToShow = len(npc.DialogText)
						npc.DialogueTree.TextFullyShown = true
						npc.DialogueTree.GraceStartTime = currentFrame
					}
				}
			} else if !npc.DialogueTree.ReadyToAdvance {
				currentFrame := int64(DM.GlobalFrameCount)
				if currentFrame-npc.DialogueTree.GraceStartTime >= npc.DialogueTree.GracePeriod {
					npc.DialogueTree.ReadyToAdvance = true
				}
			}
		}
	}
}

// CreateBasicDialogueTree creates a dialogue tree by loading from a file or falling back to defaults.
func CreateBasicDialogueTree() *DM.DialogueTree {
	dialogueFile := "assets/dialogues/npc_basic.txt"
	tree, err := Dialogue.LoadDialogueFromFile(dialogueFile)
	if err == nil {
		return tree
	}
	fmt.Printf("Warning: Could not load dialogue file '%s': %v\nFalling back to default dialogue.\n",
		dialogueFile, err)
	tree = &DM.DialogueTree{
		Nodes:          make(map[string]*DM.DialogueNode),
		CurrentNodeID:  "start",
		IsActive:       false,
		ReadyToAdvance: false,
		GraceStartTime: 0,
		GracePeriod:    60,
	}
	tree.Nodes["start"] = &DM.DialogueNode{
		ID:   "start",
		Text: "TEXT FILE NOT FOUND.",
	}
	return tree
}

// AdvanceToNextDialogue moves to the next dialogue node in sequence.
func AdvanceToNextDialogue(npc *DM.NPC, nextNodeID string) {
	if nextNodeID == "" {
		npc.DialogueTree.IsActive = false
		npc.ShowDialog = false
		return
	}
	npc.DialogueTree.CurrentNodeID = nextNodeID
	npc.DialogueTree.CharsToShow = 0
	npc.DialogueTree.TextFullyShown = false
	npc.DialogueTree.LastCharTime = int64(DM.GlobalFrameCount)
	npc.DialogueTree.ReadyToAdvance = false
	npc.DialogueTree.GraceStartTime = int64(DM.GlobalFrameCount)
	nextNode := npc.DialogueTree.Nodes[nextNodeID]
	if nextNode == nil {
		npc.DialogueTree.IsActive = false
		npc.ShowDialog = false
		return
	}
	npc.DialogText = nextNode.Text
	npc.DialogTimer = 600
	if nextNode.OnEnter != nil {
		nextNode.OnEnter(npc)
	}
}

// GetActiveDialogueNPC returns the index of the first NPC with active dialogue or -1 if none.
func (nm *NPCManager) GetActiveDialogueNPC() int {
	for i, npc := range nm.NPCs {
		if npc.ShowDialog {
			return i
		}
	}
	return -1
}
