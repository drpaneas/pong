package main

import (
	"github.com/drpaneas/rect"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math"
)

// Ball is a struct that holds information about the ball in the game
type Ball struct {
	// The position of the ball on the screen
	position *rect.Rectangle

	// The velocity (movement) of the ball
	velocity *Vector2D

	// The speed of the ball
	speed float64
}

// NewBall creates a new ball with the default values
// The ball is 20x20 pixels and is placed in the middle of the screen
// The ball has a speed of 15 pixels per second
// The ball has a velocity of 0 (not moving) in both directions
func newBall() *Ball {
	b := &Ball{
		position: rect.Rect(halfGameScreenWidth-20/2, halfGameScreenHeight-20/2, 20, 20),
		speed:    15.0,
		velocity: &Vector2D{X: 0, Y: 0},
	}

	return b
}

var endY int

// Draw draws the ball on the screen
func (b *Ball) Draw(screen *ebiten.Image) {
	// draw ball
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.position.X), float64(b.position.Y))
	vector.DrawFilledRect(screen, float32(b.position.X), float32(b.position.Y), float32(b.position.Width), float32(b.position.Height), color.White)
}

// Update updates the position of the ball based on its current velocity.
// It also checks if the ball goes out of the screen and reverses its velocity
func (b *Ball) Update() {
	b.position.X += int(math.Round(b.velocity.X))
	b.position.Y += int(math.Round(b.velocity.Y))

	// Check if ball goes out of screen
	if b.position.Top() < 0 || b.position.Bottom() > screenHeight {
		// if ball is below screen, set bottom to screen height
		if b.position.Bottom() >= screenHeight {
			b.position.Bottom(screenHeight)
		} else {
			// if ball is above screen, set top to 0
			b.position.Top(0)
		}
		b.velocity.Y *= -1 // reverse Y axis velocity
	}
}

// setInitialVelocity reduces the ball speed
// This is used when the ball is served to a player for the first time.
func (b *Ball) setInitialVelocity() {
	directionX := randFloat(-0.5, 0.5)
	directionY := randFloat(-0.5, 0.5)

	// Make sure the ball always moves in the X axis
	if directionX == 0 {
		directionX = 0.5
	}

	reducer := 0.7
	b.velocity.X = b.speed * reducer * directionX
	b.velocity.Y = b.speed * reducer * directionY

}

func (b *Ball) normalizeBallSpeed() {
	// Calculate the total speed of the ball in pixels per frame
	speed := math.Sqrt(math.Pow(b.velocity.X, 2) + math.Pow(b.velocity.Y, 2))

	// Normalize the ball speed if it's larger than desired
	if speed > b.speed {
		// Adjust the X and Y components of velocity
		factor := b.speed / speed
		b.velocity.X = b.velocity.X * factor
		b.velocity.Y = b.velocity.Y * factor
	}
}

func (b *Ball) atAngle(angle float64) float64 {
	// Convert the angle to radians
	radians := angle * math.Pi / 180
	// Calculate the new speed for Y axis
	return math.Round(math.Tan(radians) * b.velocity.X)
}

// accelerate increases the ball speed to its maximum value
func (b *Ball) accelerate() {
	signX := 1.0
	if b.velocity.X < 0 {
		signX = -1
	}
	signY := 1.0
	if b.velocity.Y < 0 {
		signY = -1
	}
	b.velocity.X = signX * b.speed
	b.velocity.Y = signY * b.speed
}
