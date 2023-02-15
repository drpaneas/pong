package main

import (
	"github.com/drpaneas/rect"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

type PaddleHolder interface {
	GetPaddle() *Paddle
}

// Paddle is a struct that holds information about a paddle in the game
type Paddle struct {
	// The position of the paddle on the screen
	position *rect.Rectangle

	// The velocity (movement) of the paddle
	velocity *Vector2D

	// The speed of the paddle
	speed float64
}

func (p *Paddle) GetPaddle() *Paddle {
	return p
}

func (p *Paddle) Draw(screen *ebiten.Image) {
	// draw player
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(p.position.X), float64(p.position.X))
	vector.DrawFilledRect(screen, float32(p.position.X), float32(p.position.Y), float32(p.position.Width), float32(p.position.Height), color.White)
}
