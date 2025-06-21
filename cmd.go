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
		defer Debug("🟢 event proxy завершается")

		for {
			select {
			case <-control: // single exit point - if channel control closed -> close scanCommand goroutine
				Debug("⛔ control канал закрыт — proxy завершён")
				close(eventChan)
				return

			default:
				ev := screen.PollEvent()
				eventChan <- ev
			}
		}
	}()

	for ev := range eventChan {
		switch ev := ev.(type) {
		case *tcell.EventKey:
			handleKeyEvent(ev, screen, &buffer, control)
		case *tcell.EventResize:
			screen.Sync()
		// case *tcell.EventMouse:
		// 	x, y := ev.Position()
		// 	if ev.Buttons()&tcell.Button1 != 0 {
		// 		if cmd, ok := handleMouseForButtons(x, y, getVisibleButtons(screen)); ok {
		// 			// handleCommand(cmd)
		// 			control <- cmd
		// 		}
		// 	}
		}
	}

}

func handleKeyEvent(ev *tcell.EventKey, screen tcell.Screen, buffer *[]rune, control chan string) {
	switch ev.Key() {
	case tcell.KeyEnter:

		cmd := strings.ToLower(strings.TrimSpace(string(*buffer)))
		control <- cmd

		clearInputLine(screen)
		*buffer = nil

	default:
		r := ev.Rune()

		if r == 127 || r == '\b' || r == '\x08' {
			if len(*buffer) > 0 {
				*buffer = (*buffer)[:len(*buffer)-1]
				writeToInputLine(screen, *buffer)
			}
		}

		if r >= 32 && r != 127 {
			*buffer = append(*buffer, r)
			writeToInputLine(screen, *buffer)
		}
	}
}
