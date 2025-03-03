package model

// MoveUp moves the selection up
func (pm *PauseMenu) MoveUp() {
	pm.CurrentOption--
	if pm.CurrentOption < 0 {
		pm.CurrentOption = len(pm.Options) - 1
	}
}

// MoveDown moves the selection down
func (pm *PauseMenu) MoveDown() {
	pm.CurrentOption = (pm.CurrentOption + 1) % len(pm.Options)
}

// GetSelectedOption returns the currently selected option
func (pm *PauseMenu) GetSelectedOption() string {
	return pm.Options[pm.CurrentOption]
}
