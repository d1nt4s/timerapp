package main

import (
	"context"
	"time"

	"github.com/gdamore/tcell/v2"
)

type TimerStatus int

const (
	Continue TimerStatus = iota
	Pause
	End
)

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

func (t *Timer) Run(cancel context.CancelFunc, screen tcell.Screen) {
	for {
		t.manage(screen)

		if t.status == End {
			t.drainControl()
			cancel()
			return
		}

		if t.status == Pause {
			time.Sleep(200 * time.Millisecond)
			continue
		}

		t.updateTime()
		drawBigTimer(screen, t.Minutes, t.Seconds, 0, tcell.StyleDefault.Foreground(tcell.ColorWhite))
		time.Sleep(time.Second)
	}
}

func (t *Timer) manage(screen tcell.Screen) {
	select {
	case cmd := <-t.control:
		switch cmd {
		case "stop":
			t.status = End
			drawMessage(screen, "⛔ Таймер остановлен", 4, tcell.StyleDefault.Foreground(tcell.ColorRed))
		case "reset":
			t.Minutes = 0
			t.Seconds = 15
			t.status = Continue
			drawMessage(screen, "🔁 Таймер сброшен", 4, tcell.StyleDefault.Foreground(tcell.ColorYellow))
		case "pause":
			t.status = Pause
			drawMessage(screen, "⏸ Таймер на паузе", 4, tcell.StyleDefault.Foreground(tcell.ColorYellow))
		case "resume":
			t.status = Continue
			drawMessage(screen, "▶️ Продолжение таймера", 4, tcell.StyleDefault.Foreground(tcell.ColorGreen))
		default:
			drawFormattedMessage(screen, 4, tcell.StyleDefault.Foreground(tcell.ColorOrange), "🤷 Неизвестная команда: %s", cmd)
		}
	default:
	}
}

func (t *Timer) updateTime() {
	if t.Seconds == 0 {
		if t.Minutes > 0 {
			t.Minutes--
			t.Seconds = 59
		}
	} else {
		t.Seconds--
	}
}

func (t *Timer) drainControl() {
Drain:
	for {
		select {
		case <-t.control:
		default:
			break Drain
		}
	}
}

func (t *Timer) Pause()  { t.control <- "pause" }
func (t *Timer) Resume() { t.control <- "resume" }
func (t *Timer) Stop()   { t.control <- "stop" }
func (t *Timer) Reset()  { t.control <- "reset" }
