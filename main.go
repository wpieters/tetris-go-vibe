package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"golang.org/x/term"

	"github.com/eiannone/keyboard"
)

const (
	boardWidth  = 10
	boardHeight = 20
	emptyCell   = " "
	frameRate   = 500 * time.Millisecond
)

// Piece types and their display characters
const (
	PieceI = 'I'
	PieceO = 'O'
	PieceT = 'T'
	PieceS = 'S'
	PieceZ = 'Z'
	PieceJ = 'J'
	PieceL = 'L'
)

// Tetromino represents a tetris piece
type Tetromino [][]rune

// TetrominoPiece represents a tetris piece and its display character
type TetrominoPiece struct {
	shape Tetromino
	char  rune
}

// Tetromino shapes
var tetrominoes = []TetrominoPiece{
	{ // I
		shape: Tetromino{
			[]rune("....."),
			[]rune("....."),
			[]rune("####."),
			[]rune("....."),
			[]rune("....."),
		},
		char: PieceI,
	},
	{ // O
		shape: Tetromino{
			[]rune("....."),
			[]rune("..##."),
			[]rune("..##."),
			[]rune("....."),
			[]rune("....."),
		},
		char: PieceO,
	},
	{ // T
		shape: Tetromino{
			[]rune("....."),
			[]rune("..#.."),
			[]rune(".###."),
			[]rune("....."),
			[]rune("....."),
		},
		char: PieceT,
	},
	{ // S
		shape: Tetromino{
			[]rune("....."),
			[]rune("..##."),
			[]rune(".##.."),
			[]rune("....."),
			[]rune("....."),
		},
		char: PieceS,
	},
	{ // Z
		shape: Tetromino{
			[]rune("....."),
			[]rune(".##.."),
			[]rune("..##."),
			[]rune("....."),
			[]rune("....."),
		},
		char: PieceZ,
	},
	{ // J
		shape: Tetromino{
			[]rune("....."),
			[]rune(".#..."),
			[]rune(".###."),
			[]rune("....."),
			[]rune("....."),
		},
		char: PieceJ,
	},
	{ // L
		shape: Tetromino{
			[]rune("....."),
			[]rune("...#."),
			[]rune(".###."),
			[]rune("....."),
			[]rune("....."),
		},
		char: PieceL,
	},
}

type Game struct {
	board        [][]string
	currentPiece Tetromino
	currentChar  rune
	pieceX       int
	pieceY       int
	score        int
	gameOver     bool
	nextPiece    Tetromino
	nextChar     rune
}

func NewGame() *Game {
	board := make([][]string, boardHeight)
	for i := range board {
		board[i] = make([]string, boardWidth)
		for j := range board[i] {
			board[i][j] = emptyCell
		}
	}

	return &Game{
		board: board,
		score: 0,
	}
}

func (g *Game) spawnPiece() {
	// If we have a next piece, use it
	if g.nextPiece != nil {
		g.currentPiece = g.nextPiece
		g.currentChar = g.nextChar
	} else {
		// First piece of the game
		piece := tetrominoes[rand.Intn(len(tetrominoes))]
		g.currentPiece = piece.shape
		g.currentChar = piece.char
	}

	// Generate next piece
	piece := tetrominoes[rand.Intn(len(tetrominoes))]
	g.nextPiece = piece.shape
	g.nextChar = piece.char
	
	// Set initial position
	g.pieceX = boardWidth/2 - 2
	g.pieceY = 0
}

func (g *Game) rotate() bool {
	if g.currentPiece == nil {
		return false
	}

	// Create a new rotated piece
	size := len(g.currentPiece)
	rotated := make(Tetromino, size)
	for i := range rotated {
		rotated[i] = make([]rune, size)
	}

	// Rotate 90 degrees clockwise
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			rotated[x][size-1-y] = g.currentPiece[y][x]
		}
	}

	// Save the current piece
	oldPiece := g.currentPiece

	// Try the rotation
	g.currentPiece = rotated
	if g.checkCollision() {
		// If collision occurs, revert the rotation
		g.currentPiece = oldPiece
		return false
	}

	return true
}

