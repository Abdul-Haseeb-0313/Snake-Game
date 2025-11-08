package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

// ----------------------------
// Utility Functions
// ----------------------------
func randomInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

// ----------------------------
// Snake Definition
// ----------------------------
type Snake struct {
	parts []struct{ x, y int }
}

func (s *Snake) move(x, y int) {
	for i := len(s.parts) - 1; i > 0; i-- {
		s.parts[i].x = s.parts[i-1].x
		s.parts[i].y = s.parts[i-1].y
	}
	s.parts[0].x += x
	s.parts[0].y += y
}

func (s *Snake) addPart() {
	n := len(s.parts)
	var x, y int
	if n == 0 {
		x = randomInt(1, 19) * 20
		y = randomInt(1, 19) * 20
	} else {
		x = s.parts[n-1].x - 20
		y = s.parts[n-1].y
	}
	s.parts = append(s.parts, struct{ x, y int }{x, y})
}

func (s *Snake) hitSelf() bool {
	head := s.parts[0]
	for i := 1; i < len(s.parts); i++ {
		if s.parts[i].x == head.x && s.parts[i].y == head.y {
			return true
		}
	}
	return false
}

type Game struct {
	snake     *Snake
	food      struct{ x, y int }
	score     int
	highScore int
	gameOver  bool
}

var movement = "left"
var tick = 0


func (g *Game) Update() error {
	// When game is over, wait for Enter to restart
	if g.gameOver {
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			g.reset()
		}
		return nil
	}

	tick++
	if tick%10 != 0 { // Slowed down for better gameplay
		return nil
	}

	// Eat food
	if g.snake.parts[0].x == g.food.x && g.snake.parts[0].y == g.food.y {
		x := randomInt(1, 19) * 20
		y := randomInt(1, 19) * 20
		g.food.x = x
		g.food.y = y
		g.snake.addPart()
		g.score++
		if g.score > g.highScore {
			g.highScore = g.score
		}
	}


	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) && movement != "left" {
		movement = "right"
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && movement != "right" {
		movement = "left"
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && movement != "down" {
		movement = "up"
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && movement != "up" {
		movement = "down"
	}

	switch movement {
	case "left":
		g.snake.move(-20, 0)
	case "right":
		g.snake.move(20, 0)
	case "up":
		g.snake.move(0, -20)
	case "down":
		g.snake.move(0, 20)
	}

	head := &g.snake.parts[0]
	if head.x < 0 {
		head.x = 380
	} else if head.x >= 400 {
		head.x = 0
	}
	if head.y < 0 {
		head.y = 380
	} else if head.y >= 400 {
		head.y = 0
	}

	
	if g.snake.hitSelf() {
		g.gameOver = true
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Dark background
	screen.Fill(color.RGBA{20, 20, 35, 255})
	
	// Draw border around gameplay area
	border := ebiten.NewImage(420, 420)
	border.Fill(color.RGBA{70, 70, 100, 255})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(40, 60)
	screen.DrawImage(border, op)
	
	// Game area background
	gameArea := ebiten.NewImage(400, 400)
	gameArea.Fill(color.RGBA{30, 30, 50, 255})
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(50, 70)
	screen.DrawImage(gameArea, op)

	if g.gameOver {
		// Semi-transparent overlay
		overlay := ebiten.NewImage(500, 500)
		overlay.Fill(color.RGBA{0, 0, 0, 200})
		screen.DrawImage(overlay, &ebiten.DrawImageOptions{})
		
		// Game Over message with better styling
		gameOverText := "GAME OVER"
		scoreText := fmt.Sprintf("Score: %d", g.score)
		highScoreText := fmt.Sprintf("High Score: %d", g.highScore)
		restartText := "Press ENTER to Restart"
		
		// Use text drawer for better control
		text.Draw(screen, gameOverText, basicfont.Face7x13, 180, 200, color.RGBA{255, 50, 50, 255})
		text.Draw(screen, scoreText, basicfont.Face7x13, 210, 230, color.White)
		text.Draw(screen, highScoreText, basicfont.Face7x13, 190, 250, color.White)
		text.Draw(screen, restartText, basicfont.Face7x13, 160, 280, color.RGBA{100, 255, 100, 255})
		
		return
	}

	// Draw the snake
	for i, part := range g.snake.parts {
		rect := ebiten.NewImage(18, 18) // Slightly smaller for gap between segments
		if i == 0 {
			rect.Fill(color.RGBA{50, 200, 50, 255}) // Head = bright green
		} else {
			// Gradient color for body
			greenVal := uint8(150 - i*2)
			if greenVal < 50 {
				greenVal = 50
			}
			rect.Fill(color.RGBA{50, greenVal, 50, 255})
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(part.x)+51, float64(part.y)+71) // Adjusted for border
		screen.DrawImage(rect, op)
	}

	
	foodImg := ebiten.NewImage(18, 18)
	foodImg.Fill(color.RGBA{255, 50, 50, 255}) // Red food
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.food.x)+51, float64(g.food.y)+71) // Adjusted for border
	screen.DrawImage(foodImg, op)

	
	scoreBox := ebiten.NewImage(500, 30)
	scoreBox.Fill(color.RGBA{40, 40, 60, 255})
	screen.DrawImage(scoreBox, &ebiten.DrawImageOptions{})
	
	scoreText := fmt.Sprintf("Score: %d", g.score)
	highScoreText := fmt.Sprintf("High Score: %d", g.highScore)
	
	text.Draw(screen, "SNAKE GAME", basicfont.Face7x13, 200, 20, color.RGBA{100, 255, 100, 255})
	text.Draw(screen, scoreText, basicfont.Face7x13, 50, 45, color.White)
	text.Draw(screen, highScoreText, basicfont.Face7x13, 350, 45, color.White)
	
	// Controls hint
	controlsText := "Use Arrow Keys to Move"
	text.Draw(screen, controlsText, basicfont.Face7x13, 170, 500, color.RGBA{150, 150, 200, 255})
}

// ----------------------------
// Game Setup + Reset
// ----------------------------
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 500, 520 // Increased height for better layout
}

func (g *Game) reset() {
	g.snake = &Snake{parts: []struct{ x, y int }{{x: 200, y: 200}}}
	g.food.x = randomInt(1, 19) * 20
	g.food.y = randomInt(1, 19) * 20
	g.score = 0
	movement = "left"
	g.gameOver = false
}


func main() {
	rand.Seed(time.Now().UnixNano())

	game := &Game{
		snake: &Snake{parts: []struct{ x, y int }{{x: 200, y: 200}}},
		score: -1,
	}
	game.reset()

	ebiten.SetWindowSize(500, 520)
	ebiten.SetWindowTitle("Snake Game")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
