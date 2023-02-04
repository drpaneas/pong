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
	"time"
)

// Game constants (global)
const (
	gameScreenWidth      = 1280
	gameScreenHeight     = 720
	ballSpeed            = 15.0
	paddleSpeed          = 15.0
	halfGameScreenWidth  = gameScreenWidth / 2
	halfGameScreenHeight = gameScreenHeight / 2
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
	// Game objects, are all rectangles (yes even the ball, because this game is older than circles)
	ball, player, enemy = createGameObjects()

	// Velocities and speeds
	ballXVelocity, ballYVelocity              int
	playerPaddleVelocity, enemyPaddleVelocity = 0, paddleSpeed

	// Scoring and font
	playerScoreCount, enemyScoreCount   int
	scoreDisplayFont, resultDisplayFont font.Face
	volleyCount                         int

	// Game state
	currentGameState GameState = playing
)

// Game struct - required by Ebiten to run the game!
type Game struct{}

// Functions
// --------------------------------------------------------------------------------------------------------------- //

// createGameObjects creates the main game objects (ball, player, enemy)
// This function is called only once, at the start of the game or when the game is reset.
func createGameObjects() (*rect.Rectangle, *rect.Rectangle, *rect.Rectangle) {
	// Place the ball (square 20x20) in the center of the screen
	ball := rect.Rect(gameScreenWidth/2-10, gameScreenHeight/2-10, 20, 20)

	// Place the player and enemy paddles (rectangle 20x110) 70 pixels away from the left and right sides of the screen, respectively.
	enemy := rect.Rect(70, halfGameScreenHeight-110/2, 20, 110)
	player := rect.Rect(gameScreenWidth-70-20, halfGameScreenHeight-110/2, 20, 110)

	return ball, player, enemy
}

// startNewRound begins a new round of the game (should be called after a player or enemy scores).
// It places the ball back in the center of the screen and serves it to a random direction with a lower speed.
func startNewRound() {
	currentGameState = firstService
	volleyCount = 0 // reset the volley count
	ball.Center(halfGameScreenWidth, halfGameScreenHeight)

	// Serve the ball to a random side, with lower speed,
	reduceBallSpeed()
}

// updateBallPosition updates the position of the ball based on its current velocity.
func updateBallPosition() {
	ball.X += ballXVelocity
	ball.Y += ballYVelocity
}

// reduceBallSpeed reduces the ball speed by a factor of 3.
// This is used when the ball is served to a player for the first time.
func reduceBallSpeed() {
	ballXVelocity = ballSpeed * randomChoice(-1, 1) / 3
	ballYVelocity = ballSpeed * randomChoice(-1, 1) / 3
}

// increaseBallSpeed increases the ball speed to its maximum value
func increaseBallSpeed() {
	signX := 1
	if ballXVelocity < 0 {
		signX = -1
	}
	signY := 1
	if ballYVelocity < 0 {
		signY = -1
	}
	ballXVelocity = signX * ballSpeed
	ballYVelocity = signY * ballSpeed
}

// ballMovement updates the position of the ball in the game and handles various events that can occur in relation to the ball.
// The function checks if the ball goes out of the screen, scores, hits player's or enemy's paddle, hits the corners of the screen,
// or changes speed, and performs appropriate actions such as reversing direction, resetting the ball, bouncing back, playing sound,
// increasing the volley count, normalizing speed, etc. The purpose of this function is to update the state of the ball in the game.
func ballMovement() {
	updateBallPosition() // update ball position

	// Check if ball goes out of screen
	if ball.Top() < 0 || ball.Bottom() > gameScreenHeight {
		// if ball is below screen, set bottom to screen height
		if ball.Bottom() >= gameScreenHeight {
			ball.Bottom(gameScreenHeight)
		} else {
			// if ball is above screen, set top to 0
			ball.Top(0)
		}
		ballYVelocity *= -1 // reverse Y axis velocity
	}

	// Check if player scores
	if ball.Left() <= 0 {
		playerScoreCount++
		if playerScoreCount == 10 {
			currentGameState = gameOver
		} else {
			startNewRound()
		}
	}

	// Check if enemy scores
	if ball.Right() >= gameScreenWidth {
		enemyScoreCount++
		if enemyScoreCount == 10 {
			currentGameState = gameOver
		} else {
			startNewRound()
		}
	}

	// Check if ball hits player paddle
	if ball.Overlaps(player) {
		volleyCount++
		ball.Right(player.Left()) // move ball so it touches paddle
		bouncePlayer()
		currentGameState = playing
	}

	// Check if ball hits enemy paddle
	if ball.Overlaps(enemy) {
		volleyCount++
		ball.Left(enemy.Right()) // move ball so it touches paddle
		bounceEnemy()
		currentGameState = playing
	}

	// Avoid ball getting stuck in screen corners
	noCornerStuck()

	// Normalize ball speed if game is playing
	if currentGameState == playing {
		normalizeBallSpeed()
	}
}

