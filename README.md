# Dopher - A Doom-style Engine in Go

<p align="center">
  <img src="./assets/logo.png" alt="Dopher Logo"/>
</p>
<p align="center">
  <img src="./assets/preview.gif" />
</p>

A sophisticated raycasting engine inspired by the original Doom, implemented in Go using SDL2.

## Description

Dopher is an advanced 3D rendering engine that utilizes raycasting techniques similar to those used in classic games like Wolfenstein 3D and Doom. It creates an immersive pseudo-3D visualization from a 2D map, featuring textured walls, sprite-based NPCs, and dynamic lighting effects.

The engine employs Digital Differential Analysis (DDA) for precise wall detection and implements a z-buffer system for proper depth sorting of walls and sprites. This allows for realistic occlusion of NPCs behind walls and accurate perspective rendering. The engine features smooth player movement with momentum-based physics, dynamic field of view adjustments during sprinting, and interactive NPCs with a typewriter-style dialog system.

At its core, Dopher combines efficient raycasting algorithms with modern rendering techniques to create a seamless retro-inspired gaming experience. The physics system provides realistic movement with momentum and collision detection, while the sprite system handles transparent NPCs that can be occluded by walls. The dynamic FOV system smoothly transitions during player movement, enhancing the sense of speed and immersion.

## Features

- First-person perspective rendering with raycasting
- Textured walls with proper perspective and distance shading
- Sprite-based NPCs with transparency and occlusion
- Interactive NPCs with typewriter-style dialog system
- Collision detection with walls and NPCs using circular hitboxes
- Distance-based shading and fog effects
- Dynamic FOV during sprinting with smooth transitions
- Physics-based movement system with momentum and friction
- Advanced texture mapping with proper coordinate calculation
- Configurable display settings and game parameters
- Simple 2D map system with extensible layout
- Z-buffer implementation for proper depth sorting
- Head bobbing effect during movement
- Smooth floor rendering with distance-based shading
- Multiple wall textures support
- Fullscreen toggle support (F key)
- Optimized rendering with goroutines
- Transparent sprite rendering with proper depth testing

## Controls

- W/A/S/D: Move around
- Left/Right Arrow: Rotate view
- Left Shift (hold): Sprint while moving
- E: Interact with NPCs
- ESC/Q: Quit game

## Prerequisites

- Go 1.24 or later
- SDL2 development libraries
- SDL2 TTF development libraries (for dialog system)

### Installing Dependencies

On Ubuntu/Debian:

```bash
sudo apt-get install libsdl2-dev libsdl2-ttf-dev
```

### Installation

A. Clone the repository:

```bash
git clone https://github.com/YourUsername/dopher.git
cd dopher
```

B. Install dependencies:

```bash
go mod tidy
```

C. Run the game:

```bash
go run main.go
```

### Project Structure

```struct
dopher/
├── assets
│   ├── beef.bmp
│   ├── dogicapixel.ttf
│   ├── logo.png
│   ├── npc.bmp
│   ├── preview.gif
│   ├── wall2.bmp
│   └── wall.bmp
├── internal
│   ├── char
│   │   ├── npc
│   │   │   ├── checker.go
│   │   │   ├── dialog.go
│   │   │   └── npc.go
│   │   └── player
│   │       ├── checker.go
│   │       └── movement.go
│   ├── core
│   │   └── run.go
│   ├── graphics
│   │   ├── casting
│   │   │   ├── cast.go
│   │   │   └── textures.go
│   │   └── renders
│   │       ├── floor.go
│   │       ├── npc.go
│   │       ├── scene.go
│   │       └── slices.go
│   └── model
│       ├── constant.go
│       ├── maps.go
│       └── models.go
├── go.mod
├── go.sum
├── main.go
└── README.md
```

### Technical Details

- Raycasting engine with DDA algorithm
- Z-buffer implementation for proper depth sorting
- Texture mapping with perspective correction
- Sprite system with transparency and occlusion
- Physics-based movement system with momentum
- Dialog system with typewriter effect
- Dynamic FOV system with smooth transitions
- Multi-threaded rendering pipeline
- Distance-based fog and shading system
- Collision detection using ray-circle intersection
- Head bobbing animation system
- Advanced floor rendering with gradient shading
- Smooth state transitions using LERP

### Acknowledgements

1. Inspired by id Software's Doom
2. Built with go-sdl2
3. Font: Dogica Pixel by Roberto Mocci
