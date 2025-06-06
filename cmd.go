package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/gdamore/tcell/v2"
)

func setupScreen() tcell.Screen {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка инициализации tcell: %v\n", err)
		os.Exit(1)
	}
	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка инициализации экрана: %v\n", err)
		os.Exit(1)
	}
	return screen
}

func scanCommand(screen tcell.Screen, control chan string) {

	var buffer []rune
	eventChan := make(chan tcell.Event)

	// Горутинa, чтобы проксировать события
	go func() {
		for {
			select {
			// case <-ctx.Done():
			// 	return
			default:
				event := screen.PollEvent()
				eventChan <- event
			}
		}
	}()

	for {
		select {
		// case <-ctx.Done():
		// 	drawMessage(screen, "scan_command: ctx.Done", 6, tcell.StyleDefault.Foreground(tcell.ColorRed))
		// 	return 
		case ev := <-eventChan:
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEnter:
					line := string(buffer)
					// fmt.Println("scan_command: Введена строка:", line)

					cmd := strings.ToLower(strings.TrimSpace(line))
					Debug("Перед отправкой команды в канал: %s" + cmd)
					control <- cmd
					Debug("После отправки команды в канал")
					if cmd == "exit" {
						return
					}

					w, h := screen.Size()
					for x := 0; x < w; x++ {
						screen.SetContent(x, h-1, ' ', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
					}
					buffer = nil 
				case tcell.KeyBackspace:
					if len(buffer) > 0 {
						buffer = buffer[:len(buffer)-1]
					}
				case tcell.KeyCtrlC:
					return
				default:
					r := ev.Rune()
					if r != 0 {
						buffer = append(buffer, r)
					}

					_, height := screen.Size()
					screen.SetContent(len(buffer), height - 1, r, nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
					screen.Show()
        			// x += runewidth.RuneWidth(r)
				}
			case *tcell.EventResize:
				screen.Sync()
			}
		}
	}
}
