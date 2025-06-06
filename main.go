package main

import (
	"context"
	"sync"
	"time"
	"github.com/gdamore/tcell/v2"
)

func main() {
	var timer Timer

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

	screen := setupScreen()
	defer screen.Fini()

	var wg sync.WaitGroup
	wg.Add(2)

	timer.control = make(chan string)
	timer.setup(0, 1)
	timer.status = Continue

	// app := NewApp();

	commandCh := make(chan string)
	acceptingTimerCommands := false

	go func() {
		defer func() {
			drawMessage(screen, "🟢 scan_command завершается", 2, tcell.StyleDefault.Foreground(tcell.ColorRed))
			wg.Done()
		}()
        scanCommand(ctx, screen, commandCh)
    }()

	Loop:
	for {
		
		drawMessage(screen, "перед select", 9, tcell.StyleDefault.Foreground(tcell.ColorRed))
		select {
		case cmd := <-commandCh:
			if acceptingTimerCommands {
				timer.control <- cmd
			} else {
				drawFormattedMessage(screen, 7, tcell.StyleDefault.Foreground(tcell.ColorYellow), "перед switch cmd: %s", cmd)
				switch cmd {
				case "exit":
					drawMessage(screen, "👋 Выход", 3, tcell.StyleDefault.Foreground(tcell.ColorRed))
					break Loop
				case "new":
					timer.setup(0, 1)
					timer.status = Continue
					acceptingTimerCommands = true

					go func() {
						defer func() {
							acceptingTimerCommands = false
							drawMessage(screen, "🟢 timer.run завершается", 3, tcell.StyleDefault.Foreground(tcell.ColorRed))
							wg.Done()
							drawMessage(screen, "⏱ Таймер завершён", 4, tcell.StyleDefault.Foreground(tcell.ColorGreen))
							drawMessage(screen, "✏️  Введите 'new' или 'exit'", 5, tcell.StyleDefault.Foreground(tcell.ColorGreen))

						}()
						timer.run(cancel, screen)

					}()
				default:
					drawMessage(screen, "🤷 Неизвестная команда", 4, tcell.StyleDefault.Foreground(tcell.ColorGreen))
				}
			}
		}
	}


	wg.Wait()
	drawMessage(screen, "👋 Программа завершена.", 5, tcell.StyleDefault.Foreground(tcell.ColorRed))

	time.Sleep(time.Second * 10)
}
