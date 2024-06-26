package main

import (
	"encoding/json"
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"os"
	"strings"
)

const (
	screenWidth     = 1080
	screenHeight    = 1080
	speed           = 20
	snakeSize       = 20
	screenOffset    = 80
	initialX        = screenOffset
	initialY        = screenOffset
	finalX          = screenWidth - screenOffset
	finalY          = screenHeight - screenOffset
	frameRate       = 60
	pausedGUIWidth  = 400
	pausedGUIHeight = 240
	defaultFontSize = 20
)

var (
	time            int
	direction       string
	purple          = rl.NewColor(102, 102, 204, 255)
	green           = rl.NewColor(153, 204, 102, 255)
	greenDark       = rl.NewColor(102, 153, 153, 255)
	gray            = rl.NewColor(204, 204, 204, 255)
	grayDark        = rl.NewColor(40, 42, 54, 255)
	mapList         []Map
	currentMap      Map
	currentMapIndex = 0
	shouldMove      = true
	goingToNextMap  = false
	midPosition     = rl.NewVector2(screenWidth/2, screenHeight/2)
	game            Game
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
	Point  int32
	Status bool
}
type Game struct {
	GameOver  bool
	Paused    bool
	Score     int32
	Snake     Snake
	Foods     []Food
	Frames    int32
	Obstacles []Obstacle
	Gaming    bool
}

type Obstacle struct {
	Shape  rl.Rectangle
	Status bool
}

type Map struct {
	Goal      int
	Time      int
	Obstacles []Obstacle
}

func loadMap() {
	content, err := os.ReadFile("./maps/maps.json")

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(content, &mapList)
	if err != nil {
		panic(err)
	}
}

func getRandomPosition() rl.Vector2 {
	x := float32(rl.GetRandomValue(initialX, finalX) / snakeSize * snakeSize)
	y := float32(rl.GetRandomValue(initialY, finalY) / snakeSize * snakeSize)

	if x == game.Snake.Head.X && y == game.Snake.Head.Y {
		return getRandomPosition()
	}

	for _, bodies := range game.Snake.Bodies {
		if bodies.rectangle.X == x && bodies.rectangle.Y == y {
			return getRandomPosition()
		}
	}

	for _, obstacles := range game.Obstacles {
		if obstacles.Shape.X == x && obstacles.Shape.Y == y {
			return getRandomPosition()
		}
	}

	return rl.NewVector2(x, y)
}
func main() {
	loadMap()
	game := Game{GameOver: false, Gaming: false}

	rl.InitWindow(screenWidth, screenHeight, "dotsnake")
	defer rl.CloseWindow()
	rl.SetTargetFPS(frameRate)

	for !rl.WindowShouldClose() {
		game.Update()
		game.Draw()
	}
}

func (g *Game) Update() {
	g.ControlsHandler()

	if g.Gaming {
		if g.Frames%60 == 0 && !g.Paused && g.Gaming {
			time++
		}

		if !g.GameOver {

			if shouldMove {
				if g.Win() {
					g.NextPhase()
				}

				if g.Lose() {
					g.GameOver = true
				}

				g.Movement()
				g.BodyXHeadCollision()
				g.ObstaclesCollision()
				g.WallCollisionValidation()
				g.FoodCollision()
			}
			g.Frames++
		} else {
			rl.DrawText("Game Over", screenWidth/2, screenHeight/2, 20, rl.White)

			if rl.IsKeyPressed(rl.KeyR) {
				g.Reset()
			}
		}
	}

}

func (g *Game) Draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.NewColor(40, 42, 54, 255))

	if g.Gaming {
		rl.DrawText(fmt.Sprintf("Time: %d", time), screenHeight/2, 0, 20, rl.White)
		rl.DrawText(fmt.Sprint(g.Score), screenHeight/4, 0, 20, rl.White)
		rl.DrawText(fmt.Sprintf("Goal: %d", currentMap.Goal), screenHeight/4*3, 0, 20, rl.White)
		rl.DrawRectangle(int32(g.Snake.Head.X), int32(g.Snake.Head.Y), snakeSize, snakeSize, purple)

		DrawNewMapTimer()
		DrawGrid()

		g.DrawBodies()
		g.DrawFruits()
		g.DrawObstacles()
		g.DrawPausedGUI()
	} else {
		DrawInitialMenu()
	}

	rl.EndDrawing()
}

func (g *Game) Init() {
	g.Gaming = true
	g.Obstacles = []Obstacle{}
	g.Snake = Snake{Head: rl.NewRectangle(screenWidth/2, finalY-snakeSize, snakeSize, snakeSize), Speed: rl.Vector2{X: speed, Y: 0}}
	currentMap = mapList[currentMapIndex]
	g.LoadMapObstacles()
}

func (g *Game) Pause() {
	shouldMove = !shouldMove
	g.Paused = !g.Paused
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

	if rl.IsKeyPressed(rl.KeyEnter) {
		g.Init()
	}

	if rl.IsKeyPressed(rl.KeyEscape) {
		os.Exit(0)
	}
}

