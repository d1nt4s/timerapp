package main

import (
	"context"
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

func (t *Timer) run(cancel context.CancelFunc, s tcell.Screen) {

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
			cancel()
			return
		}

		if t.status == Pause {
			time.Sleep(200 * time.Millisecond)
			continue
		}

		t.decrementSec()
		// drawRemainingTime(s, t.minutes, t.seconds, 0, tcell.StyleDefault.Foreground(tcell.ColorWhite))
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
			drawMessage(screen, "Таймер остановлен", 4, tcell.StyleDefault.Foreground(tcell.ColorRed))
		case "reset":
			t.setup(0, 15)
			drawMessage(screen, "🔁 Таймер сброшен", 4, tcell.StyleDefault.Foreground(tcell.ColorRed))
		case "pause":
			t.status = Pause
			drawMessage(screen, "⏸ Таймер на паузе", 4, tcell.StyleDefault.Foreground(tcell.ColorRed))
		case "resume":
			t.status = Continue
			drawMessage(screen, "▶️ Таймер продолжается", 4, tcell.StyleDefault.Foreground(tcell.ColorRed))
		default:
			drawFormattedMessage(screen, 4, tcell.StyleDefault.Foreground(tcell.ColorYellow), "🤷 Неизвестная команда: %s", cmd)

		}
	default:
	}
}

