package dialogue

import (
	"bufio"
	DM "doom/internal/model"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LoadDialogueFromFile loads a dialogue tree from a text file
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

// LoadDialoguesFromDirectory loads all dialogue files from a directory
func LoadDialoguesFromDirectory(directory string) (map[string]*DM.DialogueTree, error) {
	dialogues := make(map[string]*DM.DialogueTree)

	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("failed to read dialogue directory '%s': %v", directory, err)
	}

	for _, file := range files {
		if file.IsDir() || (!strings.HasSuffix(file.Name(), ".txt") && !strings.HasSuffix(file.Name(), ".dialogue")) {
			continue
		}

		path := filepath.Join(directory, file.Name())
		dialogueName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))

		dialogue, err := LoadDialogueFromFile(path)
		if err != nil {
			return nil, fmt.Errorf("error loading dialogue '%s': %v", dialogueName, err)
		}

		dialogues[dialogueName] = dialogue
	}

	return dialogues, nil
}
