package interfaces

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/vinidotruan/snake-go/snake"
	"github.com/vinidotruan/snake-go/utils"
)

type GameInterface interface {
	Update()
	Draw()
	Foods() []FoodInterface
	Snake() snake.Snake
	Obstacles() []utils.Obstacle
	Score() int32
	IncreaseScore(int32)
}

type FoodInterface interface {
	Status() bool
	Shape() rl.Rectangle
	Point() int32

	ChangeStatus(bool)
}

type SnakeInterface interface {
	DrawBodies(int)
	VerifyBodyXHeadCollision() bool
	VerifyWallCollisionValidation() bool
	VerifyHeadXObstaclesCollision([]utils.Obstacle) bool
	HasSnakeCollided([]utils.Obstacle) bool
	MoveBodies()

	Head() rl.Rectangle
	Bodies() []BodyInterface
}

type BodyInterface interface {
	Rectangle() rl.Rectangle
}
