package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// function to handle user input controlling the paddle up and down
func (p *Paddle) input() {
	// Up
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		p.velocity.Y = p.velocity.Y - p.speed
	} else if inpututil.IsKeyJustReleased(ebiten.KeyArrowUp) {
		p.velocity.Y = p.velocity.Y + p.speed
	}

	// Down
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		p.velocity.Y = p.velocity.Y + p.speed
	} else if inpututil.IsKeyJustReleased(ebiten.KeyArrowDown) {
		p.velocity.Y = p.velocity.Y - p.speed
	}

}
