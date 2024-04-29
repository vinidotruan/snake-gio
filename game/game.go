package game

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/vinidotruan/snake-go/food"
	"github.com/vinidotruan/snake-go/snake"
	"github.com/vinidotruan/snake-go/utils"
	"github.com/vinidotruan/snake-go/utils/interfaces"
	"os"
	"strings"
)

const (
	frameRate       = 60
	pausedGUIWidth  = 400
	pausedGUIHeight = 240
	defaultFontSize = 20
)

var (
	time      int
	direction string

	currentMap      utils.Map
	currentMapIndex = 0
	shouldMove      = true
	goingToNextMap  = false
	midPosition     = rl.NewVector2(utils.ScreenWidth/2, utils.ScreenHeight/2)
	game            interfaces.GameInterface
)

type Game struct {
	GameOver  bool
	Paused    bool
	score     int32
	snake     snake.Snake
	foods     []interfaces.FoodInterface
	Frames    int32
	obstacles []utils.Obstacle
	Gaming    bool
}

func (g *Game) Foods() []interfaces.FoodInterface {
	return g.foods
}

func (g *Game) Snake() snake.Snake {
	return g.snake
}

func (g *Game) Obstacles() []utils.Obstacle {
	return g.obstacles
}

func (g *Game) Score() int32 {
	return g.score
}

func (g *Game) IncreaseScore(score int32) {
	g.score += score
}

func Start() {
	utils.LoadMap()
	game = &Game{GameOver: false, Gaming: false}

	rl.InitWindow(utils.ScreenWidth, utils.ScreenHeight, "dotsnake")
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

				if g.Lose() ||
					g.snake.HasSnakeCollided(g.obstacles) {
					g.GameOver = true
				}

				g.Movement()

				food.VerifyFoodCollision(&*g, direction)
			}
			g.Frames++
		} else {
			rl.DrawText("Game Over", utils.ScreenWidth/2, utils.ScreenHeight/2, 20, rl.White)

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
		rl.DrawText(fmt.Sprintf("Time: %d", time), utils.ScreenHeight/2, 0, 20, rl.White)
		rl.DrawText(fmt.Sprint(g.score), utils.ScreenHeight/4, 0, 20, rl.White)
		rl.DrawText(fmt.Sprintf("Goal: %d", currentMap.Goal), utils.ScreenHeight/4*3, 0, 20, rl.White)
		rl.DrawRectangle(int32(g.snake.Head().X), int32(g.snake.Head().Y), snake.Size, snake.Size, utils.Purple)

		DrawNewMapTimer()
		DrawGrid()

		g.snake.DrawBodies(len(g.snake.Bodies()))
		food.DrawFruits(g.foods)
		g.DrawObstacles()
		g.DrawPausedGUI()
	} else {
		DrawInitialMenu()
	}

	rl.EndDrawing()
}

func (g *Game) Init() {
	g.Gaming = true
	g.obstacles = []utils.Obstacle{}
	g.snake = snake.Snake{ 
    Speed: rl.Vector2{X: utils.Speed, Y: 0},
  }
  g.snake.NewHead(rl.NewRectangle(utils.ScreenWidth/2, utils.FinalY-snake.Size, snake.Size, snake.Size))
	currentMap = utils.MapList[currentMapIndex]
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
		g.snake.Speed = rl.Vector2{X: utils.Speed, Y: 0}
		direction = "right"
	}

	if (rl.IsKeyPressed(rl.KeyLeft) || rl.IsKeyPressed(rl.KeyA)) && strings.Compare(direction, "right") != 0 {
		g.snake.Speed = rl.Vector2{X: -utils.Speed, Y: 0}
		direction = "left"
	}

	if (rl.IsKeyPressed(rl.KeyUp) || rl.IsKeyPressed(rl.KeyW)) && strings.Compare(direction, "down") != 0 {
		g.snake.Speed = rl.Vector2{X: 0, Y: -utils.Speed}
		direction = "up"
	}

	if (rl.IsKeyPressed(rl.KeyDown) || rl.IsKeyPressed(rl.KeyS)) && strings.Compare(direction, "up") != 0 {
		g.snake.Speed = rl.Vector2{X: 0, Y: +utils.Speed}
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
			g.snake.MoveBodies()
			g.snake.MoveHead()
		}
	}

}

