package main

import (
	"github.com/gdamore/tcell/v2"
)

type Button struct {
	Label    string
	Command  string
	X, Y     int
	Visible  bool
}

func drawButtons(s tcell.Screen, buttons []Button, style tcell.Style) {
	for _, b := range buttons {
		if !b.Visible {
			continue
		}
		for i, r := range b.Label {
			s.SetContent(b.X+i, b.Y, r, nil, style)
		}
	}
	s.Show()
}

func handleMouseForButtons(x, y int, buttons []Button) (string, bool) {
	for _, b := range buttons {
		if b.Visible && inRect(x, y, b.X, b.Y, len(b.Label), 1) {
			return b.Command, true
		}
	}
	return "", false
}

func getVisibleButtons(screen tcell.Screen) []Button {
	_, screenHeight := screen.Size()
	startY := screenHeight - 3

	buttons := []Button{
		{Label: "[start]", Command: "start", X: 2, Y: startY},
		{Label: "[pause]", Command: "pause", X: 12, Y: startY},
		{Label: "[resume]", Command: "resume", X: 22, Y: startY},
		{Label: "[reset]", Command: "reset", X: 34, Y: startY},
		{Label: "[stop]", Command: "stop", X: 44, Y: startY},
		{Label: "[skip]", Command: "skip", X: 54, Y: startY},
		{Label: "[help]", Command: "help", X: 64, Y: startY},
		{Label: "[exit]", Command: "exit", X: 74, Y: startY},
	}

	// Показываем только кнопки, соответствующие текущему состоянию таймера
	for i := range buttons {
		switch buttons[i].Command {
		case "start":
			// buttons[i].Visible = app.timer == nil || app.timer.status == Stopped
		case "pause":
			// buttons[i].Visible = app.timer != nil && app.timer.status == Continued
		case "resume":
			// buttons[i].Visible = app.timer != nil && app.timer.status == Paused
		case "reset", "stop", "skip":
			// buttons[i].Visible = app.timer != nil && (app.timer.status == Continued || app.timer.status == Paused)
		case "help", "exit":
			buttons[i].Visible = true
		}
	}

	return buttons
}

func inRect(mx, my, x, y, w, h int) bool {
	return mx >= x && mx < x+w && my >= y && my < y+h
}
