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
			speed:    13,
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
	e.paddle.position.X += int(math.Round(e.paddle.velocity.X))
	e.paddle.position.Y += int(math.Round(e.paddle.velocity.Y))

	// Check if the enemy is out of the screen
	if e.paddle.position.Top() < 0 {
		e.paddle.position.Top(0)
		//// Check if the velocity is negative and if so, reverse it
		//if e.paddle.velocity.Y < 0 {
		//	e.paddle.velocity.Y *= -1
		//}
	}

	if e.paddle.position.Bottom() > screenHeight {
		e.paddle.position.Bottom(screenHeight)
		//// Check if the velocity is positive and if so, reverse it
		//if e.paddle.velocity.Y > 0 {
		//	e.paddle.velocity.Y *= -1
		//}
	}

}

// bounce is making the ball bounce on the enemy paddle
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

// patrol is making the enemy paddle go randomly up and down
// taking into account the paddle's speed (to avoid jittering)
func (e *Enemy) patrol() {
	if e.randomPosition == 0 {
		halfPaddle := e.paddle.position.Height / 2
		e.randomPosition = randInt(0+halfPaddle, screenHeight-halfPaddle)
	}

	offset := 10
	if e.paddle.position.CenterY() >= e.randomPosition-offset && e.paddle.position.CenterY() <= e.randomPosition+offset {
		e.paddle.velocity.Y = 0
		e.randomPosition = 0
	} else {
		// if the distance is less than the speed, move the paddle to the random position
		if math.Abs(float64(e.randomPosition)-float64(e.paddle.position.CenterY())) < e.paddle.speed {
			e.paddle.position.CenterY(e.randomPosition)
			e.randomPosition = 0
			return
		}
		if e.paddle.position.CenterY() < e.randomPosition {
			e.paddle.velocity.Y = e.paddle.speed
		} else {
			// if the distance is less than the speed, move the paddle to the random position
			e.paddle.velocity.Y = -e.paddle.speed
		}
	}

}
