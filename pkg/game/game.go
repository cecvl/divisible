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
	PausedFrom    State
	CurrentNumber int

	Score       int
	Rounds      int
	TotalRounds int

	StartTime time.Time
	Elapsed   time.Duration
	BestTime  time.Duration
	Duration  time.Duration
}

func New() *Game {
	rand.Seed(time.Now().UnixNano())

	g := &Game{
		TotalRounds: 10,
		Duration:    3 * time.Minute,
	}
	g.Reset()
	return g
}

func (g *Game) Reset() {
	g.Score = 0
	g.Rounds = 0
	g.State = StateQuestion
	g.PausedFrom = StateQuestion
	g.StartTime = time.Now()
	g.Elapsed = 0
	if g.Duration == 0 {
		g.Duration = 3 * time.Minute
	}
	g.NextNumber()
}

func (g *Game) NextNumber() {
	g.CurrentNumber = rand.Intn(90000) + 10000
}

func (g *Game) Update() {
	if g.State == StatePaused {
		if rl.IsKeyPressed(rl.KeyP) {
			g.resume()
		}
		return
	}

	// Only update elapsed time if game is not finished
	if g.State != StateFinished {
		g.Elapsed = time.Since(g.StartTime)
	}

	// If a duration is set, finish when elapsed reaches it
	if g.Duration > 0 && g.Elapsed >= g.Duration {
		g.State = StateFinished

		if g.BestTime == 0 || g.Elapsed < g.BestTime {
			g.BestTime = g.Elapsed
		}
		return
	}

	if rl.IsKeyPressed(rl.KeyP) {
		g.pause()
		return
	}

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

func (g *Game) pause() {
	if g.State == StateQuestion || g.State == StateBonus {
		g.PausedFrom = g.State
		g.State = StatePaused
	}
}

func (g *Game) resume() {
	if g.State != StatePaused {
		return
	}

	g.StartTime = time.Now().Add(-g.Elapsed)
	g.State = g.PausedFrom
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
	padding := int32(20)
	screenWidth := int32(rl.GetScreenWidth())

	// Number (hide when finished)
	if g.State != StateFinished {
		number := fmt.Sprintf("%d", g.CurrentNumber)
		ui.DrawCentered(number, centerY-50, 60, rl.Black)
	}

	// Timer (show remaining when a duration is set, hide when finished)
	if g.State != StateFinished {
		if g.Duration > 0 {
			remaining := g.Duration - g.Elapsed
			if remaining < 0 {
				remaining = 0
			}
			ui.DrawRightAligned("Time: "+formatTime(remaining), screenWidth-padding, padding, 20, rl.DarkGray)
		} else {
			ui.DrawRightAligned("Time: "+formatTime(g.Elapsed), screenWidth-padding, padding, 20, rl.DarkGray)
		}
	}

	// Score
	ui.DrawAt(fmt.Sprintf("Score: %d", g.Score), padding, padding, 20, rl.Gray)

	switch g.State {
	case StateQuestion:
		ui.DrawCentered("Y / N", centerY+40, 20, rl.DarkGray)
		ui.DrawCentered("Press P to Pause", centerY+74, 18, rl.Gray)

	case StateBonus:
		ui.DrawCentered("+1 or +2", centerY+40, 20, rl.DarkGray)
		ui.DrawCentered("Press P to Pause", centerY+74, 18, rl.Gray)

	case StatePaused:
		rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), rl.Fade(rl.Black, 0.4))
		ui.DrawCentered("PAUSED", centerY-20, 40, rl.White)
		ui.DrawCentered("Press P to Resume", centerY+25, 20, rl.LightGray)

	case StateFinished:
		ui.DrawCentered("DONE!", centerY-100, 40, rl.Black)
		ui.DrawCentered("Final Time: "+formatTime(g.Elapsed), centerY-40, 20, rl.DarkGray)
		ui.DrawCentered("Press R to Restart", centerY+60, 20, rl.DarkGray)
	}
}