func (g *Game) DrawObstacles() {
	for l := 0; l < len(g.obstacles); l++ {
		rl.DrawRectangleRec(g.obstacles[l].Shape, utils.Gray)
	}
}

func (g *Game) Win() bool {
	return int(g.score) == currentMap.Goal
}

func (g *Game) Lose() bool {
	return int(g.score) < currentMap.Goal && time >= currentMap.Time
}

func (g *Game) Reset() {
	g.GameOver = false
	time = 0
	g.snake.Reset(rl.NewRectangle(utils.ScreenWidth/2, utils.FinalY-snake.Size, snake.Size, snake.Size))
	g.score = 0
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

func (g *Game) DrawPausedGUI() {
	if !g.Paused {
		return
	}

	content := "Continue: P"
	contentLength := rl.MeasureTextEx(rl.GetFontDefault(), content, defaultFontSize, 1)
	rectanglePosition := rl.NewVector2(midPosition.X-(pausedGUIWidth/2), midPosition.Y-(pausedGUIHeight/2))

	rl.DrawRectangle(0, 0, utils.ScreenWidth, utils.ScreenHeight, rl.NewColor(0, 0, 0, 100))
	rl.DrawRectangle(
		int32(rectanglePosition.X),
		int32(rectanglePosition.Y),
		pausedGUIWidth,
		pausedGUIHeight,
		utils.Purple,
	)

	rl.DrawText(
		content,
		int32(rectanglePosition.X)+pausedGUIWidth/2-(int32(contentLength.X)/2),
		int32(rectanglePosition.Y)+pausedGUIHeight/2-(int32(contentLength.Y)/2),
		defaultFontSize, utils.Gray)
}

func (g *Game) Seconds() {
	if g.Frames%60 == 0 {
		time++
	}
}

func (g *Game) LoadMapObstacles() {
	if len(currentMap.Obstacles) > 0 {
		for _, obstacle := range currentMap.Obstacles {
			g.obstacles = append(g.obstacles, utils.Obstacle{Shape: obstacle.Shape, Status: true})
		}
	}
}

func DrawGrid() {
	rl.DrawLine(utils.InitialX, utils.InitialY, utils.InitialX, utils.FinalY, utils.GreenDark)
	rl.DrawLine(utils.FinalX, utils.InitialY, utils.FinalX, utils.FinalY, utils.GreenDark)

	rl.DrawLine(utils.InitialX, utils.InitialY, utils.FinalX, utils.InitialY, utils.GreenDark)
	rl.DrawLine(utils.InitialX, utils.FinalY, utils.FinalX, utils.FinalY, utils.GreenDark)
	color := rl.NewColor(204, 204, 204, 20)
	for x := int32(0); x < int32(utils.ScreenWidth); x += snake.Size {
		rl.DrawLine(x, 0, x, int32(utils.ScreenHeight), color)
	}

	for y := int32(0); y < int32(utils.ScreenHeight); y += snake.Size {
		rl.DrawLine(0, y, int32(utils.ScreenWidth), y, color)
	}
}

func InitPhaseCounter() {
	if time < 3 {
		rl.DrawText(fmt.Sprintf("%d", time), int32(midPosition.X), int32(midPosition.Y), 100, utils.Purple)
	} else {
		shouldMove = true

	}

}

func DrawNewMapTimer() {
	if goingToNextMap {
		InitPhaseCounter()
	}
}

func DrawInitialMenu() {
	rl.ClearBackground(utils.GrayDark)
	containerOffset := 10
	containerSize := 400
	subcontainerSize := containerSize - 20
	container := rl.NewRectangle(midPosition.X, midPosition.Y, float32(containerSize), float32(containerSize))
	subcontainer := rl.NewRectangle(container.X+float32(containerOffset), container.Y+float32(containerOffset), float32(subcontainerSize), float32(subcontainerSize))
	rl.DrawRectangleRec(container, utils.Purple)
	rl.DrawRectangleRec(subcontainer, utils.GrayDark)
	rl.DrawText("Start: Enter", int32(midPosition.X)-100, int32(midPosition.Y)+100, 20, rl.White)
	rl.DrawText("Quit: ESC", int32(midPosition.X)-100, int32(midPosition.Y)+120, 20, rl.White)
}
