package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// function to handle user input controlling the paddle up and down
func (p *Paddle) input() {
	userMovementSpeed := 15.0 // the speed of the paddle every time the user presses a key

	// Up
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		p.velocity.Y = p.velocity.Y - userMovementSpeed
	} else if inpututil.IsKeyJustReleased(ebiten.KeyArrowUp) {
		p.velocity.Y = p.velocity.Y + userMovementSpeed
	}

	// Down
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		p.velocity.Y = p.velocity.Y + userMovementSpeed
	} else if inpututil.IsKeyJustReleased(ebiten.KeyArrowDown) {
		p.velocity.Y = p.velocity.Y - userMovementSpeed
	}

}
