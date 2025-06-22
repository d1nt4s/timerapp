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

func getButtons(screen tcell.Screen, mode AppMode) []Button {
	_, screenHeight := screen.Size()
	startY := screenHeight - 3

	var buttons []Button

	// Кнопки для SetupMode
	if mode == SetupMode {
		buttons = []Button{
			{Label: "[start]", Command: "start", X: 12, Y: startY, Visible: true},
			{Label: "[help]", Command: "help", X: 22, Y: startY, Visible: true},
			{Label: "[exit]", Command: "exit", X: 32, Y: startY, Visible: true},
		}
	}

	// Кнопки для ActiveMode
	if mode == ActiveMode {
		width, _ := screen.Size() // получаем ширину экрана
		buttonLabels := []struct {
			Label   string
			Command string
		}{
			{"[pause]", "pause"},
			{"[resume]", "resume"},
			{"[reset]", "reset"},
			{"[stop]", "stop"},
			{"[skip]", "skip"},
			{"[snooze5m]", "snooze5m"},
			{"[snooze10m]", "snooze10m"},
			{"[help]", "help"},
			{"[exit]", "exit"},
		}

		totalButtons := len(buttonLabels)
		spacePerButton := width / totalButtons

		buttons = []Button{}
		for i, item := range buttonLabels {
			x := i*spacePerButton + (spacePerButton-len(item.Label))/2 // центрируем каждую кнопку в своём слоте
			buttons = append(buttons, Button{
				Label:   item.Label,
				Command: item.Command,
				X:       x,
				Y:       startY,
				Visible: true,
			})
		}
	}


	return buttons
}


func inRect(mx, my, x, y, w, h int) bool {
	return mx >= x && mx < x+w && my >= y && my < y+h
}

func clearButtonLine(s tcell.Screen) {
	width, height := s.Size()
	y := height - 3 // строка, на которой рисуются кнопки

	for x := 0; x < width; x++ {
		s.SetContent(x, y, ' ', nil, tcell.StyleDefault)
	}

	s.Show()
}
