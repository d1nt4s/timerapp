package main

import (
	"strings"

	"github.com/gdamore/tcell/v2"
)

// scanCommand читает пользовательский ввод и отправляет команды в канал
func scanCommand(screen tcell.Screen, out chan<- string) {
	var input []rune
	_, h := screen.Size()
	inputY := h - 1

	updateInputLine := func() {
		clearLine(screen, inputY, tcell.StyleDefault)
		printAt(screen, 0, inputY, string(input), tcell.StyleDefault)
		screen.ShowCursor(len(input), inputY)
		screen.Show()
	}

	screen.ShowCursor(0, inputY)
	screen.Show()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				event := screen.PollEvent()
				eventChan <- event
			}
		}
	}()

	for {
		event := screen.PollEvent()

		switch ev := event.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEnter:
				cmd := strings.TrimSpace(string(input))
				if cmd != "" {
					out <- cmd
				}
				input = nil

			case tcell.KeyBackspace, tcell.KeyBackspace2:
				if len(input) > 0 {
					input = input[:len(input)-1]
				}

			case tcell.KeyRune:
				input = append(input, ev.Rune())

			case tcell.KeyCtrlC:
				out <- "exit"
			}

			updateInputLine()
		}
	}
}
