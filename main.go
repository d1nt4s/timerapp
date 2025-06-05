package main

import (
	"fmt"
	"sync"
	"context"
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
	timer.setup(1, 0)
	timer.status = Continue

	go func() {
		defer func() {
			fmt.Println("🟢 scan_command завершается")
			wg.Done()
		}()
        scan_command(ctx, screen, timer.control)
    }()

    go func() {
		defer func() {
			fmt.Println("🟢 timer.run завершается")
			wg.Done()
		}()
        timer.run(cancel)
    }()

	wg.Wait()
	fmt.Println("👋 Программа завершена.")
}
