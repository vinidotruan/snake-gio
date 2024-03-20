package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 800
	speed        = 10
)

type Snake struct {
	Head  rl.Vector2
	Body  []rl.Vector2
	Speed rl.Vector2
}
type Game struct {
	GameOver bool
	Score    int
	Snake    Snake
	Food     rl.Vector2
	Frames   int32
}

func main() {
	game := Game{GameOver: false}
	game.Init()

	rl.InitWindow(screenWidth, screenHeight, "raylib [core] example - basic window")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		game.Update()
		game.Draw()
	}
}

func (g *Game) Update() {
	if !g.GameOver {
		if rl.IsKeyPressed(rl.KeyRight) || rl.IsKeyPressed(rl.KeyD) {
			g.Snake.Speed = rl.Vector2{X: speed, Y: 0}
		}

		if rl.IsKeyPressed(rl.KeyLeft) || rl.IsKeyPressed(rl.KeyA) {
			g.Snake.Speed = rl.Vector2{X: -speed, Y: 0}
		}

		if rl.IsKeyPressed(rl.KeyUp) || rl.IsKeyPressed(rl.KeyW) {
			g.Snake.Speed = rl.Vector2{X: 0, Y: -speed}
		}

		if rl.IsKeyPressed(rl.KeyDown) || rl.IsKeyPressed(rl.KeyS) {
			g.Snake.Speed = rl.Vector2{X: 0, Y: +speed}
		}

		if g.Frames%5 == 0 {
			g.Snake.Head.X += g.Snake.Speed.X
			g.Snake.Head.Y += g.Snake.Speed.Y

		}

		if (g.Snake.Head.X+speed > screenWidth || g.Snake.Head.X < 0) ||
			(g.Snake.Head.Y+speed > screenHeight || g.Snake.Head.Y < 0) {
			g.GameOver = true
		}
		g.Frames++

	} else {
		rl.DrawText("Game Over", screenWidth/2, screenHeight/2, 20, rl.White)

		if rl.IsKeyPressed(rl.KeyR) {
			g.Init()
			g.GameOver = false
		}
	}

}

func (g *Game) Draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.DrawRectangle(int32(g.Snake.Head.X), int32(g.Snake.Head.Y), 20, 20, rl.Red)
	rl.EndDrawing()
}

func (g *Game) Init() {
	g.Snake = Snake{Head: rl.NewVector2(screenWidth/2, screenHeight/2), Body: []rl.Vector2{}}
}
