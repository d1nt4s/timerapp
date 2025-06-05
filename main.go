package main

import (
	"fmt"
	"sync"
	"context"
)

func main() {
	var timer Timer
	rl, err := createReadline()
	if err != nil {
		panic(err)
	}

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

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
	// go scan_command(rl, timer.control, timer.done)
        scan_command(ctx, rl, timer.control)
    }()

    go func() {
		defer func() {
			fmt.Println("🟢 timer.run завершается")
			wg.Done()
		}()
        timer.run(ctx, cancel, rl)
    }()

	wg.Wait()
	_ = rl.Close()
	fmt.Println("👋 Программа завершена.")
}
