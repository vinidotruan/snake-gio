package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"strings"
)

const (
	screenWidth  = 800
	screenHeight = 800
	speed        = 20
	snakeSize    = 20
	foodRadius   = snakeSize / 2
)

var (
	passedtime float32
	direction  string
	bodyPiece  rl.Vector2
)

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
	rl.DrawText(fmt.Sprint(direction), screenHeight/2, 0, 20, rl.White)
	rl.DrawText(fmt.Sprint(g.Score), screenHeight/4, 0, 20, rl.White)
	if !g.GameOver {
		if rl.IsKeyPressed(rl.KeyP) {
			g.Pause()
		}
		if (rl.IsKeyPressed(rl.KeyRight) || rl.IsKeyPressed(rl.KeyD)) && strings.Compare(direction, "left") != 0 {
			g.Snake.Speed = rl.Vector2{X: speed, Y: 0}
			direction = "right"
		}

		if (rl.IsKeyPressed(rl.KeyLeft) || rl.IsKeyPressed(rl.KeyA)) && strings.Compare(direction, "right") != 0 {
			g.Snake.Speed = rl.Vector2{X: -speed, Y: 0}
			direction = "left"
		}

		if (rl.IsKeyPressed(rl.KeyUp) || rl.IsKeyPressed(rl.KeyW)) && strings.Compare(direction, "down") != 0 {
			g.Snake.Speed = rl.Vector2{X: 0, Y: -speed}
			direction = "up"
		}

		if (rl.IsKeyPressed(rl.KeyDown) || rl.IsKeyPressed(rl.KeyS)) && strings.Compare(direction, "up") != 0 {
			g.Snake.Speed = rl.Vector2{X: 0, Y: +speed}
			direction = "down"
		}

		if g.Frames%5 == 0 {
			// position of body pieces
			for i := len(g.Snake.Body) - 1; i > 0; i-- {
				g.Snake.Body[i].X = g.Snake.Body[i-1].X
				g.Snake.Body[i].Y = g.Snake.Body[i-1].Y
			}

			if len(g.Snake.Body) > 0 {
				g.Snake.Body[0] = g.Snake.Head
			}
			// move snake head
			g.Snake.Head.X += g.Snake.Speed.X
			g.Snake.Head.Y += g.Snake.Speed.Y

		}
		// draw body
		if len(g.Snake.Body) > 0 {
			for i := 0; i < len(g.Snake.Body); i++ {
				rl.DrawRectangle(
					int32(g.Snake.Body[i].X),
					int32(g.Snake.Body[i].Y),
					snakeSize, snakeSize, rl.Green)
			}
		}
		// wall collision
		if (g.Snake.Head.X+speed > screenWidth || g.Snake.Head.X < 0) ||
			(g.Snake.Head.Y+speed > screenHeight || g.Snake.Head.Y < 0) {
			g.Pause()
			// g.GameOver = true
		}

		// draw fruit
		passedtime += rl.GetFrameTime()
		if int32(passedtime)%5 == 0 && g.Frames%60 == 0 {
			food := Food{Position: getRandomPosition(), Status: true}
			g.Foods = append(g.Foods, food)
		}

		if len(g.Foods)-1 >= 0 && g.Foods[len(g.Foods)-1].Status {
			rl.DrawCircle(int32(g.Foods[len(g.Foods)-1].Position.X), int32(g.Foods[len(g.Foods)-1].Position.Y), foodRadius, rl.RayWhite)
		}

		// head touch body
		for j := 0; j < len(g.Snake.Body); j++ {
			if rl.CheckCollisionRecs(rl.NewRectangle(g.Snake.Head.X, g.Snake.Head.Y, snakeSize, snakeSize), rl.NewRectangle(g.Snake.Body[j].X, g.Snake.Body[j].Y, snakeSize, snakeSize)) {
				g.GameOver = true
			}
		}

		// was fruit ate
		if rl.CheckCollisionCircleRec(g.Foods[len(g.Foods)-1].Position, foodRadius, rl.NewRectangle(g.Snake.Head.X, g.Snake.Head.Y, snakeSize, snakeSize)) && g.Foods[len(g.Foods)-1].Status {
			g.Foods[len(g.Foods)-1].Status = false
			g.Score += 5

			// get position of last body piece
			x, y := func() (float32, float32) {
				if len(g.Snake.Body) > 0 {
					return g.Snake.Body[len(g.Snake.Body)-1].X, g.Snake.Body[len(g.Snake.Body)-1].Y
				}
				return g.Snake.Head.X, g.Snake.Head.Y
			}()

			if strings.Compare(direction, "right") == 0 {
				x -= snakeSize
			} else if strings.Compare(direction, "left") == 0 {
				x += snakeSize
			} else if strings.Compare(direction, "up") == 0 {
				y += snakeSize
			} else {
				y -= snakeSize
			}

			bodyPiece = rl.NewVector2(x, y)
			rl.DrawRectangle(int32(x), int32(y), snakeSize, snakeSize, rl.Green)

			g.Snake.Body = append(g.Snake.Body, bodyPiece)

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
	fmt.Println(g.Snake.Body)
}

func (g *Game) Pause() {
	g.Snake.Speed = rl.NewVector2(0, 0)
}

func getRandomPosition() rl.Vector2 {
	return rl.NewVector2(float32(rl.GetRandomValue(0, screenWidth-20)), float32(rl.GetRandomValue(0, screenHeight-20)))
}
