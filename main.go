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
	timer.setup(5, 0)
	timer.status = Continue

	go func() {
		defer func() {
			drawMessage(screen, "🟢 scan_command завершается", 2, tcell.StyleDefault.Foreground(tcell.ColorRed))
			wg.Done()
		}()
        scan_command(ctx, screen, timer.control)
    }()

    go func() {
		defer func() {
			drawMessage(screen, "🟢 timer.run завершается", 3, tcell.StyleDefault.Foreground(tcell.ColorRed))
			wg.Done()
		}()
        timer.run(cancel, screen)
    }()

	wg.Wait()
	drawMessage(screen, "👋 Программа завершена.", 5, tcell.StyleDefault.Foreground(tcell.ColorRed))

	time.Sleep(time.Second * 10)
}
