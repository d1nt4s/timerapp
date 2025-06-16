package main

import (
	"strings"
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

var commandMap = map[string]Command{
	"stop":   CmdStop,
	"reset":  CmdReset,
	"pause":  CmdPause,
	"resume": CmdResume,
}

func ParseCommand(input string) (Command, bool) {
	cmd, ok := commandMap[strings.ToLower(strings.TrimSpace(input))]
	return cmd, ok
}

type Timer struct {
	Minutes int
	Seconds int
	control chan Command
	status  TimerStatus
}

func NewTimer(min, sec int) *Timer {
	return &Timer{
		Minutes: min,
		Seconds: sec,
		control: make(chan Command),
		status:  Continue,
	}
}

func (t *Timer) ControlChan() chan Command {
	return t.control
}

func (t *Timer) Set(min, sec int) {
	t.Minutes = min
	t.Seconds = sec
}

func (t *Timer) Run(s tcell.Screen) {
	for {
		t.handleCommands(s)

		if t.status == End {
			t.drainControlChan()
			return
		}

		if t.status == Pause {
			time.Sleep(200 * time.Millisecond)
			continue
		}

		t.tick()
		drawBigTimer(s, t.Minutes, t.Seconds, 0, tcell.StyleDefault.Foreground(tcell.ColorWhite))
		time.Sleep(time.Second)
	}
}

func (t *Timer) tick() {
	if t.Seconds == 0 {
		if t.Minutes == 0 {
			t.status = End
			return
		}
		t.Minutes--
		t.Seconds = 59
	} else {
		t.Seconds--
	}
	t.status = Continue
}

func (t *Timer) handleCommands(screen tcell.Screen) {
	select {
	case cmd := <-t.control:
		switch cmd {
		case CmdStop:
			t.status = End
			userNotice(screen, "⏹ Таймер остановлен")
		case CmdReset:
			t.Set(0, 15)
			userNotice(screen, "🔁 Таймер сброшен")
		case CmdPause:
			t.status = Pause
			userNotice(screen, "⏸ Таймер на паузе")
		case CmdResume:
			t.status = Continue
			userNotice(screen, "▶️ Таймер продолжается")
		default:
			userError(screen, "🤷 Неизвестная команда: "+string(cmd))
		}
	default:
	}
}

func (t *Timer) drainControlChan() {
	for {
		select {
		case <-t.control:
		default:
			return
		}
	}
}
