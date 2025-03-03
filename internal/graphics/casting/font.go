package casting

import (
	DM "doom/internal/model"
	"fmt"

	"github.com/veandco/go-sdl2/ttf"
)

type FontManager DM.FontManager

var GlobalFontManager = &FontManager{
	Path:          "assets/font/dogicapixel.ttf",
	Cache:         make(map[int]*ttf.Font),
	IsInitialized: false,
}

// InitFonts initializes the TTF system and prepares the font manager
func InitFonts() error {
	if GlobalFontManager.IsInitialized {
		return nil
	}
	if err := ttf.Init(); err != nil {
		return fmt.Errorf("failed to initialize TTF: %v", err)
	}
	sizes := []int{16, 24, 36}
	for _, size := range sizes {
		if _, err := GlobalFontManager.GetFont(size); err != nil {
			return err
		}
	}
	GlobalFontManager.IsInitialized = true
	return nil
}

// GetFont returns a font at the specified size, loading it if necessary
func (fm *FontManager) GetFont(size int) (*ttf.Font, error) {
	fm.Mutex.RLock()
	font, exists := fm.Cache[size]
	fm.Mutex.RUnlock()
	if exists {
		return font, nil
	}
	fm.Mutex.Lock()
	defer fm.Mutex.Unlock()
	if font, exists := fm.Cache[size]; exists {
		return font, nil
	}

	font, err := ttf.OpenFont(fm.Path, size)
	if err != nil {
		return nil, fmt.Errorf("failed to load font at size %d: %v", size, err)
	}
	fm.Cache[size] = font
	return font, nil
}

// CleanupFonts closes all loaded fonts and quits the TTF subsystem
func CleanupFonts() {
	GlobalFontManager.Mutex.Lock()
	defer GlobalFontManager.Mutex.Unlock()

	for _, font := range GlobalFontManager.Cache {
		font.Close()
	}

	GlobalFontManager.Cache = make(map[int]*ttf.Font)
	GlobalFontManager.IsInitialized = false
	ttf.Quit()
}
