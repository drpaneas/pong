package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
)

func (g *Game) Update() error {
	if g.state == paused {
		return nil
	}
	if g.state == gameOver {
		return nil
	}

	if g.state == firstService && g.ball.velocity.X == 0 && g.ball.velocity.Y == 0 { // if ball is not moving anymore
		{
			g.volleyCount = 0 // reset the volley count

			// Serve the ball to a random side, with lower speed,
			g.ball.setInitialVelocity()

			// if the velocity is going to the left side (enemy), then set the playerTurn to enemy
			if g.ball.velocity.X < 0.0 {
				g.playerTurn = playerTurnEnemy
			} else {
				g.playerTurn = playerTurnPlayer
			}

		}
	}

	if g.state == playing || g.state == firstService {

		// Move ball based on its velocity
		g.ball.Update()

		// if the velocity is going to the left side (enemy), then set the playerTurn to enemy
		if g.ball.velocity.X < 0.0 {
			g.playerTurn = playerTurnEnemy
		} else {
			g.playerTurn = playerTurnPlayer
		}

		// Check of events:
		// 1: Player scores
		if g.ball.position.Left() <= 0 {
			if err := g.ball.playSound("score"); err != nil {
				return err
			}
			g.player.score++
			if g.player.score == 10 {
				g.state = gameOver
			} else {
				g.startNewRound()
				return nil
			}
		}

		// 2. Enemy scores
		if g.ball.position.Right() >= screenWidth {
			if err := g.ball.playSound("score"); err != nil {
				return err
			}
			g.enemy.score++
			if g.enemy.score == 10 {
				g.state = gameOver
			} else {
				g.startNewRound()
				return nil
			}
		}

		// 3. Ball hits player paddle
		if g.ball.position.Overlaps(g.player.paddle.position) {
			if err := g.ball.playSound("bounce"); err != nil {
				return err
			}
			g.playerTurn = playerTurnEnemy // enemy has to play next
			g.volleyCount++
			g.ball.position.Right(g.player.paddle.position.Left()) // move ball so it touches paddle
			g.ball.accelerate()                                    // faster ball to make the game more interesting
			g.player.bounce(g.ball, g.volleyCount)
			g.state = playing
			g.ball.normalizeBallSpeed() // normalize ball speed
		}

		// 4. Ball hits enemy paddle
		if g.ball.position.Overlaps(g.enemy.paddle.position) {
			if err := g.ball.playSound("bounce"); err != nil {
				return err
			}
			g.playerTurn = playerTurnPlayer // player has to play next
			g.volleyCount++
			g.ball.position.Left(g.enemy.paddle.position.Right()) // move ball so it touches paddle
			g.ball.accelerate()                                   // faster ball to make the game more interesting
			g.enemy.bounce(g.ball, g.volleyCount)
			g.state = playing
			g.ball.normalizeBallSpeed() // normalize ball speed
		}

		// Move the player paddle based on user input
		g.player.Update()

		// Move the enemy paddle based on the AI
		// 1: Decide where to move the enemy paddle based on the ball position
		if g.playerTurn == playerTurnEnemy {
			// if the ball is at the enemy side of the screen, then move the paddle to the ball
			if g.ball.position.CenterX() < halfGameScreenWidth {
				// 1.a. If the ball is higher than the enemy paddle, move the enemy paddle up
				if g.ball.position.Bottom() < g.enemy.paddle.position.CenterY() {
					g.enemy.paddle.position.Y -= int(math.Round(g.enemy.paddle.speed))
				}

				// 1.b. If the ball is lower than the enemy paddle, move the enemy paddle down
				if g.ball.position.Top() > g.enemy.paddle.position.CenterY() {
					g.enemy.paddle.position.Y += int(math.Round(g.enemy.paddle.speed))
				}
			} else {
				if g.timer%2 == 0 {
					g.enemy.patrol()
				}
			}
		}

		if g.playerTurn == playerTurnPlayer {
			if g.timer%2 == 0 {
				g.enemy.patrol()
			}
		}

		g.timer++
		if g.timer > 60 {
			g.timer = 0
		}

		// 2: Move the enemy paddle to the previously decided location
		g.enemy.Update()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// draw dashed line in the middle (dimensions 10x60 per dash and 40px space between dashes)
	for i := 0; i < screenHeight; i += 100 {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(halfGameScreenWidth), float64(i))
		vector.StrokeLine(screen, float32(halfGameScreenWidth), float32(i), float32(halfGameScreenWidth), float32(i+60), 10, color.White)
	}

	// Loop through the gameObjects slice and call the Draw function for each object
	for _, obj := range g.gameObjects {
		obj.Draw(screen)
	}

	// draw score
	text.Draw(screen, fmt.Sprintf("%d", g.enemy.score), scoreDisplayFont, halfGameScreenWidth-360, 120, color.White)
	text.Draw(screen, fmt.Sprintf("%d", g.player.score), scoreDisplayFont, halfGameScreenWidth+360-75, 120, color.White)

	if g.state == paused {
		text.Draw(screen, "PAUSED", scoreDisplayFont, halfGameScreenWidth-100, halfGameScreenHeight-100, color.White)
	}

	if g.state == gameOver {
		if g.player.score > g.enemy.score {
			text.Draw(screen, "WINNER", resultDisplayFont, halfGameScreenWidth+450, halfGameScreenHeight, color.White)
			text.Draw(screen, "LOSER", resultDisplayFont, halfGameScreenWidth-450, halfGameScreenHeight, color.White)
		} else {
			text.Draw(screen, "WINNER", resultDisplayFont, halfGameScreenWidth-450, halfGameScreenHeight, color.White)
			text.Draw(screen, "LOSER", resultDisplayFont, halfGameScreenWidth+350, halfGameScreenHeight, color.White)
		}
	}
}

func (g *Game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}
