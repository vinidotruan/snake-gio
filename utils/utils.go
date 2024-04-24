package utils

import (
	"encoding/json"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/vinidotruan/snake-go/utils/interfaces"
	"os"
)

const (
	ScreenWidth  = 1080
	ScreenHeight = 1080
	Speed        = 20
	screenOffset = 80
	InitialX     = screenOffset
	InitialY     = screenOffset
	FinalX       = ScreenWidth - screenOffset
	FinalY       = ScreenHeight - screenOffset
)

var (
	MapList   []Map
	Purple    = rl.NewColor(102, 102, 204, 255)
	Green     = rl.NewColor(153, 204, 102, 255)
	GreenDark = rl.NewColor(102, 153, 153, 255)
	Gray      = rl.NewColor(204, 204, 204, 255)
	GrayDark  = rl.NewColor(40, 42, 54, 255)
)

type Obstacle struct {
	Shape  rl.Rectangle
	Status bool
}

type Map struct {
	Goal      int
	Time      int
	Obstacles []Obstacle
}

func LoadMap() {
	content, err := os.ReadFile("./maps/maps.json")

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(content, &MapList)
	if err != nil {
		panic(err)
	}
}

func GetRandomPosition(initialX int32, finalX int32, initialY int32, finalY int32, snakeSize int32, snake interfaces.SnakeInterface, obstacles []Obstacle) rl.Vector2 {
	x := float32(rl.GetRandomValue(initialX, finalX) / snakeSize * snakeSize)
	y := float32(rl.GetRandomValue(initialY, finalY) / snakeSize * snakeSize)

	if x == snake.Head().X && y == snake.Head().Y {
		return GetRandomPosition(initialX, finalX, initialY, finalY, snakeSize, snake, obstacles)
	}

	for _, bodies := range snake.Bodies() {
		if bodies.Rectangle().X == x && bodies.Rectangle().Y == y {
			return GetRandomPosition(initialX, finalX, initialY, finalY, snakeSize, snake, obstacles)
		}
	}

	for _, obstacle := range obstacles {
		if obstacle.Shape.X == x && obstacle.Shape.Y == y {
			return GetRandomPosition(initialX, finalX, initialY, finalY, snakeSize, snake, obstacles)
		}
	}

	return rl.NewVector2(x, y)
}