// noCornerStuck checks if the ball is in any of the four corners and, if so,
// moves the ball to the edge of the screen and reverses the speed. If the ball
// is not in any of the corners, the function returns early without doing any changes.
func noCornerStuck() {
	var x, y int
	var newX, newY int

	// Get the ball's current position
	if x, y = ball.TopLeft(); x <= 0 && y <= 0 {
		newX, newY = 0, 0
	} else if x, y = ball.TopRight(); x >= gameScreenWidth && y <= 0 {
		newX, newY = gameScreenWidth, 0
	} else if x, y = ball.BottomLeft(); x <= 0 && y >= gameScreenHeight {
		newX, newY = 0, gameScreenHeight
	} else if x, y = ball.BottomRight(); x >= gameScreenWidth && y >= gameScreenHeight {
		newX, newY = gameScreenWidth, gameScreenHeight
	} else {
		return // No corner collision detected, return early
	}

	// Move the ball back so it touches the screen but doesn't overlap
	ball.TopLeft(newX, newY)

	// Reverse the X and Y axis speed
	ballXVelocity *= -1
	ballYVelocity *= -1
}

func normalizeBallSpeed() {
	// Calculate the total speed of the ball in pixels per frame
	speed := math.Sqrt(math.Pow(float64(ballXVelocity), 2) + math.Pow(float64(ballYVelocity), 2))

	// Normalize the ball speed if it's larger than desired
	if speed > ballSpeed {
		// Adjust the X and Y components of velocity
		factor := ballSpeed / speed
		ballXVelocity = int(float64(ballXVelocity) * factor)
		ballYVelocity = int(float64(ballYVelocity) * factor)
	}
}

func bounceEnemy() {
	increaseBallSpeed()

	ballXVelocity *= -1

	part := enemy.Height / 8
	if volleyCount < 4 { // 0 1 2 3
		if ball.Top() < enemy.Top()+part {
			ballYVelocity = atAngle(-45)
		} else if ball.Top() < enemy.Top()+part*2 {
			ballYVelocity = atAngle(-30)
		} else if ball.Top() < enemy.Top()+part*3 {
			ballYVelocity = atAngle(-15)
		} else if ball.Top() < enemy.Top()+part*4 {
			ballYVelocity = atAngle(0)
		} else if ball.Top() < enemy.Top()+part*5 {
			ballYVelocity = atAngle(0)
		} else if ball.Top() < enemy.Top()+part*6 {
			ballYVelocity = atAngle(15)
		} else if ball.Top() < enemy.Top()+part*7 {
			ballYVelocity = atAngle(30)
		} else if ball.Top() < enemy.Top()+part*8 {
			ballYVelocity = atAngle(45)
		}
	}

	if volleyCount >= 4 && volleyCount < 8 { // 4 5 6 7
		if ball.Top() < enemy.Top()+part {
			ballYVelocity = atAngle(-60)
		} else if ball.Top() < enemy.Top()+part*2 {
			ballYVelocity = atAngle(-45)
		} else if ball.Top() < enemy.Top()+part*3 {
			ballYVelocity = atAngle(-30)
		} else if ball.Top() < enemy.Top()+part*4 {
			ballYVelocity = atAngle(-15)
		} else if ball.Top() < enemy.Top()+part*5 {
			ballYVelocity = atAngle(0)
		} else if ball.Top() < enemy.Top()+part*6 {
			ballYVelocity = atAngle(0)
		} else if ball.Top() < enemy.Top()+part*7 {
			ballYVelocity = atAngle(15)
		} else if ball.Top() < enemy.Top()+part*8 {
			ballYVelocity = atAngle(30)
		}
	}

	if volleyCount >= 8 {
		if ball.Top() < enemy.Top()+part {
			ballYVelocity = atAngle(-75)
		} else if ball.Top() < enemy.Top()+part*2 {
			ballYVelocity = atAngle(-60)
		} else if ball.Top() < enemy.Top()+part*3 {
			ballYVelocity = atAngle(-45)
		} else if ball.Top() < enemy.Top()+part*4 {
			ballYVelocity = atAngle(-30)
		} else if ball.Top() < enemy.Top()+part*5 {
			ballYVelocity = atAngle(-15)
		} else if ball.Top() < enemy.Top()+part*6 {
			ballYVelocity = atAngle(0)
		} else if ball.Top() < enemy.Top()+part*7 {
			ballYVelocity = atAngle(15)
		} else if ball.Top() < enemy.Top()+part*8 {
			ballYVelocity = atAngle(30)
		}
	}
}

