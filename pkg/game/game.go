package game

import (
	"fmt"
	"math/rand"
	"time"

	"games/example.com/pkg/ui"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Game struct {
	State         State
	CurrentNumber int

	Score       int
	Rounds      int
	TotalRounds int

	StartTime time.Time
	Elapsed   time.Duration
	BestTime  time.Duration
}

func New() *Game {
	rand.Seed(time.Now().UnixNano())

	g := &Game{
		TotalRounds: 10,
	}
	g.Reset()
	return g
}

func (g *Game) Reset() {
	g.Score = 0
	g.Rounds = 0
	g.State = StateQuestion
	g.StartTime = time.Now()
	g.NextNumber()
}

func (g *Game) NextNumber() {
	g.CurrentNumber = rand.Intn(90000) + 10000
}

func (g *Game) Update() {
	g.Elapsed = time.Since(g.StartTime)

	switch g.State {
	case StateQuestion:
		g.handleQuestionInput()
	case StateBonus:
		g.handleBonusInput()
	case StateFinished:
		if rl.IsKeyPressed(rl.KeyR) {
			g.Reset()
		}
	}
}

func (g *Game) handleQuestionInput() {
	if rl.IsKeyPressed(rl.KeyY) {
		g.checkAnswer(true)
	}
	if rl.IsKeyPressed(rl.KeyN) {
		g.checkAnswer(false)
	}
}

func (g *Game) handleBonusInput() {
	if rl.IsKeyPressed(rl.KeyOne) || rl.IsKeyPressed(rl.KeyKp1) {
		g.checkBonus(1)
	}
	if rl.IsKeyPressed(rl.KeyTwo) || rl.IsKeyPressed(rl.KeyKp2) {
		g.checkBonus(2)
	}
}

func (g *Game) checkAnswer(yes bool) {
	correct := IsDivisibleBy3(g.CurrentNumber)

	if yes == correct {
		if correct {
			g.Score++
			g.nextRound()
		} else {
			g.State = StateBonus
		}
	} else {
		g.nextRound()
	}
}

func (g *Game) checkBonus(choice int) {
	correct := NeededToMakeDivisible(g.CurrentNumber)

	if choice == correct {
		g.Score++
	}

	g.nextRound()
}

func (g *Game) nextRound() {
	g.Rounds++

	if g.Rounds >= g.TotalRounds {
		g.State = StateFinished

		if g.BestTime == 0 || g.Elapsed < g.BestTime {
			g.BestTime = g.Elapsed
		}
		return
	}

	g.NextNumber()
	g.State = StateQuestion
}

func formatTime(d time.Duration) string {
	totalSeconds := int(d.Seconds())

	minutes := totalSeconds / 60
	seconds := totalSeconds % 60

	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func (g *Game) Draw() {
	centerY := int32(rl.GetScreenHeight() / 2)

	// Number
	number := fmt.Sprintf("%d", g.CurrentNumber)
	ui.DrawCentered(number, centerY-50, 60, rl.Black)

	// Timer
	ui.DrawCentered("Time: "+formatTime(g.Elapsed), 20, 20, rl.DarkGray)

	// Score
	ui.DrawCentered(fmt.Sprintf("Score: %d", g.Score), 50, 20, rl.Gray)

	switch g.State {
	case StateQuestion:
		ui.DrawCentered("Y / N", centerY+40, 20, rl.DarkGray)

	case StateBonus:
		ui.DrawCentered("+1 or +2", centerY+40, 20, rl.DarkGray)

	case StateFinished:
		ui.DrawCentered("DONE!", centerY-100, 40, rl.Black)
		ui.DrawCentered("Final Time: "+formatTime(g.Elapsed), centerY-40, 20, rl.DarkGray)

		if g.BestTime > 0 {
			ui.DrawCentered("Best: "+formatTime(g.BestTime), centerY, 20, rl.Gray)
		}

		ui.DrawCentered("Press R to Restart", centerY+60, 20, rl.DarkGray)
	}
}
