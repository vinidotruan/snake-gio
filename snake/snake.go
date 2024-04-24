package snake

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/vinidotruan/snake-go/utils"
	"github.com/vinidotruan/snake-go/utils/interfaces"
)

const (
	Size = 20
)

type Snake struct {
	head   rl.Rectangle
	bodies []interfaces.BodyInterface
	Speed  rl.Vector2
}

type Body struct {
	rectangle rl.Rectangle
}

func (s *Snake) Head() rl.Rectangle {
	return s.head
}

func (s *Snake) Bodies() []interfaces.BodyInterface {
	return s.bodies
}

func (s *Snake) DrawBodies(bodyCount int) {
	if bodyCount > 0 {
		for k := 0; k < len(s.Bodies()); k++ {
			rl.DrawRectangle(
				int32(s.Bodies()[k].Rectangle().X),
				int32(s.Bodies()[k].Rectangle().Y),
				Size,
				Size,
				utils.Green,
			)

		}
	}
}

func (s *Snake) VerifyBodyXHeadCollision() bool {
	for j := len(s.Bodies()) - 1; j > 0; j-- {
		if rl.CheckCollisionRecs(s.head, s.Bodies()[j].Rectangle()) {
			return true
		}
	}

	return false
}

func (s *Snake) VerifyWallCollisionValidation() bool {
	if (s.head.X >= utils.FinalX || s.head.X < utils.InitialX) ||
		(s.head.Y >= utils.FinalY || s.head.Y < utils.InitialY) {
		return true
	}

	return false
}

func (s *Snake) VerifyHeadXObstaclesCollision(obstacles []utils.Obstacle) bool {
	for _, obstacle := range obstacles {
		if rl.CheckCollisionRecs(s.head, obstacle.Shape) {
			return true
		}
	}
	return false
}

func (s *Snake) HasSnakeCollided(obstacles []utils.Obstacle) bool {
	return s.VerifyHeadXObstaclesCollision(obstacles) || s.VerifyBodyXHeadCollision() || s.VerifyWallCollisionValidation()
}

func (s *Snake) MoveBodies() {
	firstBody := s.Bodies()[0].Rectangle()
	for i := len(s.Bodies()) - 1; i > 0; i-- {
		b := s.Bodies()[i].Rectangle()
		b.X = s.Bodies()[i-1].Rectangle().X
		b.Y = s.Bodies()[i-1].Rectangle().Y
	}

	if len(s.Bodies()) > 0 {
		firstBody.X = s.head.X
		firstBody.Y = s.head.Y
	}
}

func (s *Snake) MoveHead() {
	s.head.X += s.Speed.X
	s.head.Y += s.Speed.Y
}

func (b Body) Rectangle() rl.Rectangle {
	return b.rectangle
}

func (b *Body) NewBody(rectangle rl.Rectangle) {
	b.rectangle = rectangle
}
