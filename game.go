package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"log"
)

// GameObject is an interface that holds common fields and methods for all game objects
type GameObject interface {
	Update()
	Draw(screen *ebiten.Image)
}

// Game is the main struct for our game that holds all the important information
type Game struct {
	// The player's and enemy's score in the game
	playerScore, enemyScore int

	// The font used to display the score on the screen
	font font.Face

	// The current state of the game (playing, paused, etc)
	state GameState

	// The number of times the ball has been hit back and forth
	// the more times it is hit, the faster it goes to increase the difficulty
	volleyCount int

	// The ball in the game
	ball *Ball

	// The player's paddle
	player *Player

	// The enemy's paddle
	enemy *Enemy

	// A slice to store all the game objects
	gameObjects []GameObject
}

func newGame() *Game {
	// Load the font
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	face, _ := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	// Create the game
	game := &Game{
		font:   face,
		state:  firstService,
		ball:   newBall(),
		player: newPlayer(),
		enemy:  newEnemy(),
	}

	// Add the objects to the gameObjects slice
	game.gameObjects = append(game.gameObjects, game.ball, game.player, game.enemy)

	return game
}

// startNewRound begins a new round of the game (should be called after a player or enemy scores).
// It places the ball back in the center of the screen and serves it to a random direction with a lower speed.
func (g *Game) startNewRound() {
	g.volleyCount = 0 // reset the volley count
	g.ball.position.Center(halfGameScreenWidth, halfGameScreenHeight)

	// Serve the ball to a random side, with lower speed,
	g.ball.setInitialVelocity()
}
