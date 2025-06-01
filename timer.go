package main

import (
	"fmt"
	"time"
)

type Timer struct {
	seconds int
	minutes int
	control chan string
	status Status
}

type Status int

const (
	Continue Status = iota
	Pause
	End
)

func (t *Timer) setup(sec int, min int) {
	t.seconds = sec
	t.minutes = min
}

func (t *Timer) decrementSec() {
	if t.seconds == 0 {
		if t.minutes == 0 {
			t.status = End
		}
		t.minutes--
		t.seconds = 60
	}
	t.seconds--
	t.status = Continue
}

func (t *Timer) run() {
	for t.status == Continue {
		t.decrementSec()
		fmt.Printf("\033[2K\r Minutes: %d, Seconds: %d ", t.minutes, t.seconds)
		t.manage()
		time.Sleep(time.Second)
	}
}

func (t *Timer) manage() {
	select {
	case cmd := <-t.control:
		switch cmd {
		case "stop":
			t.setup(0, 0)
			t.status = End
		case "reset":
			t.setup(0, 0)
		case "pause":
			t.status = Pause
		case "resume":
			t.status = Continue
		default:
			fmt.Println("⚠️ Неизвестная команда:", cmd)
		}
	default:
	}
}