func bouncePlayer() {
	increaseBallSpeed()

	// Chop the player in 8 parts and assign a different angle to each part
	ballXVelocity *= -1 // Reverse the ball direction on X axis
	part := player.Height / 8

	if volleyCount < 4 {
		if ball.Y < player.Y+part {
			ballYVelocity = atAngle(-135)
		} else if ball.Y < player.Y+part*2 {
			ballYVelocity = atAngle(-150)
		} else if ball.Y < player.Y+part*3 {
			ballYVelocity = atAngle(-165)
		} else if ball.Y < player.Y+part*4 {
			ballYVelocity = atAngle(180)
		} else if ball.Y < player.Y+part*5 {
			ballYVelocity = atAngle(180)
		} else if ball.Y < player.Y+part*6 {
			ballYVelocity = atAngle(165)
		} else if ball.Y < player.Y+part*7 {
			ballYVelocity = atAngle(150)
		} else if ball.Y < player.Y+part*8 {
			ballYVelocity = atAngle(135)
		}
	}

	if volleyCount >= 4 && volleyCount < 8 {
		if ball.Y < player.Y+part {
			ballYVelocity = atAngle(-135)
		} else if ball.Y < player.Y+part*2 {
			ballYVelocity = atAngle(-150)
		} else if ball.Y < player.Y+part*3 {
			ballYVelocity = atAngle(-165)
		} else if ball.Y < player.Y+part*4 {
			ballYVelocity = atAngle(180)
		} else if ball.Y < player.Y+part*5 {
			ballYVelocity = atAngle(180)
		} else if ball.Y < player.Y+part*6 {
			ballYVelocity = atAngle(165)
		} else if ball.Y < player.Y+part*7 {
			ballYVelocity = atAngle(150)
		} else if ball.Y < player.Y+part*8 {
			ballYVelocity = atAngle(135)
		}
	}

	if volleyCount >= 8 {
		if ball.Y < player.Y+part {
			ballYVelocity = atAngle(-135)
		} else if ball.Y < player.Y+part*2 {
			ballYVelocity = atAngle(-150)
		} else if ball.Y < player.Y+part*3 {
			ballYVelocity = atAngle(-165)
		} else if ball.Y < player.Y+part*4 {
			ballYVelocity = atAngle(180)
		} else if ball.Y < player.Y+part*5 {
			ballYVelocity = atAngle(180)
		} else if ball.Y < player.Y+part*6 {
			ballYVelocity = atAngle(165)
		} else if ball.Y < player.Y+part*7 {
			ballYVelocity = atAngle(150)
		} else if ball.Y < player.Y+part*8 {
			ballYVelocity = atAngle(135)
		}
	}

}

func atAngle(angle float64) int {
	// Convert the angle to radians
	radians := angle * math.Pi / 180
	// Calculate the new speed for Y axis
	return int(math.Round(math.Tan(radians) * float64(ballXVelocity)))
}

