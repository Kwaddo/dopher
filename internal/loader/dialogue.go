package loader

import (
	"bufio"
	DM "doom/internal/models/global"
	"fmt"
	"os"
	"strings"
)

// CreateBasicDialogueTree creates a dialogue tree by loading from a file or falling back to defaults.
func CreateBasicDialogueTree() *DM.DialogueTree {
	dialogueFile := "assets/dialogues/npc_basic.txt"
	tree, err := LoadDialogueFromFile(dialogueFile)
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

// LoadDialogueFromFile loads a dialogue tree from a text file.
func LoadDialogueFromFile(filename string) (*DM.DialogueTree, error) {
	tree := &DM.DialogueTree{
		Nodes:          make(map[string]*DM.DialogueNode),
		CurrentNodeID:  "start",
		IsActive:       false,
		ReadyToAdvance: false,
		GraceStartTime: 0,
		GracePeriod:    60,
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open dialogue file '%s': %v", filename, err)
	}
	defer file.Close()
	var currentNode *DM.DialogueNode
	var currentNodeID string
	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			if currentNode != nil {
				tree.Nodes[currentNodeID] = currentNode
			}
			currentNodeID = line[1 : len(line)-1]
			currentNode = &DM.DialogueNode{
				ID: currentNodeID,
			}
			continue
		}
		if currentNode == nil {
			return nil, fmt.Errorf("line %d: content found before node definition", lineNum)
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("line %d: invalid format, expected 'Key: Value'", lineNum)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch strings.ToLower(key) {
		case "text":
			currentNode.Text = value
		case "nextid":
			currentNode.NextID = value
		default:
			return nil, fmt.Errorf("line %d: unknown key '%s'", lineNum, key)
		}
	}
	if currentNode != nil {
		tree.Nodes[currentNodeID] = currentNode
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading dialogue file: %v", err)
	}
	if len(tree.Nodes) == 0 {
		return nil, fmt.Errorf("no dialogue nodes found in file")
	}
	if _, exists := tree.Nodes["start"]; !exists {
		return nil, fmt.Errorf("dialogue file must contain a 'start' node")
	}
	return tree, nil
}