func (g *Game) Movement() {
	if shouldMove == true {
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

}

func (g *Game) BodyXHeadCollision() {
	for j := len(g.Snake.Bodies) - 1; j > 0; j-- {
		if rl.CheckCollisionRecs(g.Snake.Head, g.Snake.Bodies[j].rectangle) {
			g.GameOver = true
		}
	}

}

func (g *Game) WallCollisionValidation() {
	if (g.Snake.Head.X >= finalX || g.Snake.Head.X < initialX) ||
		(g.Snake.Head.Y >= finalY || g.Snake.Head.Y < initialY) {
		g.GameOver = true
	}

}

func (g *Game) ObstaclesCollision() {
	for _, obstacle := range g.Obstacles {
		if rl.CheckCollisionRecs(g.Snake.Head, obstacle.Shape) {
			g.GameOver = true
		}
	}
}

func (g *Game) SpawnFood() {
	position := getRandomPosition()
	food := Food{Shape: rl.NewRectangle(position.X, position.Y, snakeSize, snakeSize), Status: true, Point: 5}
	g.Foods = append(g.Foods, food)
}

func (g *Game) FoodCollision() {
	if len(g.Foods) == 0 {
		g.SpawnFood()
	}
	if rl.CheckCollisionRecs(g.Foods[len(g.Foods)-1].Shape, g.Snake.Head) && g.Foods[len(g.Foods)-1].Status {
		g.Foods[len(g.Foods)-1].Status = false
		g.Score += g.Foods[len(g.Foods)-1].Point

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

func (g *Game) DrawObstacles() {
	for l := 0; l < len(g.Obstacles); l++ {
		rl.DrawRectangleRec(g.Obstacles[l].Shape, gray)
	}
}

func (g *Game) Win() bool {
	return int(g.Score) == currentMap.Goal
}

func (g *Game) Lose() bool {
	return int(g.Score) < currentMap.Goal && time >= currentMap.Time
}

func (g *Game) Reset() {
	g.GameOver = false
	time = 0
	g.Snake.Head = rl.NewRectangle(screenWidth/2, finalY-snakeSize, snakeSize, snakeSize)
	g.Score = 0
	g.Snake.Bodies = []Body{}
	shouldMove = false
	goingToNextMap = true
	g.Seconds()
}

func (g *Game) NextPhase() {
	currentMapIndex++
	g.Reset()
	g.Init()
	DrawNewMapTimer()
}

func InitPhaseCounter() {
	if time < 3 {
		rl.DrawText(fmt.Sprintf("%d", time), int32(midPosition.X), int32(midPosition.Y), 100, purple)
	} else {
		shouldMove = true

	}

}

func DrawNewMapTimer() {
	if goingToNextMap {
		InitPhaseCounter()
	}
}

func (g *Game) Seconds() {
	if g.Frames%60 == 0 {
		time++
	}
}

func (g *Game) LoadMapObstacles() {
	if len(currentMap.Obstacles) > 0 {
		for _, obstacle := range currentMap.Obstacles {
			g.Obstacles = append(g.Obstacles, Obstacle{Shape: obstacle.Shape, Status: true})
		}
	}
}

func (g *Game) DrawPausedGUI() {
	if !g.Paused {
		return
	}

	content := "Continue: P"
	contentLength := rl.MeasureTextEx(rl.GetFontDefault(), content, defaultFontSize, 1)
	rectanglePosition := rl.NewVector2(midPosition.X-(pausedGUIWidth/2), midPosition.Y-(pausedGUIHeight/2))

	rl.DrawRectangle(0, 0, screenWidth, screenHeight, rl.NewColor(0, 0, 0, 100))
	rl.DrawRectangle(
		int32(rectanglePosition.X),
		int32(rectanglePosition.Y),
		pausedGUIWidth,
		pausedGUIHeight,
		purple,
	)

	rl.DrawText(
		content,
		int32(rectanglePosition.X)+pausedGUIWidth/2-(int32(contentLength.X)/2),
		int32(rectanglePosition.Y)+pausedGUIHeight/2-(int32(contentLength.Y)/2),
		defaultFontSize, gray)
}

func DrawInitialMenu() {
	rl.ClearBackground(grayDark)
	containerOffset := 10
	containerSize := 400
	subcontainerSize := containerSize - 20
	container := rl.NewRectangle(midPosition.X, midPosition.Y, float32(containerSize), float32(containerSize))
	subcontainer := rl.NewRectangle(container.X+float32(containerOffset), container.Y+float32(containerOffset), float32(subcontainerSize), float32(subcontainerSize))
	rl.DrawRectangleRec(container, purple)
	rl.DrawRectangleRec(subcontainer, grayDark)
	rl.DrawText("Start: Enter", int32(midPosition.X)-100, int32(midPosition.Y)+100, 20, rl.White)
	rl.DrawText("Quit: ESC", int32(midPosition.X)-100, int32(midPosition.Y)+120, 20, rl.White)
}
