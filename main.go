package main

import (
	"fmt"
	"github.com/drpaneas/rect"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	"log"
	"math"
	"math/rand"
)

// Structures
// --------------------------------------------------------------------------------------------------------------- //

// Vector2D is a struct that stores X and Y values for a position
type Vector2D struct {
	X float64
	Y float64
}

// Game is the main struct for our game that holds all the important information
type Game struct {
	// The player's and enemy's score in the game
	playerScore, enemyScore int

	// The font used to display the score on the screen
	font font.Face

	// The current state of the game (playing, paused, etc)
	state GameState

	// The number of times the ball has been hit back and forth
	// the more times it is hit, the faster it goes to increase the difficulty
	volleyCount int

	// The ball in the game
	ball *Ball

	// The player's paddle
	player *Player

	// The enemy's paddle
	enemy *Enemy
}

// Ball is a struct that holds information about the ball in the game
type Ball struct {
	// The position of the ball on the screen
	position *rect.Rectangle

	// The velocity (movement) of the ball
	velocity *Vector2D

	// The speed of the ball
	speed float64
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

// Player is a struct that holds information about the player's paddle and score
type Player struct {
	// The player's paddle
	paddle *Paddle

	// The player's score
	score int
}

// Enemy is a struct that holds information about the enemy's paddle and score
type Enemy struct {
	// The enemy's paddle
	paddle *Paddle

	// The enemy's score
	score int
}

// Constructors and methods
// --------------------------------------------------------------------------------------------------------------- //
// newGame creates a new game instance
func newGame() *Game {
	// Load the font
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	face, _ := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	// Create the game
	return &Game{
		font:   face,
		state:  firstService,
		ball:   newBall(),
		player: newPlayer(),
		enemy:  newEnemy(),
	}
}

func newPlayer() *Player {
	return &Player{
		paddle: &Paddle{
			position: rect.Rect(screenWidth-70-20, halfGameScreenHeight-110/2, 20, 110),
			velocity: &Vector2D{X: 0, Y: 0},
			speed:    15.0,
		},
		score: 0,
	}
}

func newEnemy() *Enemy {
	return &Enemy{
		paddle: &Paddle{
			position: rect.Rect(70, halfGameScreenHeight-110/2, 20, 110),
			velocity: &Vector2D{X: 0, Y: 0},
			speed:    15.0,
		},
		score: 0,
	}
}

func newBall() *Ball {
	b := &Ball{
		position: rect.Rect(halfGameScreenWidth-20/2, halfGameScreenHeight-20/2, 20, 20),
		speed:    15.0,
		velocity: &Vector2D{X: 0, Y: 0},
	}

	return b
}

// Implement the Interface for ball and paddle
// --------------------------------------------------------------------------------------------------------------- //
// Ball.Update() and Ball.Draw(), Paddle.Update() and Paddle.Draw()

func (b *Ball) Draw(screen *ebiten.Image) {
	// draw ball
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.position.X), float64(b.position.Y))
	vector.DrawFilledRect(screen, float32(b.position.X), float32(b.position.Y), float32(b.position.Width), float32(b.position.Height), color.White)
}

func (paddle *Paddle) Draw(screen *ebiten.Image) {
	// draw player
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(paddle.position.X), float64(paddle.position.X))
	vector.DrawFilledRect(screen, float32(paddle.position.X), float32(paddle.position.Y), float32(paddle.position.Width), float32(paddle.position.Height), color.White)
}

// --------------------------------------------------------------------------------------------------------------- //

// Game constants (global)
const (
	screenWidth          = 1280
	screenHeight         = 720
	halfGameScreenWidth  = screenWidth / 2
	halfGameScreenHeight = screenHeight / 2
)

// GameState is the current state of the game (playing, paused, game over)
type GameState int

const (
	playing GameState = iota
	paused
	gameOver
	firstService
)

// Game variables (global)
var (
	// Scoring and font
	playerScoreCount, enemyScoreCount   int
	scoreDisplayFont, resultDisplayFont font.Face
)

// Functions
// --------------------------------------------------------------------------------------------------------------- //
// startNewRound begins a new round of the game (should be called after a player or enemy scores).
// It places the ball back in the center of the screen and serves it to a random direction with a lower speed.
func (g *Game) startNewRound() {
	g.volleyCount = 0 // reset the volley count
	g.ball.position.Center(halfGameScreenWidth, halfGameScreenHeight)

	// Serve the ball to a random side, with lower speed,
	g.ball.velocity.X = g.ball.speed * randomChoice(-1, 1) / 3
	g.ball.velocity.Y = g.ball.speed * randomChoice(-1, 1) / 3
}

