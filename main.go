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
	screenOffset = 80

	initialX  = screenOffset
	initialY  = screenOffset
	finalX    = screenWidth - screenOffset
	finalY    = screenHeight - screenOffset
	frameRate = 60
)

var (
	pastime     float32
	direction   string
	image       *rl.Image
	headTexture rl.Texture2D
	purple      = rl.NewColor(102, 102, 204, 255)
	green       = rl.NewColor(153, 204, 102, 255)
	greenDark   = rl.NewColor(102, 153, 153, 255)
	gray        = rl.NewColor(204, 204, 204, 255)
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

	rl.InitWindow(screenWidth, screenHeight, "dotsnake")
	defer rl.CloseWindow()
	rl.SetTargetFPS(frameRate)

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
		g.FoodCollision()

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
	rl.ClearBackground(rl.NewColor(40, 42, 54, 255))
	rl.DrawRectangle(int32(g.Snake.Head.X), int32(g.Snake.Head.Y), snakeSize, snakeSize, purple)

	DrawGrid()
	g.DrawBodies()
	g.DrawFruits()

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
}

func (g *Game) Movement() {
	if g.Frames%3 == 0 {
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
	if (g.Snake.Head.X+speed > finalX || g.Snake.Head.X-speed < initialX) ||
		(g.Snake.Head.Y+speed > finalY || g.Snake.Head.Y < initialY) {
		g.GameOver = true
	}

}

func (g *Game) SpawnFood() {
	position := getRandomPosition()
	food := Food{Shape: rl.NewRectangle(position.X, position.Y, snakeSize, snakeSize), Status: true}
	g.Foods = append(g.Foods, food)
}

func (g *Game) FoodCollision() {

	if len(g.Foods) == 0 {
		g.SpawnFood()
	}
	if rl.CheckCollisionRecs(g.Foods[len(g.Foods)-1].Shape, g.Snake.Head) && g.Foods[len(g.Foods)-1].Status {
		g.Foods[len(g.Foods)-1].Status = false
		g.Score += 5

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
		g.SpawnFood()

	}
}

func getRandomPosition() rl.Vector2 {
	x := float32(rl.GetRandomValue(initialX, finalX) / snakeSize * snakeSize)
	y := float32(rl.GetRandomValue(initialY, finalY) / snakeSize * snakeSize)

	return rl.NewVector2(x, y)
}

func DrawGrid() {

	rl.DrawLine(initialX, initialY, initialX, finalY, greenDark)
	rl.DrawLine(finalX, initialY, finalX, finalY, greenDark)

	rl.DrawLine(initialX, initialY, finalX, initialY, greenDark)
	rl.DrawLine(initialX, finalY, finalX, finalY, greenDark)
	color := rl.NewColor(204, 204, 204, 20)
	for x := int32(0); x < int32(screenWidth); x += snakeSize {
		rl.DrawLine(x, 0, x, int32(screenHeight), color)
	}

	for y := int32(0); y < int32(screenHeight); y += snakeSize {
		rl.DrawLine(0, y, int32(screenWidth), y, color)
	}
}

func (g *Game) DrawBodies() {
	if len(g.Snake.Bodies) > 0 {
		for k := 0; k < len(g.Snake.Bodies); k++ {
			rl.DrawRectangle(int32(g.Snake.Bodies[k].rectangle.X), int32(g.Snake.Bodies[k].rectangle.Y), snakeSize, snakeSize, green)

		}
	}
}

func (g *Game) DrawFruits() {
	if len(g.Foods)-1 >= 0 && g.Foods[len(g.Foods)-1].Status {
		fruit := rl.NewRectangle(g.Foods[len(g.Foods)-1].Shape.X, g.Foods[len(g.Foods)-1].Shape.Y, snakeSize, snakeSize)
		rl.DrawRectangle(int32(fruit.X), int32(fruit.Y), int32(fruit.Width), int32(fruit.Height), rl.White)
	}
}
