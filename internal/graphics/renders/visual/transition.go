package renders

import (
	DM "doom/internal/models/global"

	"github.com/veandco/go-sdl2/sdl"
)

// RenderTransition renders the current transition effect (fade to black and back)
func RenderTransition(renderer *sdl.Renderer) {
	if DM.TransitionState == DM.TransitionInactive {
		return
	}
	var alpha uint8
	if DM.TransitionState == DM.TransitionFadeOut {
		alpha = uint8(255.0 * (DM.TransitionProgress / DM.TransitionDuration))
	} else {
		alpha = uint8(255.0 * (1.0 - DM.TransitionProgress/DM.TransitionDuration))
	}
	renderer.SetDrawColor(0, 0, 0, alpha)
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	renderer.FillRect(&sdl.Rect{
		X: 0,
		Y: 0,
		W: int32(DM.ScreenWidth),
		H: int32(DM.ScreenHeight),
	})
}

// UpdateTransition updates the transition state and executes the callback when needed
func UpdateTransition() bool {
	if DM.TransitionState == DM.TransitionInactive {
		return false
	}
	DM.TransitionProgress++
	if DM.TransitionProgress >= DM.TransitionDuration {
		if DM.TransitionState == DM.TransitionFadeOut {
			if DM.TransitionCallback != nil {
				DM.TransitionCallback()
			}
			DM.TransitionState = DM.TransitionFadeIn
			DM.TransitionProgress = 0
		} else {
			DM.TransitionState = DM.TransitionInactive
			DM.TransitionProgress = 0
			return false
		}
	}
	return true
}

// StartTransition initiates a fade-out/fade-in transition with the specified callback.
func StartTransition(callback func()) {
	if DM.TransitionState == DM.TransitionInactive {
		DM.TransitionState = DM.TransitionFadeOut
		DM.TransitionProgress = 0
		DM.TransitionCallback = callback
	}
}
