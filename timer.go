package main

import (
	"context"
	"fmt"
	"time"
)

type Timer struct {
	seconds int
	minutes int
	control chan string
	status  Status
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
			return
		}
		t.minutes--
		t.seconds = 60
	}
	t.seconds--
	t.status = Continue
}

func (t *Timer) run(cancel context.CancelFunc) {
	for {
		t.manage()

		if t.status == End {
		Drain:
			for {
				select {
				case <-t.control:
				default:
					break Drain
				}
			}
			cancel()
			return
		}

		if t.status == Pause {
			time.Sleep(200 * time.Millisecond)
			continue
		}

		t.decrementSec()
		fmt.Printf("\033[1A\033[2K⏳ Осталось: %d мин %02d сек\n", t.minutes, t.seconds)

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
			fmt.Printf("\n Таймер остановлен\n")
		case "reset":
			t.setup(0, 15)
			fmt.Printf("\n🔁 Таймер сброшен\n")
		case "pause":
			t.status = Pause
			fmt.Printf("\n⏸ Таймер на паузе\n")
		case "resume":
			t.status = Continue
			fmt.Printf("\n▶️ Таймер продолжается\n")
		case "exit":
			fmt.Println("t")
			t.status = End
		default:
			fmt.Printf("\n🤷 Неизвестная команда: %s\n", cmd)
		}
	default:
	}
}
