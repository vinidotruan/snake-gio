package main

import (
	"encoding/json"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth = 900 
	screenHeight = 600
	snakeSize = 20
	frameRate = 200
)

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

