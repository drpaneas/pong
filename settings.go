package main

import "golang.org/x/image/font"

// Game constants (global)
const (
	screenWidth          = 1280
	screenHeight         = 720
	halfGameScreenWidth  = screenWidth / 2
	halfGameScreenHeight = screenHeight / 2
)

// Game variables (global)
var (
	// Scoring and font
	scoreDisplayFont, resultDisplayFont font.Face
)
