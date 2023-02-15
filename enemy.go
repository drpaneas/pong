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

// newEnemy creates a new enemy and returns a pointer to it
func newEnemy() *Enemy {
	return &Enemy{
		paddle: &Paddle{
			position: rect.Rect(70, halfGameScreenHeight-110/2, 20, 110),
			velocity: &Vector2D{X: 0, Y: 0},
			speed:    12.0,
		},
		score: 0,
	}
}

// GetPaddle returns the enemy's paddle
func (e *Enemy) GetPaddle() *Paddle {
	return e.paddle
}

// Draw draws the enemy's paddle on the screen
func (e *Enemy) Draw(screen *ebiten.Image) {
	e.paddle.Draw(screen)
}

// Update updates the enemy's paddle
func (e *Enemy) Update() {
	// Check if the enemy is out of the screen
	if e.paddle.position.Top() < 0 {
		e.paddle.position.Top(0)
	}

	if e.paddle.position.Bottom() > screenHeight {
		e.paddle.position.Bottom(screenHeight)
	}
}

// bounce() is making the ball bounce on the enemy paddle
func (e *Enemy) bounce(ball *Ball, volleyCount int) {
	ball.velocity.X *= -1 // reverse the ball direction on X axis
	part := float64(e.paddle.position.Height / 8.0)
	var sl []float64
	if volleyCount < 4 {
		sl = []float64{-45, -30, -15, 0, 0, 15, 30, 45}
	} else if volleyCount >= 4 && volleyCount < 8 {
		sl = []float64{-60, -45, -30, -15, 0, 0, 15, 30}
	} else {
		sl = []float64{-75, -60, -45, -30, -15, 0, 15, 30}
	}

	for i := 0; i < 8; i++ {
		if ball.position.Top() < e.paddle.position.Top()+int(math.Round(part*(float64(i)+1))) {
			ball.velocity.Y = ball.atAngle(sl[i])
			break
		}
	}
}

// patrol() is making the enemy paddle go randomly up and down
// if it is reaching the top or bottom of the screen, it will change direction
func (e *Enemy) patrol() {
	if e.randomPosition == 0 {
		halfPaddle := e.paddle.position.Height / 2
		e.randomPosition = randInt(0+halfPaddle, screenHeight-halfPaddle)
	}

	// the position is higher than the paddle, move up
	if e.randomPosition < e.paddle.position.CenterY() {
		e.paddle.position.Y -= int(math.Round(e.paddle.speed))
	}

	// the position is lower than the paddle, move down
	if e.randomPosition > e.paddle.position.CenterY() {
		e.paddle.position.Y += int(math.Round(e.paddle.speed))
	}

	// if the enemy paddle Y is +-15from the random position, calculate a new random position
	if e.randomPosition-15 < e.paddle.position.CenterY() && e.randomPosition+15 > e.paddle.position.CenterY() {
		e.randomPosition = 0
	}

}
