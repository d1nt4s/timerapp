package main

import (
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

func ParseCommand(input string) (Command, string, bool) {
	cleaned := strings.ToLower(strings.TrimSpace(input))
	cmd, ok := commandMap[cleaned]
	return cmd, cleaned, ok
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
		status:  Continued,
	}
}

func (t *Timer) Set(min, sec int) {
	t.Minutes = min
	t.Seconds = sec
}

func (t *Timer) Run(s tcell.Screen) TimerResult {
	for {
		t.handleCommands(s)

		switch t.status {
		case Stopped:
			t.drainControlChan()
			return TimerStopped
		case ExitApp:
			t.drainControlChan()
			return TimerExitApp
		case Paused:
			time.Sleep(200 * time.Millisecond)
			continue
		}

		t.tick()
		drawCenteredBigTimer(s, t.Minutes, t.Seconds, tcell.StyleDefault.Foreground(tcell.ColorWhite))
		time.Sleep(time.Second)
	}
}

func (t *Timer) tick() {
	if t.Seconds == 0 {
		if t.Minutes == 0 {
			t.status = Stopped
			return
		}
		t.Minutes--
		t.Seconds = 59
	} else {
		t.Seconds--
	}
	t.status = Continued
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
