package main

// GameState is the current state of the game (playing, paused, game over)
type GameState int

const (
	playing GameState = iota
	paused
	gameOver
	firstService
)

type playerTurn int

const (
	playerTurnPlayer playerTurn = iota
	playerTurnEnemy
)
