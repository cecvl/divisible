package ui

import rl "github.com/gen2brain/raylib-go/raylib"

func DrawCentered(text string, y int32, size int32, color rl.Color) {
	w := int32(rl.MeasureText(text, int32(size)))
	sw := int32(rl.GetScreenWidth())

	x := int32((sw - w) / 2)

	rl.DrawText(text, x, y, size, color)
}
