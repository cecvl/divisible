package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"games/example.com/pkg/game"
)

func main() {
	rl.InitWindow(800, 450, "D I V I S I B L E")
	defer rl.CloseWindow()

	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()

	rl.SetTargetFPS(60)

	g := game.New()
	g.InitAudio()
	defer g.Close()

	for !rl.WindowShouldClose() {
		g.Update()

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		g.Draw()

		rl.EndDrawing()
	}
}
