package core

import (
	DM "doom/internal/model"
	Graphics "doom/internal/graphics/renders"
	"github.com/veandco/go-sdl2/sdl"
)

// HandleEvents handles certain events for the game.
func HandleEvents(window *sdl.Window, renderer *sdl.Renderer) bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			return true
		case *sdl.KeyboardEvent:
			if t.State == sdl.PRESSED {
				if DM.GlobalGameState.IsPaused {
					switch t.Keysym.Sym {
					case sdl.K_UP:
						Graphics.GlobalPauseMenu.MoveUp()
					case sdl.K_DOWN:
						Graphics.GlobalPauseMenu.MoveDown()
					case sdl.K_RETURN, sdl.K_SPACE:
						switch Graphics.GlobalPauseMenu.GetSelectedOption() {
						case "Resume":
							DM.GlobalGameState.IsPaused = false
							return false
						case "Quit":
							return true
						}
					case sdl.K_ESCAPE:
						DM.GlobalGameState.IsPaused = false
					}
					continue
				}
				switch t.Keysym.Sym {
				case sdl.K_ESCAPE:
					DM.GlobalGameState.IsPaused = !DM.GlobalGameState.IsPaused
				case sdl.K_f:
					flags := window.GetFlags()
					if flags&sdl.WINDOW_FULLSCREEN_DESKTOP == sdl.WINDOW_FULLSCREEN_DESKTOP {
						window.SetFullscreen(0)
						DM.ScreenWidth = 1500
						DM.ScreenHeight = 900
						renderer.SetViewport(&sdl.Rect{X: 0, Y: 0, W: int32(DM.ScreenWidth), H: int32(DM.ScreenHeight)})
						DM.ZBuffer = make([]float64, int(DM.ScreenWidth))
					} else {
						window.SetFullscreen(sdl.WINDOW_FULLSCREEN_DESKTOP)
						screenSurface, err := window.GetSurface()
						if err != nil {
							panic(err)
						}
						DM.ScreenWidth = float64(screenSurface.W)
						DM.ScreenHeight = float64(screenSurface.H)
						renderer.SetViewport(&sdl.Rect{X: 0, Y: 0, W: int32(DM.ScreenWidth), H: int32(DM.ScreenHeight)})
						DM.ZBuffer = make([]float64, int(DM.ScreenWidth))
					}
				case sdl.K_TAB:
					DM.ShowMiniMap = !DM.ShowMiniMap
					DM.ShowMegaMap = false
				case sdl.K_m:
					DM.ShowMegaMap = !DM.ShowMegaMap
					DM.ShowMiniMap = false
				}
			}
		}
	}
	return false
}
