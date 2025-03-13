package renders

import (
	Casts "doom/internal/graphics/casting"
	DM "doom/internal/model"
	"fmt"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

// RenderCountdown renders the battle countdown with enhanced visual effects
func RenderCountdown(renderer *sdl.Renderer) {
	if DM.CountdownState == DM.CountdownInactive {
		return
	}
	numberToShow := DM.CountdownNumber - int(DM.CountdownProgress/60)
	frameInCurrentNumber := int(DM.CountdownProgress) % 60
	if DM.CountdownProgress >= 180 {
		RenderFightText(renderer, frameInCurrentNumber)
		return
	}
	scale := 1.0
	if frameInCurrentNumber < 15 {
		scale = 1.3 - float64(frameInCurrentNumber)/50
	} else if frameInCurrentNumber > 45 {
		scale = 1.0 - float64(frameInCurrentNumber-45)/75
	}
	var alpha uint8
	if frameInCurrentNumber < 15 {
		progress := float64(frameInCurrentNumber) / 15.0
		alpha = uint8(255 * progress * progress)
	} else if frameInCurrentNumber > 45 {
		progress := float64(60-frameInCurrentNumber) / 15.0
		alpha = uint8(255 * progress * progress)
	} else {
		alpha = 255
	}
	RenderCountdownNumber(renderer, numberToShow, alpha, scale, frameInCurrentNumber)
}

// RenderCountdownNumber renders a single countdown number with glow effect.
func RenderCountdownNumber(renderer *sdl.Renderer, number int, alpha uint8, scale float64, frameInNumber int) {
	countdownFont, err := Casts.GlobalFontManager.GetFont(120)
	if err != nil {
		return
	}
	numberText := fmt.Sprintf("%d", number)
	var textColor sdl.Color
	switch number {
	case 3:
		textColor = sdl.Color{R: 255, G: 50, B: 50, A: alpha}
	case 2:
		textColor = sdl.Color{R: 255, G: 200, B: 50, A: alpha}
	case 1:
		textColor = sdl.Color{R: 50, G: 255, B: 50, A: alpha}
	default:
		textColor = sdl.Color{R: 255, G: 255, B: 255, A: alpha}
	}
	surface, err := countdownFont.RenderUTF8Solid(numberText, textColor)
	if err != nil {
		return
	}
	defer surface.Free()
	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return
	}
	defer texture.Destroy()
	pulseOffset := math.Sin(float64(frameInNumber)/10) * 5
	scaledWidth := int32(float64(surface.W) * scale)
	scaledHeight := int32(float64(surface.H) * scale)

	bounceOffset := 0.0
	if frameInNumber > 5 && frameInNumber < 35 {
		bounceAmplitude := 30 * (1.0 - float64(frameInNumber-5)/30.0)
		bounceOffset = math.Sin(float64(frameInNumber-5)/5) * bounceAmplitude
	}
	texture.SetBlendMode(sdl.BLENDMODE_BLEND)
	texture.SetAlphaMod(alpha)
	centerX := DM.ScreenWidth/2 - float64(scaledWidth)/2
	centerY := DM.ScreenHeight/2 - float64(scaledHeight)/2 - bounceOffset - pulseOffset
	renderer.Copy(texture, nil, &sdl.Rect{
		X: int32(centerX),
		Y: int32(centerY),
		W: scaledWidth,
		H: scaledHeight,
	})
}

// RenderFightText renders the "FIGHT!" text with dramatic animation.
func RenderFightText(renderer *sdl.Renderer, frameInNumber int) {
	fightFont, err := Casts.GlobalFontManager.GetFont(120)
	if err != nil {
		return
	}
	fightProgress := float64(frameInNumber-40) / 20.0
	if fightProgress > 1.0 {
		fightProgress = 1.0
	}
	r := uint8(255)
	g := uint8(50 + 200*math.Sin(fightProgress*math.Pi))
	b := uint8(50)
	alpha := uint8(255 * fightProgress)
	fightSurface, err := fightFont.RenderUTF8Blended("FIGHT!", sdl.Color{R: r, G: g, B: b, A: alpha})
	if err != nil {
		return
	}
	defer fightSurface.Free()
	fightTexture, err := renderer.CreateTextureFromSurface(fightSurface)
	if err != nil {
		return
	}
	defer fightTexture.Destroy()
	scale := 1.0 + math.Sin(fightProgress*math.Pi/2)*0.3
	scaledWidth := int32(float64(fightSurface.W) * scale)
	scaledHeight := int32(float64(fightSurface.H) * scale)
	shakeAmount := (1.0 - fightProgress) * 4.0
	shakeX := int32(math.Sin(float64(frameInNumber)*0.8) * shakeAmount)
	shakeY := int32(math.Cos(float64(frameInNumber)*1.2) * shakeAmount)
	fightTexture.SetBlendMode(sdl.BLENDMODE_BLEND)
	centerX := DM.ScreenWidth/2 - float64(scaledWidth)/2
	centerY := DM.ScreenHeight/2 - float64(scaledHeight)/2
	posX := int32(centerX) + shakeX
	posY := int32(centerY) + shakeY
	fightTexture.SetAlphaMod(alpha)
	renderer.Copy(fightTexture, nil, &sdl.Rect{
		X: posX,
		Y: posY,
		W: scaledWidth,
		H: scaledHeight,
	})
}

// UpdateCountdown updates the countdown state.
func UpdateCountdown() {
	if DM.CountdownState == DM.CountdownInactive {
		return
	}
	DM.CountdownProgress++
	if DM.CountdownProgress >= DM.CountdownDuration+120 {
		DM.CountdownState = DM.CountdownInactive
		DM.CountdownProgress = 0
		DM.CountdownFreeze = false
	} else if DM.CountdownProgress >= DM.CountdownDuration {
		DM.CountdownFreeze = false
	}
}

// StartCountdown initiates the battle countdown.
func StartCountdown() {
	DM.CountdownState = DM.CountdownActive
	DM.CountdownProgress = 0
	DM.CountdownFreeze = true
}
