package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

const (
	screenWidth  = 640
	screenHeight = 480
	ballSpeed    = 5
	paddleSpeed  = 6
)

type Object struct {
	X, Y, W, H int
}

type Paddle struct {
	Object
}

type Ball struct {
	Object
	dxdt int // x velocity per tick
	dydt int // y velocity per tick
}

type Game struct {
	paddle    Paddle
	ball      Ball
	score     int
	highscore int
}

func main() {
	ebiten.SetWindowTitle("Pong")
	ebiten.SetWindowSize(screenWidth, screenHeight)

	g := &Game{
		paddle: Paddle{
			Object: Object{
				X: 600,
				Y: 200,
				W: 15,
				H: 100,
			},
		},
		ball: Ball{
			Object: Object{
				X: 0,
				Y: 0,
				W: 15,
				H: 15,
			},
			dxdt: ballSpeed,
			dydt: ballSpeed,
		},
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen,
		float32(g.paddle.X), float32(g.paddle.Y),
		float32(g.paddle.W), float32(g.paddle.H),
		color.White, false,
	)

	vector.DrawFilledRect(screen,
		float32(g.ball.X), float32(g.ball.Y),
		float32(g.ball.W), float32(g.ball.H),
		color.White, false,
	)

	scoreStr := "Score: " + fmt.Sprint(g.score)
	text.Draw(screen, scoreStr, basicfont.Face7x13, 10, 20, color.White)

	highScoreStr := "High Score: " + fmt.Sprint(g.highscore)
	text.Draw(screen, highScoreStr, basicfont.Face7x13, 10, 40, color.White)
}

func (g *Game) Update() error {
	g.paddle.MoveOnKeyPress()
	g.ball.Move()
	g.CollideWithBall()
	g.CollideWithPaddle()
	return nil
}

func (p *Paddle) MoveOnKeyPress() {
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		p.Y += paddleSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		p.Y -= paddleSpeed
	}

	// Keep paddle within screen bounds
	if p.Y < 0 {
		p.Y = 0
	}
	if p.Y+p.H > screenHeight {
		p.Y = screenHeight - p.H
	}
}

func (b *Ball) Move() {
	b.X += b.dxdt
	b.Y += b.dydt
}

func (g *Game) Reset() {
	g.ball.X = 0
	g.ball.Y = 0
	g.ball.dxdt = ballSpeed
	g.ball.dydt = ballSpeed
	g.score = 0
}

func (g *Game) CollideWithBall() {
	if g.ball.X >= screenWidth {
		g.Reset()
	} else if g.ball.X <= 0 {
		g.ball.dxdt = ballSpeed
	}
	if g.ball.Y <= 0 {
		g.ball.dydt = ballSpeed
	} else if g.ball.Y+g.ball.H >= screenHeight {
		g.ball.dydt = -ballSpeed
	}
}

func (g *Game) CollideWithPaddle() {
	if g.ball.X+g.ball.W >= g.paddle.X &&
		g.ball.Y+g.ball.H >= g.paddle.Y &&
		g.ball.Y <= g.paddle.Y+g.paddle.H {
		g.ball.dxdt = -g.ball.dxdt
		g.score++
		if g.score > g.highscore {
			g.highscore = g.score
		}
	}
}
