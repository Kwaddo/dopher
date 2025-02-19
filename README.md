# Dopher - A Doom-style Engine in Go

![alt text](./assets/logo.png)

A simple raycasting engine inspired by the original Doom, implemented in Go using SDL2.

## Description

Dopher is a basic 3D rendering engine that uses raycasting techniques similar to those used in classic games like Wolfenstein 3D and Doom. It creates a pseudo-3D visualization from a 2D map with textured walls and sprite-based NPCs.

## Features

- First-person perspective rendering
- Textured walls with proper perspective
- Sprite-based NPCs with transparency
- Collision detection with walls and NPCs
- Distance-based shading and fog
- Dynamic FOV during sprinting
- Smooth movement and rotation
- Configurable display settings
- Simple 2D map system

## Controls

- W/A/S/D: Move around
- Left/Right Arrow: Rotate view
- Left Shift: Sprint
- ESC/Q: Quit game

## Prerequisites

- Go 1.24 or later
- SDL2 development libraries

### Installing SDL2

On Ubuntu/Debian:

```bash
sudo apt-get install libsdl2-dev
```

### Installation

- Clone the repository:

```bash
git clone https://github.com/Kwaddo/dopher.git
cd dopher
```

- Install dependencies:

```bash
go mod tidy
```

- Run the game! Use arrow keys to move around.

```bash
go run main.go
```

### Project Structure

```struct
dopher/
├── main.go                 # Main entry point
├── assets/                # Game assets
│   ├── wall.bmp          # Wall texture
│   └── npc.bmp           # NPC sprite
├── internal/
│   ├── constants/        # Game constants and configurations
│   │   ├── constant.go   # Screen dimensions and game settings
│   │   ├── maps.go       # Map definitions
│   │   └── models.go     # Struct definitions
│   └── core/            # Game logic
│       ├── cast.go      # Raycasting implementation
│       ├── checkers.go  # Collision detection
│       ├── npc.go       # NPC management
│       ├── player.go    # Player controls
│       ├── render.go    # Rendering engine
│       ├── run.go       # Game loop
│       └── textures.go  # Texture management
├── go.mod
└── go.sum
```

### Acknowledgements

1) Inspired by id Software's Doom
2) Built with go-sdl2