func playerMovement() {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		playerPaddleVelocity -= paddleSpeed
	} else if inpututil.IsKeyJustReleased(ebiten.KeyArrowUp) {
		playerPaddleVelocity += paddleSpeed
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		playerPaddleVelocity += paddleSpeed
	} else if inpututil.IsKeyJustReleased(ebiten.KeyArrowDown) {
		playerPaddleVelocity -= paddleSpeed
	}

	player.Y += playerPaddleVelocity

	// Check if the player is out of the screen
	if player.Top() < 0 {
		player.Top(0)
	}

	if player.Bottom() > gameScreenHeight {
		player.Bottom(gameScreenHeight)
	}
}

func enemyMovement() {
	enemyAI()

	// Check if the enemy is out of the screen
	if enemy.Top() < 0 {
		enemy.Top(0)
	}

	if enemy.Bottom() > gameScreenHeight {
		enemy.Bottom(gameScreenHeight)
	}
}

func enemyAI() {
	//// Enemy follows the ball on Y axis
	var tmp int
	if rand.Float64() > 0.1 {
		tmp = paddleSpeed * -1
	} else {
		tmp = paddleSpeed
	}

	// if half of the enemy height is below the center of ball, move down
	if enemy.CenterY() < ball.CenterY() {
		enemy.Y += int(enemyPaddleVelocity) - tmp
	}

	// if half of the enemy height is above the center of the ball, move up
	if enemy.CenterY() > ball.CenterY() {
		enemy.Y -= int(enemyPaddleVelocity) - tmp
	}
}

// function that returns a random number between min and max
func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

// Function that returns randomly either a or b. If a and b are equal, it returns a.
func randomChoice(a, b int) int {
	if a == b {
		return a
	}
	if rand.Intn(2) == 0 {
		return a
	}
	return b
}

func (g *Game) Update() error {
	if currentGameState == paused {
		return nil
	}
	if currentGameState == gameOver {
		return nil
	}
	ballMovement()
	playerMovement()
	enemyMovement()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// draw ball
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(ball.X), float64(ball.Y))
	vector.DrawFilledRect(screen, float32(ball.X), float32(ball.Y), float32(ball.Width), float32(ball.Height), color.White)

	// draw dashed line in the middle (dimensions 10x60 per dash and 40px space between dashes)
	for i := 0; i < gameScreenHeight; i += 100 {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(halfGameScreenWidth), float64(i))
		vector.StrokeLine(screen, float32(halfGameScreenWidth), float32(i), float32(halfGameScreenWidth), float32(i+60), 10, color.White)
	}

	// draw player
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(player.X), float64(player.X))
	vector.DrawFilledRect(screen, float32(player.X), float32(player.Y), float32(player.Width), float32(player.Height), color.White)

	// draw enemy
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(enemy.X), float64(enemy.Y))
	vector.DrawFilledRect(screen, float32(enemy.X), float32(enemy.Y), float32(enemy.Width), float32(enemy.Height), color.White)

	// draw score
	text.Draw(screen, fmt.Sprintf("%d", enemyScoreCount), scoreDisplayFont, halfGameScreenWidth-360, 120, color.White)
	text.Draw(screen, fmt.Sprintf("%d", playerScoreCount), scoreDisplayFont, halfGameScreenWidth+360-75, 120, color.White)

	if currentGameState == paused {
		text.Draw(screen, "PAUSED", scoreDisplayFont, halfGameScreenWidth-100, halfGameScreenHeight-100, color.White)
	}

	if currentGameState == gameOver {
		if playerScoreCount > enemyScoreCount {
			text.Draw(screen, "WINNER", resultDisplayFont, halfGameScreenWidth+450, halfGameScreenHeight, color.White)
			text.Draw(screen, "LOSER", resultDisplayFont, halfGameScreenWidth-450, halfGameScreenHeight, color.White)
		} else {
			text.Draw(screen, "WINNER", resultDisplayFont, halfGameScreenWidth-450, halfGameScreenHeight, color.White)
			text.Draw(screen, "LOSER", resultDisplayFont, halfGameScreenWidth+350, halfGameScreenHeight, color.White)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return gameScreenWidth, gameScreenHeight
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

	startNewRound() // Initialize ball position and speed
}

func main() {
	ebiten.SetWindowSize(gameScreenWidth, gameScreenHeight)
	ebiten.SetWindowTitle("Pong")

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
