# Dopher - A Doom-style Engine in Go

<p align="center">
  <img src="./assets/logo.png" alt="Dopher Logo"/>
</p>
<p align="center">
  <img src="./assets/preview.gif" />
</p>

A sophisticated raycasting engine inspired by the original Doom, implemented in Go using SDL2.

## Description

Dopher is an advanced 3D rendering engine that utilizes raycasting techniques similar to those used in classic games like Wolfenstein 3D and Doom. It creates an immersive pseudo-3D visualization from a 2D map, featuring textured walls, sprite-based NPCs, dynamic lighting effects, and realistic floor rendering with distance-based shading.

The engine employs Digital Differential Analysis (DDA) for precise wall detection and implements a z-buffer system for proper depth sorting of walls and sprites. This allows for realistic occlusion of NPCs behind walls and accurate perspective rendering. The engine features smooth player movement with momentum-based physics, dynamic field of view adjustments during sprinting and dashing, and interactive NPCs with a sophisticated dialogue tree system.

At its core, Dopher combines efficient raycasting algorithms with modern rendering techniques to create a seamless retro-inspired gaming experience. The physics system provides realistic movement with momentum, crouching capabilities, and collision detection, while the sprite system handles transparent NPCs that can be occluded by walls. The dynamic FOV system smoothly transitions during player movement, enhancing the sense of speed and immersion.

## Features

- First-person perspective rendering with raycasting
- Complete UI system with main menu, options, and pause menus
- Multiple maps support with seamless transitions
- Textured walls with proper perspective and distance shading
- Realistic floor rendering with gradient darkness based on distance
- Weapon system with muzzle flash effects
- Interactive dialogue trees with animated text display
- Frame-distributed enemy AI with health and combat system
- Smooth state transitions with fade effects
- Performance optimizations with double-buffered rendering
- Adaptive screen resolution and fullscreen handling
- Configurable game settings through options menu
- Texture caching for dialogue text rendering with LRU eviction
- Enemy pathfinding with obstacle avoidance and AI states
- Collision detection with walls and NPCs using circular hitboxes
- Distance-based shading and fog effects
- Dynamic FOV during sprinting and dashing with smooth LERP transitions
- Physics-based movement system with momentum and friction
- Crouching mechanics that adjust player height and speed
- Dash ability with cooldown and FOV effect
- Battle countdown system for enemy encounters
- Frame rate control system for consistent gameplay
- Z-buffer implementation for proper depth sorting
- Head bobbing effect during movement
- Multiple wall textures support
- Fullscreen toggle support (F key)
- Multi-threaded rendering with goroutines
- Transparent sprite rendering with proper depth testing
- Minimap toggle (TAB key)
- Fullscreen megamap toggle (M key)
- Real-time player position and direction indicators
- Continuous integration with GitHub Actions

## Controls

- W/A/S/D: Move around
- Left/Right Arrow: Rotate view
- Left Shift (hold): Sprint while moving
- Left Ctrl (hold): Crouch to move slower and reduce height
- Mouse: Aim weapon
- Left Mouse Button: Fire weapon
- Space: Dash in movement direction (with cooldown and FOV effect)
- E: Interact with NPCs and advance dialogue
- Enter: Exit dialogue immediately
- ESC: Toggle pause menu / Exit dialogue
- TAB: Toggle minimap
- M: Toggle fullscreen map
- F: Toggle fullscreen mode

## Game Menus

### Main Menu

- Start Game: Begin a new game
- Options: Configure game settings
- Quit: Exit the game

### Options Menu

- FOV: Adjust field of view (60°-110°)
- Mouse Sensitivity: Adjust rotation speed
- Movement Speed: Change player movement speed
- Head Bobbing: Toggle head bobbing effect
- Back: Return to main menu

### Pause Menu

- Resume: Return to gameplay with transition effect
- Return to Menu: Go back to main menu with transition effect
- Quit: Exit the game

## Prerequisites

- Go 1.24 or later
- SDL2 development libraries
- SDL2 TTF development libraries (for text rendering)

### Installing Dependencies

On Ubuntu/Debian:

```bash
sudo apt-get install libsdl2-dev libsdl2-ttf-dev
```

### Installation

A. Clone the repository:

```bash
git clone https://github.com/Kwaddo/dopher.git
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
│   ├── dialogues
│   │   └── npc_basic.txt
│   ├── font
│   │   └── dogicapixel.ttf
│   ├── preview
│   │   ├── logo.png
│   │   └── preview.gif
│   └── textures
│       ├── beef.bmp
│       ├── dictator.bmp
│       ├── gun.bmp
│       ├── muzzleflash.bmp
│       ├── npc.bmp
│       ├── wall2.bmp
│       └── wall.bmp
├── internal
│   ├── character
│   │   ├── dialogue
│   │   │   ├── loader.go
│   │   │   └── render.go
│   │   ├── npc
│   │   │   ├── checker.go
│   │   │   ├── dialogue.go
│   │   │   ├── enemy.go
│   │   │   └── npc.go
│   │   └── player
│   │       ├── battle.go
│   │       ├── checker.go
│   │       ├── gun.go
│   │       ├── inputs.go
│   │       ├── mobility.go
│   │       └── player.go
│   ├── core
│   │   ├── events.go
│   │   ├── initializer.go
│   │   └── loop.go
│   ├── graphics
│   │   ├── casting
│   │   │   ├── cast.go
│   │   │   ├── font.go
│   │   │   └── textures.go
│   │   └── renders
│   │       ├── general
│   │       │   ├── game.go
│   │       │   ├── npc.go
│   │       │   ├── scene.go
│   │       │   └── slices.go
│   │       ├── ui
│   │       │   ├── floor.go
│   │       │   ├── gun.go
│   │       │   ├── megamap.go
│   │       │   ├── minimap.go
│   │       │   └── roof.go
│   │       └── visual
│   │           ├── countdown.go
│   │           └── transition.go
│   ├── model
│   │   ├── constant.go
│   │   ├── maps.go
│   │   └── models.go
│   └── ui
│       ├── menu.go
│       ├── options.go
│       └── pause.go
├── go.mod
├── go.sum
├── main.go
└── README.md
```

### Technical Details

- Raycasting engine with DDA algorithm
- Multi-map support with map-specific NPCs
- Z-buffer implementation for proper depth sorting
- Double-buffered rendering for smooth frame rates
- Texture mapping with perspective correction
- Sprite system with transparency and occlusion
- Animated text dialogue system with variable speed
- Physics-based movement system with momentum
- Dialogue system with texture caching and LRU eviction
- Frame-distributed enemy AI for performance optimization
- Enemy health system with damage and death states
- Weapon system with fire rate and muzzle flash effects
- Dynamic FOV system with smooth LERP transitions for various movement states
- Multi-threaded rendering pipeline
- Distance-based floor rendering with gradient shading
- Collision detection using ray-circle intersection
- Head bobbing animation system synchronized with movement
- Configurable game settings through in-game menu
- Visual state transitions with fade effects
- Battle countdown system for enemy encounters
- Dynamic map system with multiple visualization modes
- Real-time player position tracking
- Directional indicators with triangle rendering
- Multi-layered map rendering with transparency
- Global font management system with size-based caching
- Centralized game state management
- Complete menu system with main, options, and pause states
- Crouching with height adjustment and speed reduction
- Dash mechanics with cooldown, FOV effect, and directional acceleration
- TTF font rendering and text display system
- Structured game initialization and resource cleanup
- Continuous integration with GitHub Actions

### Acknowledgements

1. Inspired by id Software's Doom
2. Built with go-sdl2
3. Font: Dogica Pixel by Roberto Mocci