func (g *Game) draw() {
	clearScreen()

	// Draw border and board
	fmt.Println("+" + strings.Repeat("-", boardWidth*2) + "+")

	// Create a temporary board for drawing
	tempBoard := make([][]string, len(g.board))
	for i := range g.board {
		tempBoard[i] = make([]string, len(g.board[i]))
		copy(tempBoard[i], g.board[i])
	}

	// Draw current piece on temp board
	if g.currentPiece != nil {
		for y := 0; y < len(g.currentPiece); y++ {
			for x := 0; x < len(g.currentPiece[y]); x++ {
				if g.currentPiece[y][x] == '#' {
					newY := g.pieceY + y
					newX := g.pieceX + x
					if newY >= 0 && newY < boardHeight && newX >= 0 && newX < boardWidth {
						tempBoard[newY][newX] = string(g.currentChar)
					}
				}
			}
		}
	}

	// Draw the board
	for _, row := range tempBoard {
		fmt.Print("|")
		for _, cell := range row {
			fmt.Printf("%s ", cell)
		}
		fmt.Println("|")
	}

	fmt.Println("+" + strings.Repeat("-", boardWidth*2) + "+")
	
	// Draw next piece preview
	fmt.Println("\nNext piece:")
	for _, row := range g.nextPiece {
		for _, cell := range row {
			if cell == '.' {
				fmt.Print("  ")
			} else {
				fmt.Printf("%c ", g.nextChar)
			}
		}
		fmt.Println()
	}
	
	fmt.Printf("\nScore: %d\n", g.score)
	fmt.Println("\nControls:")
	fmt.Println("← → : Move left/right")
	fmt.Println("↓ : Move down")
	fmt.Println("↑ : Rotate")
	fmt.Println("q : Quit")
}

func clearScreen() {
	os.Stdout.Sync()
	// Use simple ANSI escape sequences
	fmt.Print("\033[2J\033[H") // Clear screen and move cursor to top-left
}

func restoreScreen() {
	fmt.Print("\033[2J\033[H") // Clear screen and move cursor to top-left
}

func (g *Game) moveLeft() bool {
	g.pieceX--
	if g.checkCollision() {
		g.pieceX++
		return false
	}
	return true
}

func (g *Game) moveRight() bool {
	g.pieceX++
	if g.checkCollision() {
		g.pieceX--
		return false
	}
	return true
}

func (g *Game) moveDown() bool {
	g.pieceY++
	if g.checkCollision() {
		g.pieceY--
		g.lockPiece()
		return false
	}
	return true
}

func (g *Game) checkCollision() bool {
	for y := 0; y < len(g.currentPiece); y++ {
		for x := 0; x < len(g.currentPiece[y]); x++ {
			if g.currentPiece[y][x] == '#' {
				newY := g.pieceY + y
				newX := g.pieceX + x

				if newX < 0 || newX >= boardWidth || newY >= boardHeight {
					return true
				}

				if newY >= 0 && g.board[newY][newX] != emptyCell {
					return true
				}
			}
		}
	}
	return false
}

func (g *Game) lockPiece() {
	for y := 0; y < len(g.currentPiece); y++ {
		for x := 0; x < len(g.currentPiece[y]); x++ {
			if g.currentPiece[y][x] == '#' {
				newY := g.pieceY + y
				newX := g.pieceX + x
				if newY >= 0 && newY < boardHeight && newX >= 0 && newX < boardWidth {
					g.board[newY][newX] = string(g.currentChar)
				}
			}
		}
	}

	g.clearLines()
	g.spawnPiece()
	if g.checkCollision() {
		g.gameOver = true
	}
}

func (g *Game) clearLines() {
	for y := boardHeight - 1; y >= 0; y-- {
		full := true
		for x := 0; x < boardWidth; x++ {
			if g.board[y][x] == emptyCell {
				full = false
				break
			}
		}

		if full {
			// Move all lines above down
			for y2 := y; y2 > 0; y2-- {
				for x := 0; x < boardWidth; x++ {
					g.board[y2][x] = g.board[y2-1][x]
				}
			}
			// Clear top line
			for x := 0; x < boardWidth; x++ {
				g.board[0][x] = emptyCell
			}
			g.score += 100
			y++ // Check the same line again as everything moved down
		}
	}
}

func main() {
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Println("Warning: Not running in a real terminal!")
	}

	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Open keyboard
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		keyboard.Close()
		restoreScreen()
	}()

	game := NewGame()
	game.spawnPiece()

	// Create a ticker for the game loop
	ticker := time.NewTicker(frameRate)
	defer ticker.Stop()

	// Create a channel for keyboard events
	keyChan := make(chan keyboard.Key)
	go func() {
		for {
			char, key, err := keyboard.GetKey()
			if err != nil {
				panic(err)
			}
			if char == 'q' {
				game.gameOver = true
				return
			}
			keyChan <- key
		}
	}()

	// Game loop
	for !game.gameOver {
		game.draw()

		select {
		case key := <-keyChan:
			switch key {
			case keyboard.KeyArrowLeft:
				game.moveLeft()
			case keyboard.KeyArrowRight:
				game.moveRight()
			case keyboard.KeyArrowDown:
				game.moveDown()
			case keyboard.KeyArrowUp:
				game.rotate()
			}
		case <-ticker.C:
			game.moveDown()
		}
	}

	fmt.Println("Game Over! Final score:", game.score)
}
