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
			userError(a.screen, "💥 Ошибка при загрузке настроек", false)
			return false
		}
		a.timer = NewTimer(settings.PomodoroMinutes, settings.PomodoroSeconds)
		a.startTimer()

	case strings.HasPrefix(cmd, "set_timer"):
		if min, sec, ok := parseTimeFromSetCommand(cmd, "set_timer"); ok {
			a.applyNewSettings(min, sec, false)
		} else {
			userError(a.screen, "Введите в формате set_timer mm:ss", true)
		}

	case strings.HasPrefix(cmd, "set_pause"):
		if min, sec, ok := parseTimeFromSetCommand(cmd, "set_pause"); ok {
			a.applyNewSettings(min, sec, true)
		} else {
			userError(a.screen, "Введите в формате set_pause mm:ss", true)
		}

	default:
		userError(a.screen, "⭔ Неизвестная команда "+cmd, true)
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
			clearAllExceptInputLine(screen)
			userNotice(screen, "⏹ Таймер остановлен", false)
		case CmdReset:
			settings, err := LoadSettings()
			if err != nil {
				userError(screen, "💥 Ошибка при загрузке настроек", false)
			}
			t.Set(settings.PomodoroMinutes, settings.PomodoroSeconds)
			userNotice(screen, "🔁 Таймер сброшен", true)
		case CmdPause:
			t.status = Paused
			userNotice(screen, "⏸ Таймер на паузе", false)
		case CmdResume:
			t.status = Continued
			userNotice(screen, "▶️ Таймер продолжается", true)
		case CmdExit:
			t.status = ExitApp
			userNotice(screen, "❌ Запрошен выход из программы", false)
		case CmdSkip:
			t.changeMode(screen)
			userNotice(screen, "Пропуск...", true)
		default:
			userError(screen, "⭔ Неизвестная команда: "+string(cmd), true)
		}
	default:
	}
}
