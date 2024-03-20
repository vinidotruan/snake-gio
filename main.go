package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 800
	speed        = 10
	snakeSize    = 20
	foodRadius   = snakeSize / 2
)

var passedtime float32

type Snake struct {
	Head  rl.Vector2
	Body  []rl.Vector2
	Speed rl.Vector2
}

type Food struct {
	Position rl.Vector2
	Status   bool
}
type Game struct {
	GameOver bool
	Score    int
	Snake    Snake
	Foods    []Food
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

	rl.DrawText(fmt.Sprint(g.Score), 10, screenHeight/2, 20, rl.White)
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

		g.Snake.Head.X += g.Snake.Speed.X
		g.Snake.Head.Y += g.Snake.Speed.Y

		if (g.Snake.Head.X+speed > screenWidth || g.Snake.Head.X < 0) ||
			(g.Snake.Head.Y+speed > screenHeight || g.Snake.Head.Y < 0) {
			g.GameOver = true
		}

		passedtime += rl.GetFrameTime()
		if int32(passedtime)%5 == 0 && g.Frames%60 == 0 {
			food := Food{Position: getRandomPosition(), Status: true}
			g.Foods = append(g.Foods, food)
		}

		if len(g.Foods)-1 >= 0 && g.Foods[len(g.Foods)-1].Status {
			rl.DrawCircle(int32(g.Foods[len(g.Foods)-1].Position.X), int32(g.Foods[len(g.Foods)-1].Position.Y), foodRadius, rl.RayWhite)
		}

		if rl.CheckCollisionCircleRec(g.Foods[len(g.Foods)-1].Position, foodRadius, rl.NewRectangle(g.Snake.Head.X, g.Snake.Head.Y, snakeSize, snakeSize)) && g.Foods[len(g.Foods)-1].Status {
			g.Foods[len(g.Foods)-1].Status = false
			g.Score += 5
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
	rl.DrawRectangle(int32(g.Snake.Head.X), int32(g.Snake.Head.Y), snakeSize, snakeSize, rl.Red)

	rl.EndDrawing()
}

func (g *Game) Init() {
	g.Snake = Snake{Head: rl.NewVector2(screenWidth/2, screenHeight/2), Body: []rl.Vector2{}}
}

func getRandomPosition() rl.Vector2 {
	return rl.NewVector2(float32(rl.GetRandomValue(0, screenWidth-20)), float32(rl.GetRandomValue(0, screenHeight-20)))
}
