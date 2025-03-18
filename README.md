# Tetris Game

A classic Tetris game implementation in Go, created using Windsurf IDE.

## Features

- Classic Tetris gameplay with all 7 piece types (I, O, T, S, Z, J, L)
- Real-time piece rotation and movement
- Line clearing and scoring system
- Terminal-based graphics
- Smooth animations and controls

## How to Play

1. Run the game using one of these methods:
   - Build and run directly: `./build.sh`
   - Run in Docker: `docker build -t tetris . && docker run -it tetris`

2. Controls:
   - Left arrow: Move piece left
   - Right arrow: Move piece right
   - Down arrow: Move piece down
   - Up arrow: Rotate piece clockwise
   - Space: Drop piece instantly
   - Q: Quit game

3. Game Rules:
   - Stack pieces to create complete horizontal lines
   - Complete lines will be cleared, earning you points
   - The game ends when pieces stack to the top of the screen
   - Try to create as many complete lines as possible!

## Requirements

- Go 1.18 or later
- Terminal with ANSI escape code support

## Development

This project was created using Windsurf IDE, a powerful AI-assisted development environment. The game is built using Go's standard library and some external packages for terminal handling and keyboard input.

## License

This project is open source and available under the MIT License.
