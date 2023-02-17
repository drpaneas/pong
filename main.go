package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	// Configure the game window
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	ebiten.SetWindowTitle("Pong")
	ebiten.SetFullscreen(false)

	game := newGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
