package main

import (
	"fmt"
	"time"

	"github.com/chzyer/readline"
)

type Timer struct {
	seconds int
	minutes int
	control chan string
	done chan bool
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

func (t *Timer) run(rl *readline.Instance) {
	for {
		t.manage(rl)

		if t.status == End {
			return
		}

		if t.status == Pause {
			time.Sleep(200 * time.Millisecond)
			continue
		}

		t.decrementSec()
		// fmt.Printf("\033[2K\r Minutes: %d, Seconds: %d ", t.minutes, t.seconds)
		rl.Write([]byte(fmt.Sprintf("\033[1A\033[2K⏳ Осталось: %d мин %02d сек\n", t.minutes, t.seconds)))
		rl.Refresh()

		time.Sleep(time.Second)
	}
}

func (t *Timer) manage(rl *readline.Instance) {
	select {
	case cmd := <-t.control:
		switch cmd {
		case "stop":
			t.setup(0, 0)
			t.status = End
			rl.Write([]byte("\n Таймер остановлен\n"))
		case "reset":
			t.setup(0, 0)
			rl.Write([]byte("\n🔁 Таймер сброшен\n"))
		case "pause":
			t.status = Pause
			rl.Write([]byte("\n⏸ Таймер на паузе\n"))
		case "resume":
			t.status = Continue
			rl.Write([]byte("\n▶️ Таймер продолжается\n"))
		case "exit":
			t.status = End
			t.done <- true
		default:
			rl.Write([]byte(fmt.Sprintf("\n🤷 Неизвестная команда: %s\n", cmd)))
		}
	default:
	}
}

