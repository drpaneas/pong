package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

// GameObject is considered anything that can be updated and drawn on the screen
type GameObject interface {
	Update()
	Draw(screen *ebiten.Image)
}

// Game is the main struct for our game that holds all the important information
type Game struct {
	// The score in the game
	score Score

	// The current state of the game (playing, paused, etc)
	state GameState

	// The current turn of the player (user or computer)
	turn playerTurn

	// The number of times the ball has been hit back and forth
	// the more times it is hit, the faster it goes to increase the difficulty
	volleyCount int

	// The ball in the game
	ball *Ball

	// The player's paddle
	player *Player

	// The enemy's paddle
	enemy *Enemy

	// A slice to store all the game objects (ball, player, enemy)
	// used to update and draw them all at once
	objects []GameObject

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

	// Add the objects to the objects slice
	game.objects = append(game.objects, game.ball, game.player, game.enemy)

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
	if g.score.player == pointsToWin {
		g.state = gameOver
	} else if g.score.enemy == pointsToWin {
		g.state = gameOver
	} else {
		g.startNewRound()
	}
}

func (g *Game) isGameOver() bool {
	return g.score.player == pointsToWin || g.score.enemy == pointsToWin
}

// handleEnemyAttack handles the enemy's AI paddle movement.
func (g *Game) handleEnemyAttack() {
	// Calculate in which Y there will be collision
	// slope of the ball's trajectory
	slope := g.ball.velocity.Y / g.ball.velocity.X

	// Y-intercept of the ball's trajectory
	yIntercept := float64(g.ball.position.Y) - slope*float64(g.ball.position.X)

	// predict the Y position of the ball when it reaches the center of the paddle
	predictedY := slope*float64(g.enemy.paddle.position.X) + yIntercept

	// Check if the paddle is already at the predicted Y position
	// taking into account the paddle's speed (to avoid jittering)
	offset := 10
	if g.enemy.paddle.position.CenterX() >= int(predictedY)-offset && g.enemy.paddle.position.CenterX() <= int(predictedY)+offset {
		// stop moving
		g.enemy.paddle.velocity.Y = 0
		return
	} else {
		// If the paddle is not at the predicted Y position, move it towards the predicted Y position
		// If the paddle is lower than the predicted Y position, move it up
		if g.enemy.paddle.position.CenterY() > int(predictedY) {
			// if the distance is less than the paddle's speed, stop
			if g.enemy.paddle.position.CenterY()-int(predictedY) < int(g.enemy.paddle.speed) {
				g.enemy.paddle.velocity.Y = 0
				return
			}
			// move it up
			g.enemy.paddle.velocity.Y = randFloat(-g.enemy.paddle.speed/2, -g.enemy.paddle.speed)
		}

		// If the paddle is higher than the predicted Y position, move it down
		if g.enemy.paddle.position.CenterY() < int(predictedY) {
			// if the distance is less than the paddle's speed, stop
			if int(predictedY)-g.enemy.paddle.position.CenterY() < int(g.enemy.paddle.speed) {
				g.enemy.paddle.velocity.Y = 0
				return
			}
			// move it down
			g.enemy.paddle.velocity.Y = randFloat(g.enemy.paddle.speed/2, g.enemy.paddle.speed)
		}
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
		g.turn = computer
		g.ball.position.Right(g.player.paddle.position.Left())
		g.player.bounce(g.ball, g.volleyCount)
	case g.enemy.paddle:
		g.turn = user
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
		g.score.player++
		g.checkWinCondition()
	}

	if g.ball.position.Right() >= screenWidth {
		if err := g.ball.playSound("score"); err != nil {
			return err
		}
		g.score.enemy++
		g.checkWinCondition()
	}
	return nil
}
