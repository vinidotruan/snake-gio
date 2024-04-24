package food

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/vinidotruan/snake-go/snake"
	"github.com/vinidotruan/snake-go/utils"
	"github.com/vinidotruan/snake-go/utils/interfaces"
	"strings"
)

type Food struct {
	shape  rl.Rectangle
	point  int32
	status bool
}

func (f Food) Shape() rl.Rectangle {
	return f.shape
}

func (f Food) Point() int32 {
	return f.point
}

func (f Food) Status() bool {
	return f.status
}

func (f *Food) ChangeStatus(newStatus bool) {
	f.status = newStatus
}

func DrawFruits(foods []interfaces.FoodInterface) {
	if len(foods)-1 >= 0 && foods[len(foods)-1].Status() == true {
		fruit := rl.NewRectangle(foods[len(foods)-1].Shape().X, foods[len(foods)-1].Shape().Y, snake.Size, snake.Size)
		rl.DrawRectangle(int32(fruit.X), int32(fruit.Y), int32(fruit.Width), int32(fruit.Height), rl.White)
	}
}

func VerifyFoodCollision(g interfaces.GameInterface, direction string) {
	s := g.Snake()
	b := s.Bodies()
	if len(g.Foods()) == 0 {
		SpawnFood(&s, g)
	}
	if rl.CheckCollisionRecs(g.Foods()[len(g.Foods())-1].Shape(), s.Head()) &&
		g.Foods()[len(g.Foods())-1].Status() {
		g.Foods()[len(g.Foods())-1].ChangeStatus(false)
		g.IncreaseScore(g.Foods()[len(g.Foods())-1].Point())

		x, y := func() (float32, float32) {
			if len(b) > 0 {
				lastBodyPiece := b[len(b)-1]
				return lastBodyPiece.Rectangle().X, lastBodyPiece.Rectangle().Y
			}
			return s.Head().X, s.Head().Y
		}()

		if strings.Compare(direction, "right") == 0 {
			x -= snake.Size
		} else if strings.Compare(direction, "left") == 0 {
			x += snake.Size
		} else if strings.Compare(direction, "up") == 0 {
			y += snake.Size
		} else {
			y -= snake.Size
		}

		rect := rl.NewRectangle(x, y, snake.Size, snake.Size)
		nb := snake.Body{}
		nb.NewBody(rect)
		b = append(b, nb)
		SpawnFood(&s, g)

	}
}

func SpawnFood(s interfaces.SnakeInterface, g interfaces.GameInterface) {
	position := utils.GetRandomPosition(utils.InitialX, utils.FinalX, utils.InitialY, utils.FinalY, snake.Size, s, g.Obstacles())
	newFood := Food{shape: rl.NewRectangle(position.X, position.Y, snake.Size, snake.Size), status: true, point: 5}
	foods := g.Foods()
	foods = append(foods, &newFood)
}
