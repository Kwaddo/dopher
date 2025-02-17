# Dopher - A Doom-style Engine in Go

A simple raycasting engine inspired by the original Doom, implemented in Go using SDL2.

## Description

Dopher is a basic 3D rendering engine that uses raycasting techniques similar to those used in classic games like Wolfenstein 3D and Doom. It creates a pseudo-3D visualization from a 2D map.

## Features

- First-person perspective rendering
- Collision detection with walls
- Distance-based shading
- Smooth movement and rotation
- Configurable display settings
- Simple 2D map system

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
├── internal/
│   ├── constants/         # Game constants and configurations
│   │   ├── constant.go    # Screen dimensions and game settings
│   │   ├── maps.go        # Map definitions
│   │   └── player.go      # Player struct definition
│   └── funcs/            # Game logic
│       ├── cast.go       # Raycasting implementation
│       ├── checkers.go   # Collision detection
│       └── run.go        # Game loop and rendering
├── go.mod
└── go.sum
```

### Acknowledgements

1) Inspired by id Software's Doom
2) Built with go-sdl2
