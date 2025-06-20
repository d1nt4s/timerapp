package main

import (
	"strings"
	"github.com/gdamore/tcell/v2"
)

// UI handling

func (a *App) handleCommand(cmd string) bool {
	switch {
	case cmd == "exit":
		return true

	case cmd == "start":
		settings, err := LoadSettings()
		if err != nil {
			userError(a.screen, "💥 Ошибка при загрузке настроек")
			return false
		}
		a.timer = NewTimer(settings.PomodoroMinutes, settings.PomodoroSeconds)
		a.startTimer()

	case strings.HasPrefix(cmd, "set_timer"):
		if min, sec, ok := parseTimeFromSetCommand(cmd, "set_timer"); ok {
			a.applyNewSettings(min, sec, false)
		} else {
			userError(a.screen, "Введите в формате set_timer mm:ss")
		}

	case strings.HasPrefix(cmd, "set_pause"):
		if min, sec, ok := parseTimeFromSetCommand(cmd, "set_pause"); ok {
			a.applyNewSettings(min, sec, true)
		} else {
			userError(a.screen, "Введите в формате set_pause mm:ss")
		}

	default:
		userError(a.screen, "⭔ Неизвестная команда "+cmd)
	}

	return false
}

// Timer Handling

func (t *Timer) handleCommands(screen tcell.Screen) {
	select {
	case cmd := <-t.control:
		switch cmd {
		case CmdStop:
			t.status = Stopped
			userNotice(screen, "⏹ Таймер остановлен")
		case CmdReset:
			t.Set(0, 15)
			userNotice(screen, "🔁 Таймер сброшен")
		case CmdPause:
			t.status = Paused
			userNotice(screen, "⏸ Таймер на паузе")
		case CmdResume:
			t.status = Continued
			userNotice(screen, "▶️ Таймер продолжается")
		case CmdExit:
			t.status = ExitApp
			userNotice(screen, "❌ Запрошен выход из программы")
		default:
			userError(screen, "⭔ Неизвестная команда: "+string(cmd))
		}
	default:
	}
}