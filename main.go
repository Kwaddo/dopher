package main

import (
	"fmt"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	screenWidth     = 800
	screenHeight    = 600
	fov             = math.Pi / 3
	numRays         = 120
	maxDepth        = 800.0
	moveSpeed       = 5.0
	rotateSpeed     = 0.1
	collisionBuffer = 20.0 // Buffer distance from walls
)

// Add map definition
var worldMap = [][]int{
	{1, 1, 1, 1, 1, 1, 1, 1},
	{1, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 1, 0, 0, 1, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 1, 0, 0, 1, 0, 1},
	{1, 0, 1, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 1},
	{1, 1, 1, 1, 1, 1, 1, 1},
}

type Player struct {
	x, y  float64
	angle float64
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Println("Failed to initialize SDL:", err)
		return
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Doom in Go", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, screenWidth, screenHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println("Could not create window:", err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println("Could not create renderer:", err)
		return
	}
	defer renderer.Destroy()

	// Initialize player in the middle of the screen
	player := Player{
		x:     150, // Position player in first room
		y:     150,
		angle: 0,
	}

	// Main loop
	running := true
	for running {
		// Handle events
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}

		// Handle keyboard input
		keys := sdl.GetKeyboardState()

		// Rotate with left/right arrows
		if keys[sdl.SCANCODE_LEFT] == 1 {
			player.angle -= rotateSpeed
		}
		if keys[sdl.SCANCODE_RIGHT] == 1 {
			player.angle += rotateSpeed
		}

		// Move with up/down arrows
		if keys[sdl.SCANCODE_UP] == 1 {
			newX := player.x + math.Cos(player.angle)*moveSpeed
			newY := player.y + math.Sin(player.angle)*moveSpeed
			if !checkCollision(newX, newY) {
				player.x = newX
				player.y = newY
			}
		}
		if keys[sdl.SCANCODE_DOWN] == 1 {
			newX := player.x - math.Cos(player.angle)*moveSpeed
			newY := player.y - math.Sin(player.angle)*moveSpeed
			if !checkCollision(newX, newY) {
				player.x = newX
				player.y = newY
			}
		}

		// Render Scene
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		// Draw 3D view
		rayAngle := player.angle - fov/2
		for i := 0; i < numRays; i++ {
			distance := castRay(player.x, player.y, rayAngle)

			// Fix fisheye effect
			distance = distance * math.Cos(rayAngle-player.angle)

			// Calculate wall height
			wallHeight := (screenHeight / distance) * 50
			if wallHeight > screenHeight {
				wallHeight = screenHeight
			}

			// Draw wall slice
			wallTop := (screenHeight - wallHeight) / 2
			wallRect := sdl.Rect{
				X: int32(i * (screenWidth / numRays)),
				Y: int32(wallTop),
				W: int32(screenWidth/numRays + 1),
				H: int32(wallHeight),
			}

			// Color based on distance
			intensity := uint8(255 - math.Min(255, distance/2))
			renderer.SetDrawColor(intensity, intensity/2, intensity/2, 255)
			renderer.FillRect(&wallRect)

			rayAngle += fov / float64(numRays)
		}

		renderer.Present()
		sdl.Delay(16) // ~60 FPS
	}
}

func castRay(startX, startY, angle float64) float64 {
	rayX := math.Cos(angle)
	rayY := math.Sin(angle)

	distance := 0.0
	for distance < maxDepth {
		x := startX + rayX*distance
		y := startY + rayY*distance

		mapX := int(x / 100)
		mapY := int(y / 100)

		if mapX >= 0 && mapX < len(worldMap[0]) && mapY >= 0 && mapY < len(worldMap) {
			if worldMap[mapY][mapX] == 1 {
				return distance
			}
		}

		distance += 1.0
	}
	return maxDepth
}

func checkCollision(x, y float64) bool {
	// Check multiple points around the player
	for angle := 0.0; angle < 2*math.Pi; angle += math.Pi / 4 {
		checkX := x + math.Cos(angle)*collisionBuffer
		checkY := y + math.Sin(angle)*collisionBuffer

		mapX := int(checkX / 100)
		mapY := int(checkY / 100)

		if mapX >= 0 && mapX < len(worldMap[0]) && mapY >= 0 && mapY < len(worldMap) {
			if worldMap[mapY][mapX] == 1 {
				return true
			}
		}
	}
	return false
}
