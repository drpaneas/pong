package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"math"
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

	// The current state of the game (playing, paused, etc)
	state GameState

	// The current turn of the player (player or enemy)
	playerTurn playerTurn

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

	// a timer to delay the patrol movement of the enemy paddle
	timer int

	// HUD for the game (used to display score and the result)
	hud *HUD
}

func newGame() *Game {
	newHud, err := newHUD()
	if err != nil {
		log.Fatal(err)
	}

	// Create the game
	game := &Game{
		state:  firstService,
		ball:   newBall(),
		player: newPlayer(),
		enemy:  newEnemy(),
		hud:    newHud,
	}

	// Add the objects to the gameObjects slice
	game.gameObjects = append(game.gameObjects, game.ball, game.player, game.enemy)

	return game
}

// startNewRound begins a new round of the game (should be called after a player or enemy scores).
// It places the ball back in the center of the screen and serves it to a random direction with a lower speed.
func (g *Game) startNewRound() {
	g.volleyCount = 0 // reset the volley count

	// Stop the ball
	g.ball.velocity.X = 0
	g.ball.velocity.Y = 0

	// Place the ball in the center of the screen
	g.ball.position.Center(halfGameScreenWidth, randInt(20, screenHeight-20))

	// Serve the ball to a random side, with lower speed,
	g.ball.setInitialVelocity()

}

// checkWinCondition checks if either the player or enemy has won the game.
func (g *Game) checkWinCondition() {
	if g.player.score == pointsToWin {
		g.state = gameOver
	} else if g.enemy.score == pointsToWin {
		g.state = gameOver
	} else {
		g.startNewRound()
	}
}

func (g *Game) isGameOver() bool {
	return g.player.score == pointsToWin || g.enemy.score == pointsToWin
}

// handleEnemyAttack handles the enemy's paddle movement.
// If the ball is on the enemy's side of the screen, the enemy will move towards the ball.
// If the ball is on the player's side of the screen, the enemy will patrol back and forth.
// The enemy will only patrol every other frame to make it easier to hit the ball.
// This is done by using the timer variable.
// The timer is incremented every frame and when it reaches 60, it is reset to 0.
// This means that the enemy will patrol every other second.
// This is done to make the game more enjoyable.
func (g *Game) handleEnemyAttack() {
	if g.ball.position.CenterX() < halfGameScreenWidth {
		if g.ball.position.Bottom() < g.enemy.paddle.position.CenterY() {
			g.enemy.paddle.position.Y -= int(math.Round(g.enemy.paddle.speed))
		}
		if g.ball.position.Top() > g.enemy.paddle.position.CenterY() {
			g.enemy.paddle.position.Y += int(math.Round(g.enemy.paddle.speed))
		}
	} else {
		if g.timer%2 == 0 {
			g.enemy.patrol()
		}
	}
}

// updateTimer updates the timer by 1 every frame.
func (g *Game) updateTimer() {
	g.timer++
	if g.timer > 60 {
		g.timer = 0
	}
}

// handleBallCollision handles the collision of the ball with the paddles only.
func (g *Game) handlePaddleCollision(holder PaddleHolder) error {
	if err := g.ball.playSound("paddle"); err != nil {
		return err
	}

	g.volleyCount++
	g.ball.accelerate(1)

	switch holder.GetPaddle() {
	case g.player.paddle:
		g.playerTurn = playerTurnEnemy
		g.ball.position.Right(g.player.paddle.position.Left())
		g.player.bounce(g.ball, g.volleyCount)
	case g.enemy.paddle:
		g.playerTurn = playerTurnPlayer
		g.ball.position.Left(g.enemy.paddle.position.Right())
		g.enemy.bounce(g.ball, g.volleyCount)
	}

	return nil
}

// handleFirstService handles the first service of the game.
// The first service is when the ball is in the center of the screen and not moving.
// When the ball is in this state, the game will serve the ball to a random direction.
// The ball will also be given a random speed.
// The game will then change to the playing state.
func (g *Game) handleFirstService() error {
	if g.ball.velocity.X == 0 && g.ball.velocity.Y == 0 {
		g.volleyCount = 0
		g.ball.setInitialVelocity()

		if g.ball.velocity.X < 0.0 {
			g.playerTurn = playerTurnEnemy
		} else {
			g.playerTurn = playerTurnPlayer
		}
		g.state = playing
	}

	return nil
}

// handleScore handles the scoring of the game.
//  1. If the ball goes off the left side of the screen, the player scores.
//  2. If the ball goes off the right side of the screen, the enemy scores.
//  3. If either player scores, the game checks if the game is over.
func (g *Game) handleScore() error {
	if g.ball.position.Left() <= 0 {
		if err := g.ball.playSound("score"); err != nil {
			return err
		}
		g.player.score++
		g.checkWinCondition()
	}

	if g.ball.position.Right() >= screenWidth {
		if err := g.ball.playSound("score"); err != nil {
			return err
		}
		g.enemy.score++
		g.checkWinCondition()
	}
	return nil
}
