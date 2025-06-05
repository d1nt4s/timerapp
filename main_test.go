// main_test.go
package main

import (
	"context"
	"testing"
	"time"
)

func TestTimerRace(t *testing.T) {
	rl, err := createReadline()
	if err != nil {
		t.Fatal(err)
	}
	defer rl.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var timer Timer
	timer.control = make(chan string)
	timer.setup(1, 0)
	timer.status = Continue

	go timer.run(ctx, cancel, rl)

	go func() {
		timer.control <- "pause"
		time.Sleep(100 * time.Millisecond)
		timer.control <- "resume"
		time.Sleep(100 * time.Millisecond)
		timer.control <- "stop"
	}()

	// Ждём завершения таймера
	time.Sleep(2 * time.Second)
}
