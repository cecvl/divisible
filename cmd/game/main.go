package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"games/example.com/pkg/game"
)

func main() {
	rl.InitWindow(800, 450, "Divisible by 3")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	g := game.New()

	for !rl.WindowShouldClose() {
		g.Update()

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		g.Draw()

		rl.EndDrawing()
	}
}
