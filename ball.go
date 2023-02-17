package main

import (
	"errors"
	"github.com/drpaneas/rect"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"log"
	"math"
)

const maxBallSpeed = 15

// Ball is a struct that holds information about the ball in the game
type Ball struct {
	// The position of the ball on the screen
	position *rect.Rectangle

	// The velocity (movement) of the ball
	velocity *Vector2D

	// sounds map
	sounds map[string]*Sound
}

// NewBall creates a new ball with the default values
// The ball is 20x20 pixels and is placed in the middle of the screen
// The ball has a velocity of 0 (not moving) in both directions
func newBall() *Ball {
	b := &Ball{
		position: rect.Rect(halfGameScreenWidth-20/2, halfGameScreenHeight-20/2, 20, 20),
		velocity: &Vector2D{X: 0, Y: 0},
	}

	var err error
	b.sounds, err = LoadSounds()
	if err != nil {
		errSound := errors.New("error loading sounds")
		log.Fatal(errors.Join(errSound, err))
	}

	return b
}

// Draw draws the ball on the screen
func (b *Ball) Draw(screen *ebiten.Image) {
	// draw ball
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.position.X), float64(b.position.Y))
	vector.DrawFilledRect(screen, float32(b.position.X), float32(b.position.Y), float32(b.position.Width), float32(b.position.Height), color.White)
}

// Update updates the position of the ball based on its current velocity.
func (b *Ball) Update() {
	b.position.X += int(math.Round(b.velocity.X))
	b.position.Y += int(math.Round(b.velocity.Y))
}

func (b *Ball) handleBallWallCollision() {
	// Check if ball goes out of screen
	if b.position.Top() < 0 || b.position.Bottom() > screenHeight {
		if err := b.playSound("wall"); err != nil {
			return
		}
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
	directionX := randomChoice(randFloat(-2, -1), randFloat(1, 2))
	directionY := randFloat(-2, 2)

	reducer := 0.25
	b.velocity.X = maxBallSpeed * reducer * directionX
	b.velocity.Y = maxBallSpeed * reducer * directionY

}

func (b *Ball) normalizeBallSpeed() {
	// Calculate the total speed of the ball in pixels per frame
	speed := math.Sqrt(math.Pow(b.velocity.X, 2) + math.Pow(b.velocity.Y, 2))

	// Normalize the ball speed if it's larger than desired
	if speed > maxBallSpeed {
		// Adjust the X and Y components of velocity
		factor := maxBallSpeed / speed
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
func (b *Ball) accelerate(amount float64) {
	signX := 1.0
	if b.velocity.X < 0 {
		signX = -1
	}
	signY := 1.0
	if b.velocity.Y < 0 {
		signY = -1
	}
	b.velocity.X = signX * maxBallSpeed * amount
	b.velocity.Y = signY * maxBallSpeed * amount
}

func (b *Ball) playSound(name string) error {
	if s, ok := b.sounds[name]; ok {
		if err := s.Play(); err != nil {
			return err
		}
	}
	return nil
}
