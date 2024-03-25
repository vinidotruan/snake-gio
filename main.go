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
)

var (
	pastime          float32
	direction        string
	image            *rl.Image
	headTexture      rl.Texture2D
	bodyTexture      rl.Texture2D
	currentDirection float32 = 0
	targetDirection  float32
)

type Snake struct {
	Head   rl.Rectangle
	Bodies []Body
	Speed  rl.Vector2
}

type Body struct {
	rectangle rl.Rectangle
}

type Food struct {
	Shape  rl.Rectangle
	Status bool
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
	image = rl.LoadImage("assets/head_v1.png")
	headTexture = rl.LoadTextureFromImage(image)
	rl.UnloadImage(image)

	image = rl.LoadImage("assets/body_v1.png")
	bodyTexture = rl.LoadTextureFromImage(image)
	rl.UnloadImage(image)
	for !rl.WindowShouldClose() {
		game.Update()
		game.Draw()
	}
}

func (g *Game) Update() {
	rl.DrawText(fmt.Sprint(direction), screenHeight/2, 0, 20, rl.White)
	rl.DrawText(fmt.Sprint(g.Score), screenHeight/4, 0, 20, rl.White)

	if !g.GameOver {
		g.ControlsHandler()
		g.Movement()
		g.BodyXHeadCollision()
		g.WallCollisionValidation()

		// Get new fruit position
		pastime += rl.GetFrameTime()

		if int32(pastime)%5 == 0 && g.Frames%60 == 0 {
			position := getRandomPosition()
			food := Food{Shape: rl.NewRectangle(position.X, position.Y, snakeSize, snakeSize), Status: true}
			g.Foods = append(g.Foods, food)
		}

		// was fruit ate
		lastFruit := g.Foods[len(g.Foods)-1]
		if rl.CheckCollisionRecs(lastFruit.Shape, g.Snake.Head) && lastFruit.Status {
			fmt.Println("Teste")
			rl.DrawText("Fruta comida", 0, 0, 10, rl.White)
			lastFruit.Status = false
			g.Score += 5

			// get position of last body piece
			x, y := func() (float32, float32) {
				if len(g.Snake.Bodies) > 0 {
					lastBodyPiece := g.Snake.Bodies[len(g.Snake.Bodies)-1]
					return lastBodyPiece.rectangle.X, lastBodyPiece.rectangle.Y
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

			g.Snake.Bodies = append(g.Snake.Bodies, Body{rectangle: rl.NewRectangle(x, y, snakeSize, snakeSize)})
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
	rl.ClearBackground(rl.NewColor(40, 42, 54, 1))
	rl.DrawTexture(headTexture, int32(g.Snake.Head.X), int32(g.Snake.Head.Y), rl.White)

	// Draw body
	if len(g.Snake.Bodies) > 0 {
		for k := len(g.Snake.Bodies) - 1; k > 0; k-- {
			rl.DrawTextureRec(bodyTexture, g.Snake.Bodies[k].rectangle, rl.NewVector2(g.Snake.Bodies[k].rectangle.X, g.Snake.Bodies[k].rectangle.Y), rl.White)
		}
	}

	// Draw fruit
	if len(g.Foods)-1 >= 0 && g.Foods[len(g.Foods)-1].Status {
		fruit := rl.NewRectangle(g.Foods[len(g.Foods)-1].Shape.X, g.Foods[len(g.Foods)-1].Shape.Y, snakeSize, snakeSize)
		rl.DrawRectangle(int32(fruit.X), int32(fruit.Y), int32(fruit.Width), int32(fruit.Height), rl.White)
	}
	//rl.DrawTextureEx(headTexture, rl.NewVector2(g.Snake.Head.X, g.Snake.Head.Y), targetDirection, 10, rl.White)
	rl.EndDrawing()
}

func (g *Game) Init() {
	g.Snake = Snake{Head: rl.NewRectangle(screenWidth/2, screenHeight/2, snakeSize, snakeSize)}
}

func (g *Game) Pause() {
	g.Snake.Speed = rl.NewVector2(0, 0)
}

func (g *Game) ControlsHandler() {
	if rl.IsKeyPressed(rl.KeyP) {
		g.Pause()
	}
	if (rl.IsKeyPressed(rl.KeyRight) || rl.IsKeyPressed(rl.KeyD)) && strings.Compare(direction, "left") != 0 {
		g.Snake.Speed = rl.Vector2{X: speed, Y: 0}
		direction = "right"
		targetDirection = 90
	}

	if (rl.IsKeyPressed(rl.KeyLeft) || rl.IsKeyPressed(rl.KeyA)) && strings.Compare(direction, "right") != 0 {
		g.Snake.Speed = rl.Vector2{X: -speed, Y: 0}
		direction = "left"
		targetDirection = -90
	}

	if (rl.IsKeyPressed(rl.KeyUp) || rl.IsKeyPressed(rl.KeyW)) && strings.Compare(direction, "down") != 0 {
		g.Snake.Speed = rl.Vector2{X: 0, Y: -speed}
		direction = "up"
		targetDirection = 0
	}

	if (rl.IsKeyPressed(rl.KeyDown) || rl.IsKeyPressed(rl.KeyS)) && strings.Compare(direction, "up") != 0 {
		g.Snake.Speed = rl.Vector2{X: 0, Y: +speed}
		direction = "down"
		targetDirection = 180
	}
}

func (g *Game) Movement() {

	if g.Frames%5 == 0 {
		for i := len(g.Snake.Bodies) - 1; i > 0; i-- {
			g.Snake.Bodies[i].rectangle.X = g.Snake.Bodies[i-1].rectangle.X
			g.Snake.Bodies[i].rectangle.Y = g.Snake.Bodies[i-1].rectangle.Y
		}

		if len(g.Snake.Bodies) > 0 {
			g.Snake.Bodies[0].rectangle.X = g.Snake.Head.X
			g.Snake.Bodies[0].rectangle.Y = g.Snake.Head.Y
		}
		// move snake head
		g.Snake.Head.X += g.Snake.Speed.X
		g.Snake.Head.Y += g.Snake.Speed.Y
	}

}

func (g *Game) BodyXHeadCollision() {
	for j := len(g.Snake.Bodies) - 1; j > 0; j-- {
		if rl.CheckCollisionRecs(g.Snake.Head, g.Snake.Bodies[j].rectangle) {
			g.GameOver = true
		}
	}

}

func (g *Game) WallCollisionValidation() {
	if (g.Snake.Head.X+speed > screenWidth || g.Snake.Head.X < 0) || (g.Snake.Head.Y+speed > screenHeight || g.Snake.Head.Y < 0) {
		g.GameOver = true
	}

}

func getRandomPosition() rl.Vector2 {
	return rl.NewVector2(float32(rl.GetRandomValue(0, screenWidth-20)), float32(rl.GetRandomValue(0, screenHeight-20)))
}
