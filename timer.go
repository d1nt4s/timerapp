package main

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

type Timer struct {
	Minutes      int
	Seconds      int
	control      chan Command
	status       TimerStatus
	mode         TimerMode
	pauseCounter int
}

func NewTimer(min, sec int, modes ...TimerMode) *Timer {
	mode := Pomodoro
	if len(modes) > 0 {
		mode = modes[0]
	}
	return &Timer{
		Minutes:      min,
		Seconds:      sec,
		control:      make(chan Command),
		status:       Continued,
		mode:         mode,
		pauseCounter: 0,
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
		case Finished:
			t.drainControlChan()
			t.changeMode(s)
			return TimerFinished
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
			t.status = Finished
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

func (t *Timer) changeMode(s tcell.Screen) {
	settings, err := LoadSettings()
	if err != nil {
		userError(s, "💥 Ошибка при загрузке настроек. Выставлены базовые настройки.", false)
	}
	switch t.mode {
	case Pomodoro:
		if t.pauseCounter == settings.LongBreakInterval && settings.LongBreakInterval != 0 {
			t.Set(settings.LongPauseMinutes, settings.LongPauseSeconds)
			t.mode = LongPause
		} else {
			t.Set(settings.PauseMinutes, settings.PauseSeconds)
			t.mode = Pause
		}
	case Pause:
		t.pauseCounter++
		t.Set(settings.PomodoroMinutes, settings.PomodoroSeconds)
		t.mode = Pomodoro
	case LongPause:
		t.pauseCounter = 0
		t.Set(settings.PomodoroMinutes, settings.PomodoroSeconds)
		t.mode = Pomodoro
	}

	t.status = Continued
}