// move updates the position of the ball based on its current velocity.
func (b *Ball) move() {
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

// setInitialVelocity reduces the ball speed by a factor of 3.
// This is used when the ball is served to a player for the first time.
func (b *Ball) setInitialVelocity() {
	b.velocity.X = b.speed * randomChoice(-1, 1) / 3
	b.velocity.Y = b.speed * randomChoice(-1, 1) / 3
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

func (enemy *Enemy) bounce(ball *Ball, volleyCount int) {
	ball.velocity.X *= -1 // reverse the ball direction on X axis
	part := float64(enemy.paddle.position.Height / 8.0)
	var sl []float64
	if volleyCount < 4 {
		sl = []float64{-45, -30, -15, 0, 0, 15, 30, 45}
	} else if volleyCount >= 4 && volleyCount < 8 {
		sl = []float64{-60, -45, -30, -15, 0, 0, 15, 30}
	} else if volleyCount >= 8 {
		sl = []float64{-75, -60, -45, -30, -15, 0, 15, 30}
	} else {
		sl = []float64{-90, -45, -30, -15, 0, 15, 30, 45}
	}

	for i := 0; i < 8; i++ {
		if ball.position.Top() < enemy.paddle.position.Top()+int(math.Round(part*(float64(i)+1))) {
			ball.velocity.Y = ball.atAngle(sl[i])
			break
		}
	}
}

func (player *Player) bounce(ball *Ball, volleyCount int) {

	// Chop the player in 8 parts and assign a different angle to each part
	ball.velocity.X *= -1 // Reverse the ball direction on X axis
	part := player.paddle.position.Height / 8

	var sl []float64
	if volleyCount < 4 {
		sl = []float64{-135, -150, -165, -180, -180, 165, 150, 135}
	} else if volleyCount >= 4 && volleyCount < 8 {
		sl = []float64{-150, -165, -180, -180, 165, 150, 135, 120}
	} else if volleyCount >= 8 {
		sl = []float64{-165, -180, -180, 165, 150, 135, 120, 105}
	} else {
		sl = []float64{-180, -180, 165, 150, 135, 120, 105, 90}
	}

	for i := 0; i < 8; i++ {
		if ball.position.Top() < player.paddle.position.Top()+part*(i+1) {
			ball.velocity.Y = ball.atAngle(sl[i])
			break
		}
	}

}

func (b *Ball) atAngle(angle float64) float64 {
	// Convert the angle to radians
	radians := angle * math.Pi / 180
	// Calculate the new speed for Y axis
	return math.Round(math.Tan(radians) * b.velocity.X)
}

func (player *Player) move() {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		player.paddle.velocity.Y -= player.paddle.speed
	} else if inpututil.IsKeyJustReleased(ebiten.KeyArrowUp) {
		player.paddle.velocity.Y += player.paddle.speed
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		player.paddle.velocity.Y += player.paddle.speed
	} else if inpututil.IsKeyJustReleased(ebiten.KeyArrowDown) {
		player.paddle.velocity.Y -= player.paddle.speed
	}

	player.paddle.position.Y += int(math.Round(player.paddle.velocity.Y))

	// Check if the player is out of the screen
	if player.paddle.position.Top() < 0 {
		player.paddle.position.Top(0)
	}

	if player.paddle.position.Bottom() > screenHeight {
		player.paddle.position.Bottom(screenHeight)
	}
}

func (enemy *Enemy) move() {

	// Check if the enemy is out of the screen
	if enemy.paddle.position.Top() < 0 {
		enemy.paddle.position.Top(0)
	}

	if enemy.paddle.position.Bottom() > screenHeight {
		enemy.paddle.position.Bottom(screenHeight)
	}

}

// Function that returns randomly either a or b. If a and b are equal, it returns value 'a'.
func randomChoice(a, b float64) float64 {
	if a == b {
		return a
	}
	if rand.Intn(2) == 0 {
		return a
	}
	return b
}

func (g *Game) Update() error {
	if g.state == paused {
		return nil
	}
	if g.state == gameOver {
		return nil
	}

	if g.state == firstService && g.ball.velocity.X == 0 && g.ball.velocity.Y == 0 { // if ball is not moving anymore
		{
			g.volleyCount = 0 // reset the volley count
			g.ball.position.Center(halfGameScreenWidth, halfGameScreenHeight)

			// Serve the ball to a random side, with lower speed,
			g.ball.velocity.X = g.ball.speed * randomChoice(-1, 1) / 3
			g.ball.velocity.Y = g.ball.speed * randomChoice(-1, 1) / 3

		}
	}

	// Update the position of the ball in the game based on its velocity (check for out of screen as well)
	g.ball.position.X += int(math.Round(g.ball.velocity.X))
	g.ball.position.Y += int(math.Round(g.ball.velocity.Y))

	// Check if ball goes out of screen
	if g.ball.position.Top() < 0 || g.ball.position.Bottom() > screenHeight {
		// if ball is below screen, set bottom to screen height
		if g.ball.position.Bottom() >= screenHeight {
			g.ball.position.Bottom(screenHeight)
		} else {
			// if ball is above screen, set top to 0
			g.ball.position.Top(0)
		}
		g.ball.velocity.Y *= -1 // reverse Y axis velocity
	}

	// After moving the ball to the new location, check for various events:

	// Check 1: if player scores
	if g.ball.position.Left() <= 0 {
		playerScoreCount++
		if playerScoreCount == 10 {
			g.state = gameOver
		} else {
			g.startNewRound()
		}
	}

	// Check 2: if enemy scores
	if g.ball.position.Right() >= screenWidth {
		enemyScoreCount++
		if enemyScoreCount == 10 {
			g.state = gameOver
		} else {
			g.startNewRound()
		}
	}

	// Check 3: if ball hits player paddle
	if g.ball.position.Overlaps(g.player.paddle.position) {
		g.volleyCount++
		g.ball.position.Right(g.player.paddle.position.Left()) // move ball so it touches paddle
		g.ball.accelerate()                                    // faster ball to make the game more interesting
		g.player.bounce(g.ball, g.volleyCount)
		g.state = playing
		g.ball.normalizeBallSpeed() // normalize ball speed
	}

	// Check 4: if ball hits enemy paddle
	if g.ball.position.Overlaps(g.enemy.paddle.position) {
		g.volleyCount++
		g.ball.position.Left(g.enemy.paddle.position.Right()) // move ball so it touches paddle
		g.ball.accelerate()                                   // faster ball to make the game more interesting
		g.enemy.bounce(g.ball, g.volleyCount)
		g.state = playing
		g.ball.normalizeBallSpeed() // normalize ball speed
	}

	// Move the player paddle based on the input from the user (if any)
	g.player.move()

	// Move the enemy paddle based on the AI
	// Step 1: Decide where to move the enemy paddle based on the ball position
	var tmp float64
	if rand.Float64() > 0.1 {
		tmp = g.enemy.paddle.speed * -1
	} else {
		tmp = g.enemy.paddle.speed
	}
	// if half of the enemy height is below the center of the ball, move down
	if g.enemy.paddle.position.CenterY() < g.ball.position.CenterY() {
		g.enemy.paddle.position.Y += int(math.Round(g.enemy.paddle.speed - tmp))
	}

	// if half of the enemy height is above the center of the ball, move up
	if g.enemy.paddle.position.CenterY() > g.ball.position.CenterY() {
		g.enemy.paddle.position.Y -= int(math.Round(g.enemy.paddle.speed - tmp))
	}

	// Step 2: Move the enemy paddle to the previously decided location
	g.enemy.move()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// draw dashed line in the middle (dimensions 10x60 per dash and 40px space between dashes)
	for i := 0; i < screenHeight; i += 100 {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(halfGameScreenWidth), float64(i))
		vector.StrokeLine(screen, float32(halfGameScreenWidth), float32(i), float32(halfGameScreenWidth), float32(i+60), 10, color.White)
	}

	g.ball.Draw(screen)
	g.player.paddle.Draw(screen)
	g.enemy.paddle.Draw(screen)

	// draw score
	text.Draw(screen, fmt.Sprintf("%d", enemyScoreCount), scoreDisplayFont, halfGameScreenWidth-360, 120, color.White)
	text.Draw(screen, fmt.Sprintf("%d", playerScoreCount), scoreDisplayFont, halfGameScreenWidth+360-75, 120, color.White)

	if g.state == paused {
		text.Draw(screen, "PAUSED", scoreDisplayFont, halfGameScreenWidth-100, halfGameScreenHeight-100, color.White)
	}

	if g.state == gameOver {
		if playerScoreCount > enemyScoreCount {
			text.Draw(screen, "WINNER", resultDisplayFont, halfGameScreenWidth+450, halfGameScreenHeight, color.White)
			text.Draw(screen, "LOSER", resultDisplayFont, halfGameScreenWidth-450, halfGameScreenHeight, color.White)
		} else {
			text.Draw(screen, "WINNER", resultDisplayFont, halfGameScreenWidth-450, halfGameScreenHeight, color.White)
			text.Draw(screen, "LOSER", resultDisplayFont, halfGameScreenWidth+350, halfGameScreenHeight, color.White)
		}
	}
}

func (g *Game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

func init() {
	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	scoreDisplayFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    76,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	resultDisplayFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    18,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Pong")
	game := newGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
