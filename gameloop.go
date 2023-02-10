package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
	"math/rand"
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
			g.ball.position.Center(halfGameScreenWidth, halfGameScreenHeight)

			// Serve the ball to a random side, with lower speed,
			g.ball.setInitialVelocity()

		}
	}

	// Move ball based on its velocity
	g.ball.Update()

	// Check of events:
	// 1: Player scores
	if g.ball.position.Left() <= 0 {
		g.player.score++
		if g.player.score == 10 {
			g.state = gameOver
		} else {
			g.startNewRound()
		}
	}

	// 2. Enemy scores
	if g.ball.position.Right() >= screenWidth {
		g.enemy.score++
		if g.enemy.score == 10 {
			g.state = gameOver
		} else {
			g.startNewRound()
		}
	}

	// 3. Ball hits player paddle
	if g.ball.position.Overlaps(g.player.paddle.position) {
		g.volleyCount++
		g.ball.position.Right(g.player.paddle.position.Left()) // move ball so it touches paddle
		g.ball.accelerate()                                    // faster ball to make the game more interesting
		g.player.bounce(g.ball, g.volleyCount)
		g.state = playing
		g.ball.normalizeBallSpeed() // normalize ball speed
	}

	// 4. Ball hits enemy paddle
	if g.ball.position.Overlaps(g.enemy.paddle.position) {
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
	var tmp float64
	if rand.Float64() > 0.1 {
		tmp = g.enemy.paddle.speed * -1
	} else {
		tmp = g.enemy.paddle.speed
	}
	// 1.a. Move enemy down, if it is higher than the ball
	if g.enemy.paddle.position.Top() < g.ball.position.Bottom() {
		g.enemy.paddle.position.Y += int(math.Round(g.enemy.paddle.speed - tmp))
	}

	// 1.b. Move enemy up, if it is lower than the ball
	if g.enemy.paddle.position.Bottom() > g.ball.position.Top() {
		g.enemy.paddle.position.Y -= int(math.Round(g.enemy.paddle.speed - tmp))
	}

	// 2: Move the enemy paddle to the previously decided location
	g.enemy.Update()

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
