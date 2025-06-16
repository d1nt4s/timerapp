package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
)

type TimerStatus int

const (
	Continue TimerStatus = iota
	Pause
	End
)

type Command string

const (
	CmdStop   Command = "stop"
	CmdReset  Command = "reset"
	CmdPause  Command = "pause"
	CmdResume Command = "resume"
)

var TimerCommandChan = make(chan Command)

type Timer struct {
	Minutes int
	Seconds int
	control chan string
	status  TimerStatus
}

func NewTimer(min, sec int) *Timer {
	return &Timer{
		Minutes: min,
		Seconds: sec,
		control: make(chan string),
		status:  Continue,
	}
}

func (t *Timer) setTimer(min, sec int) {
	t.Minutes = min
	t.Seconds = sec
}

func (t *Timer) decrementSec() {
	if t.Seconds == 0 {
		if t.Minutes == 0 {
			t.status = End
			return
		}
		t.Minutes--
		t.Seconds = 60
	}
	t.Seconds--
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
		drawBigTimer(s, t.Minutes, t.Seconds, 0, tcell.StyleDefault.Foreground(tcell.ColorWhite))

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
			t.setTimer(0, 15)
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
