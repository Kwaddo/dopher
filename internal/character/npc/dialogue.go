package npc

import (
	DM "doom/internal/model"
)

// StartDialogue initiates dialogue with an NPC
func (nm *NPCManager) StartDialogue(npcIndex int) bool {
	if npcIndex < 0 || npcIndex >= len(nm.NPCs) {
		return false
	}

	npc := nm.NPCs[npcIndex]

	// If NPC has a dialogue tree
	if npc.DialogueTree != nil {
		npc.DialogueTree.IsActive = true

		// Reset grace period
		npc.DialogueTree.ReadyToAdvance = false
		npc.DialogueTree.GraceStartTime = int64(DM.GlobalFrameCount)
		npc.DialogueTree.GracePeriod = 60 // About 1 second at 60 FPS

		// Start with the "start" node by default
		if npc.DialogueTree.CurrentNodeID == "" {
			npc.DialogueTree.CurrentNodeID = "start"
		}

		// Get current node
		currentNode := npc.DialogueTree.Nodes[npc.DialogueTree.CurrentNodeID]
		if currentNode == nil {
			return false
		}

		// Set dialogue text and show it
		npc.DialogText = currentNode.Text
		npc.ShowDialog = true
		npc.DialogTimer = 600 // Longer timer for interactive dialogue

		// Call OnEnter callback if it exists
		if currentNode.OnEnter != nil {
			currentNode.OnEnter(npc)
		}

		return true
	} else {
		// Fallback for NPCs without dialogue trees - show simple dialogue
		npc.ShowDialog = true
		npc.DialogTimer = 180
		return true
	}
}

// SelectDialogueOption selects a dialogue option and advances the conversation
func (nm *NPCManager) SelectDialogueOption(npcIndex int, optionIndex int) bool {
	if npcIndex < 0 || npcIndex >= len(nm.NPCs) {
		return false
	}

	npc := nm.NPCs[npcIndex]

	if npc.DialogueTree == nil || !npc.DialogueTree.IsActive {
		return false
	}

	currentNode := npc.DialogueTree.Nodes[npc.DialogueTree.CurrentNodeID]
	if currentNode == nil || optionIndex < 0 || optionIndex >= len(currentNode.Options) {
		return false
	}

	// Get the selected option
	option := currentNode.Options[optionIndex]

	// Check if the option has a condition and if it's met
	if option.Condition != nil && !option.Condition(npc) {
		return false
	}

	// If next node is empty, end dialogue
	if option.NextNode == "" {
		npc.DialogueTree.IsActive = false
		npc.ShowDialog = false
		return true
	}

	// Move to the next node
	npc.DialogueTree.CurrentNodeID = option.NextNode

	// Reset grace period for new text
	npc.DialogueTree.ReadyToAdvance = false
	npc.DialogueTree.GraceStartTime = int64(DM.GlobalFrameCount)

	// Update the dialogue text
	nextNode := npc.DialogueTree.Nodes[option.NextNode]
	if nextNode == nil {
		// Node not found, end dialogue
		npc.DialogueTree.IsActive = false
		npc.ShowDialog = false
		return false
	}

	npc.DialogText = nextNode.Text
	npc.DialogTimer = 600 // Reset timer

	// Call OnEnter callback if it exists
	if nextNode.OnEnter != nil {
		nextNode.OnEnter(npc)
	}

	return true
}

// EndDialogue ends the current dialogue
func (nm *NPCManager) EndDialogue(npcIndex int) bool {
	if npcIndex < 0 || npcIndex >= len(nm.NPCs) {
		return false
	}

	npc := nm.NPCs[npcIndex]

	if npc.DialogueTree != nil {
		npc.DialogueTree.IsActive = false
	}

	npc.ShowDialog = false
	return true
}

// CreateBasicDialogueTree creates a simple dialogue tree for testing
func CreateBasicDialogueTree() *DM.DialogueTree {
	tree := &DM.DialogueTree{
		Nodes:          make(map[string]*DM.DialogueNode),
		CurrentNodeID:  "start",
		IsActive:       false,
		ReadyToAdvance: false,
		GraceStartTime: 0,
		GracePeriod:    60, // About 1 second at 60 FPS
	}

	// Add start node
	tree.Nodes["start"] = &DM.DialogueNode{
		ID:   "start",
		Text: "Hello there! What can I help you with?",
		Options: []DM.DialogueOption{
			{Text: "Tell me about yourself", NextNode: "about"},
			{Text: "I need information", NextNode: "info"},
			{Text: "Goodbye", NextNode: ""},
		},
	}

	// Add about node
	tree.Nodes["about"] = &DM.DialogueNode{
		ID:   "about",
		Text: "I'm just an NPC living in this world. Not much to tell really.",
		Options: []DM.DialogueOption{
			{Text: "Tell me more", NextNode: "about_more"},
			{Text: "Let's talk about something else", NextNode: "start"},
		},
	}

	// Add more about info
	tree.Nodes["about_more"] = &DM.DialogueNode{
		ID:   "about_more",
		Text: "I was created for this game. That's all there is to know!",
		Options: []DM.DialogueOption{
			{Text: "Back to main topics", NextNode: "start"},
			{Text: "End conversation", NextNode: ""},
		},
	}

	// Add info node
	tree.Nodes["info"] = &DM.DialogueNode{
		ID:   "info",
		Text: "What kind of information are you looking for?",
		Options: []DM.DialogueOption{
			{Text: "About this area", NextNode: "area_info"},
			{Text: "About enemies", NextNode: "enemy_info"},
			{Text: "Back", NextNode: "start"},
		},
	}

	// Add area info
	tree.Nodes["area_info"] = &DM.DialogueNode{
		ID:   "area_info",
		Text: "You're in a dangerous place. Be careful as you explore.",
		Options: []DM.DialogueOption{
			{Text: "Thanks for the warning", NextNode: "start"},
		},
	}

	// Add enemy info
	tree.Nodes["enemy_info"] = &DM.DialogueNode{
		ID:   "enemy_info",
		Text: "The enemies here are fierce. Shoot them before they get to you!",
		Options: []DM.DialogueOption{
			{Text: "I'll be careful", NextNode: "start"},
		},
	}

	return tree
}
