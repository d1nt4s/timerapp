package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
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

func (t *Timer) run(s tcell.Screen) {

	for {
		t.manage(s)

		if t.status == End {
		Drain:
			for {
				select {
				case <-t.control:
				default:
					break Drain
				}
			}
			fmt.Println()
			return
		}

		if t.status == Pause {
			time.Sleep(200 * time.Millisecond)
			continue
		}

		t.decrementSec()
		drawBigTimer(s, t.minutes, t.seconds, 0, tcell.StyleDefault.Foreground(tcell.ColorWhite))

		time.Sleep(time.Second)
	}
}

func (t *Timer) manage(screen tcell.Screen) {
	select {
	case cmd := <-t.control:
		switch cmd {
		case "stop":
			t.status = End
			userNotice(screen, "Таймер остановлен")
		case "reset":
			t.setup(0, 15)
			userNotice(screen, "🔁 Таймер сброшен")
		case "pause":
			t.status = Pause
			userNotice(screen, "⏸ Таймер на паузе")
		case "resume":
			t.status = Continue
			userNotice(screen, "▶️ Таймер продолжается")
		default:
			userError(screen, "🤷 Неизвестная команда "+cmd)

		}
	default:
	}
}
