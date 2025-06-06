package main

import (
	"context"
	"sync"
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

	app := NewApp();
	acceptingTimerCommands := false

	go func() {
		defer func() {
			Debug("🟢 scanCommand завершается")
			wg.Done()
		}()
        scanCommand(ctx, screen, app.uiCommandCh)
    }()

	Loop:
	for {
		
		select {
		case cmd := <-app.uiCommandCh:
			if acceptingTimerCommands {
				timer.control <- cmd
			} else {
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
							Debug("🟢 timer.run завершается")
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
}
