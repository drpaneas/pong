package main

import (
	"github.com/drpaneas/rect"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

// Player is a struct that holds information about the player's paddle and score
type Player struct {
	// The player's paddle
	paddle *Paddle

	// The player's score
	score int
}

func newPlayer() *Player {
	return &Player{
		paddle: &Paddle{
			position: rect.Rect(screenWidth-70-20, halfGameScreenHeight-110/2, 20, 110),
			velocity: &Vector2D{X: 0, Y: 0},
			speed:    15.0,
		},
		score: 0,
	}
}

func (p *Player) GetPaddle() *Paddle {
	return p.paddle
}

func (player *Player) Draw(screen *ebiten.Image) {
	player.paddle.Draw(screen)
}

func (player *Player) Update() {
	// 1. Get the player input and update the paddle velocity
	player.paddle.input()

	// 2. Update the paddle position based on its velocity
	player.paddle.position.Y += int(math.Round(player.paddle.velocity.Y))

	// 3. Check if the player is out of the screen
	// and if so, set the paddle position to the screen edge
	if player.paddle.position.Top() < 0 {
		player.paddle.position.Top(0)
	}

	if player.paddle.position.Bottom() > screenHeight {
		player.paddle.position.Bottom(screenHeight)
	}
}

func (player *Player) bounce(ball *Ball, volleyCount int) {
	// Chop the player in 8 parts and assign a different angle to each part
	ball.velocity.X *= -1 // Reverse the ball direction on X axis
	part := player.paddle.position.Height / 8

	var sl []float64
	if volleyCount < 4 {
		sl = []float64{-135, -150, -165, -180, -180, 165, 150, 135}
	} else if volleyCount >= 4 && volleyCount < 8 {
		sl = []float64{-150, -165, -180, -180, 165, 150, 135, 120}
	} else {
		sl = []float64{-165, -180, -180, 165, 150, 135, 120, 105}
	}

	for i := 0; i < 8; i++ {
		if ball.position.Top() < player.paddle.position.Top()+part*(i+1) {
			ball.velocity.Y = ball.atAngle(sl[i])
			break
		}
	}

}
