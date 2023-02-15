package main

import (
	"github.com/drpaneas/rect"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

// Enemy is a struct that holds information about the enemy's paddle and score
type Enemy struct {
	// The enemy's paddle
	paddle *Paddle

	// The enemy's score
	score int

	// Random goto position during patrol
	randomPosition int
}

func newEnemy() *Enemy {
	return &Enemy{
		paddle: &Paddle{
			position: rect.Rect(70, halfGameScreenHeight-110/2, 20, 110),
			velocity: &Vector2D{X: 0, Y: 0},
			speed:    14.0,
		},
		score: 0,
	}
}

func (enemy *Enemy) Draw(screen *ebiten.Image) {
	enemy.paddle.Draw(screen)
}

func (enemy *Enemy) Update() {
	// Check if the enemy is out of the screen
	if enemy.paddle.position.Top() < 0 {
		enemy.paddle.position.Top(0)
	}

	if enemy.paddle.position.Bottom() > screenHeight {
		enemy.paddle.position.Bottom(screenHeight)
	}
}

func (enemy *Enemy) bounce(ball *Ball, volleyCount int) {
	ball.velocity.X *= -1 // reverse the ball direction on X axis
	part := float64(enemy.paddle.position.Height / 8.0)
	var sl []float64
	if volleyCount < 4 {
		sl = []float64{-45, -30, -15, 0, 0, 15, 30, 45}
	} else if volleyCount >= 4 && volleyCount < 8 {
		sl = []float64{-60, -45, -30, -15, 0, 0, 15, 30}
	} else if volleyCount >= 8 {
		sl = []float64{-75, -60, -45, -30, -15, 0, 15, 30}
	} else {
		sl = []float64{-90, -45, -30, -15, 0, 15, 30, 45}
	}

	for i := 0; i < 8; i++ {
		if ball.position.Top() < enemy.paddle.position.Top()+int(math.Round(part*(float64(i)+1))) {
			ball.velocity.Y = ball.atAngle(sl[i])
			break
		}
	}
}

// patrol() is making the enemy paddle go randomly up and down
// if it is reaching the top or bottom of the screen, it will change direction
func (enemy *Enemy) patrol() {
	if enemy.randomPosition == 0 {
		halfPaddle := enemy.paddle.position.Height / 2
		enemy.randomPosition = randInt(0+halfPaddle, screenHeight-halfPaddle)
	}

	// the position is higher than the paddle, move up
	if enemy.randomPosition < enemy.paddle.position.CenterY() {
		enemy.paddle.position.Y -= int(math.Round(enemy.paddle.speed))
	}

	// the position is lower than the paddle, move down
	if enemy.randomPosition > enemy.paddle.position.CenterY() {
		enemy.paddle.position.Y += int(math.Round(enemy.paddle.speed))
	}

	// if the enemy paddle Y is +-15from the random position, calculate a new random position
	if enemy.randomPosition-15 < enemy.paddle.position.CenterY() && enemy.randomPosition+15 > enemy.paddle.position.CenterY() {
		enemy.randomPosition = 0
	}

}
