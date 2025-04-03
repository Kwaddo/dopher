package core

import (
	DM "doom/internal/global"
	Visual "doom/internal/graphics/renders/visual"
	Menu "doom/internal/ui"

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
				if DM.GlobalGameState.InMainMenu {
					switch t.Keysym.Sym {
					case sdl.K_UP:
						Menu.GlobalMainMenu.MoveUp()
					case sdl.K_DOWN:
						Menu.GlobalMainMenu.MoveDown()
					case sdl.K_RETURN, sdl.K_SPACE:
						switch Menu.GlobalMainMenu.GetSelectedOption() {
						case "Start Game":
							DM.GlobalGameState.InMainMenu = false
						case "Options":
							DM.GlobalGameState.InOptionsMenu = true
							DM.GlobalGameState.InMainMenu = false
						case "Quit":
							return true
						}
					}
					continue
				}
				if DM.GlobalGameState.InOptionsMenu {
					switch t.Keysym.Sym {
					case sdl.K_UP:
						Menu.GlobalOptionsMenu.MoveUp()
					case sdl.K_DOWN:
						Menu.GlobalOptionsMenu.MoveDown()
					case sdl.K_LEFT:
						Menu.GlobalOptionsMenu.DecreaseSetting()
					case sdl.K_RIGHT:
						Menu.GlobalOptionsMenu.IncreaseSetting()
					case sdl.K_RETURN, sdl.K_SPACE:
						if Menu.GlobalOptionsMenu.GetSelectedOption() == "Back" {
							DM.GlobalGameState.InOptionsMenu = false
							DM.GlobalGameState.InMainMenu = true
						}
					case sdl.K_ESCAPE:
						DM.GlobalGameState.InOptionsMenu = false
						DM.GlobalGameState.InMainMenu = true
					}
					continue
				}
				if DM.GlobalGameState.IsPaused {
					switch t.Keysym.Sym {
					case sdl.K_UP:
						Menu.GlobalPauseMenu.MoveUp()
					case sdl.K_DOWN:
						Menu.GlobalPauseMenu.MoveDown()
					case sdl.K_RETURN, sdl.K_SPACE:
						switch Menu.GlobalPauseMenu.GetSelectedOption() {
						case "Resume":
							Visual.StartTransition(func() {
								DM.GlobalGameState.IsPaused = false
							})
							return false
						case "Return to Menu":
							Visual.StartTransition(func() {
								DM.GlobalGameState.IsPaused = false
								DM.GlobalGameState.InMainMenu = true
							})
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
					} else {
						window.SetFullscreen(sdl.WINDOW_FULLSCREEN_DESKTOP)
						w, h, err := renderer.GetOutputSize()
						if err != nil {
							panic(err)
						}
						DM.ScreenWidth = float64(w)
						DM.ScreenHeight = float64(h)
					}
					renderer.SetViewport(&sdl.Rect{X: 0, Y: 0, W: int32(DM.ScreenWidth), H: int32(DM.ScreenHeight)})
					DM.ZBuffer = make([]float64, int(DM.ScreenWidth))
					DM.NeedToRecreateBuffers = true
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
