package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

func (g *Game) Update() error {
	switch g.state {
	case paused:
		return nil

	case gameOver:
		return nil

	case firstService:
		if err := g.handleFirstService(); err != nil {
			return err
		}

	case playing:

		if g.ball.velocity.X < 0.0 {
			g.turn = computer
		} else {
			g.turn = user
		}

		// Make the ball speed up after the first 4 volleys
		if g.volleyCount < 4 {
			g.ball.normalizeBallSpeed()
		}

		// Collision logic has 3 parts:
		// 	1. Check if the ball is colliding with the player's paddle
		// 	2. Check if the ball is colliding with the enemy's paddle
		// 	3. Check if the ball is colliding with the top or bottom wall
		if g.ball.position.Overlaps(g.player.paddle.position) {
			if err := g.handlePaddleCollision(g.player.paddle); err != nil {
				return err
			}
		} else if g.ball.position.Overlaps(g.enemy.paddle.position) {
			if err := g.handlePaddleCollision(g.enemy.paddle); err != nil {
				return err
			}
		} else {
			g.ball.handleBallWallCollision()
		}

		// If someone scores,
		//  1. update the score for this guy and reset the ball
		//  2. check if the game is over and if so, change the game state
		if err := g.handleScore(); err != nil {
			return err
		}

		// AI logic for the enemy has two parts:
		// 	1. If the enemy is not serving, it will patrol the screen
		// 	2. If the enemy is serving, it will attack (meaning, it will move towards the ball)
		if g.turn == computer {
			g.handleEnemyAttack()
		} else {
			g.enemy.patrol()
		}

		// Lastly, update the ball, player and enemy positions
		g.ball.Update()
		g.player.Update()
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
	text.Draw(screen, fmt.Sprintf("%d", g.score.enemy), g.hud.ScoreDisplayFont, halfGameScreenWidth-360, 120, color.White)
	text.Draw(screen, fmt.Sprintf("%d", g.score.player), g.hud.ScoreDisplayFont, halfGameScreenWidth+360-75, 120, color.White)

	if g.state == paused {
		text.Draw(screen, "PAUSED", g.hud.ScoreDisplayFont, halfGameScreenWidth-100, halfGameScreenHeight-100, color.White)
	}

	if g.state == gameOver {
		if g.score.player > g.score.enemy {
			text.Draw(screen, "WINNER", g.hud.ResultDisplayFont, halfGameScreenWidth+450, halfGameScreenHeight, color.White)
			text.Draw(screen, "LOSER", g.hud.ResultDisplayFont, halfGameScreenWidth-450, halfGameScreenHeight, color.White)
		} else {
			text.Draw(screen, "WINNER", g.hud.ResultDisplayFont, halfGameScreenWidth-450, halfGameScreenHeight, color.White)
			text.Draw(screen, "LOSER", g.hud.ResultDisplayFont, halfGameScreenWidth+350, halfGameScreenHeight, color.White)
		}
	}
}

func (g *Game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}
