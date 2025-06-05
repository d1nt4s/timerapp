package main

import (
	"context"
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

func scan_command(ctx context.Context, screen tcell.Screen, control chan string) {

	var buffer []rune
	eventChan := make(chan tcell.Event)

	// Горутинa, чтобы проксировать события
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
		select {
		case <-ctx.Done():
			fmt.Println("scan_command: ctx.Done")
			return // безопасный выход
		case ev := <-eventChan:
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEnter:
					line := string(buffer)
					fmt.Println("scan_command: Введена строка:", line)

					cmd := strings.ToLower(strings.TrimSpace(line))
					fmt.Println("scan_command: scan⏳ Перед отправкой команды в канал")
					control <- cmd
					fmt.Println("scan_command: ⏳ Перед отправкой команды в канал")
					if cmd == "exit" {
						return
					}

					buffer = nil // очистить буфер
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
				}
			case *tcell.EventResize:
				screen.Sync()
			}
		}
	}
}
