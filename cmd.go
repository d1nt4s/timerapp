package main

import (
	"strings"

	"github.com/gdamore/tcell/v2"
)

func scanCommand(screen tcell.Screen, control chan string) {
	var buffer []rune
	eventChan := make(chan tcell.Event)

	// Асинхронная проксирующая горутина
	go func() {
		for {
			event := screen.PollEvent()
			eventChan <- event
		}
	}()

	for ev := range eventChan {
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if handleKeyEvent(ev, screen, &buffer, control) {
				return // пользователь ввёл "exit" или нажал Ctrl+C
			}
		case *tcell.EventResize:
			screen.Sync()
		}
	}
}

func handleKeyEvent(ev *tcell.EventKey, screen tcell.Screen, buffer *[]rune, control chan string) (exit bool) {
	switch ev.Key() {
	case tcell.KeyEnter:
		cmd := strings.ToLower(strings.TrimSpace(string(*buffer)))
		Debug("Перед отправкой команды в канал: %s" + cmd)
		control <- cmd
		Debug("После отправки команды в канал")

		if cmd == "exit" {
			return true
		}

		clearInputLine(screen)
		*buffer = nil

	case tcell.KeyBackspace:
		if len(*buffer) > 0 {
			*buffer = (*buffer)[:len(*buffer)-1]
		}

	case tcell.KeyCtrlC:
		return true

	default:
		r := ev.Rune()
		if r != 0 {
			*buffer = append(*buffer, r)
			writeToInputLine(screen, *buffer)
		}
	}

	return false
}

func writeToInputLine(screen tcell.Screen, buffer []rune) {
	width, height := screen.Size()
	for x := 0; x < width; x++ {
		screen.SetContent(x, height-1, ' ', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
	}
	for i, r := range buffer {
		screen.SetContent(i, height-1, r, nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
	}
	screen.Show()
}

func clearInputLine(screen tcell.Screen) {
	width, height := screen.Size()
	for x := 0; x < width; x++ {
		screen.SetContent(x, height-1, ' ', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
	}
	screen.Show()
}
