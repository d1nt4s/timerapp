package main

import (
	"strings"

	"github.com/gdamore/tcell/v2"
)

// UI handling

func (a *App) handleCommand(cmd Command, cmd_string string) bool {
	switch cmd {
	case CmdExit:
		return true

	case CmdStart:
		settings, err := LoadSettings()
		if err != nil {
			userError(a.screen, "💥 Ошибка при загрузке настроек", false)
			return false
		}
		a.timer = NewTimer(settings.PomodoroMinutes, settings.PomodoroSeconds)
		a.startTimer()
	
	case CmdHelp:
		drawLongNotice(a.screen, "Управляй настройкой таймера через следующие команды: set_timer mm:ss - для установки pomodoro, set_pause mm:ss - для установки паузы, set_longpause mm:ss - для установки долгой паузы, set_interval {value} - через сколько пауз будет длинная пауза. Для просмотра команд таймера введите help в режиме таймера.")

	case CmdSetPause, CmdSetTimer, CmdSetInterval, CmdSetLongPause:
		a.handleSetCommand(cmd_string)

	// case CmdSettings:
	// 	showSettingsModal(a)

	default:
		userError(a.screen, "⭔ Неизвестная команда "+string(cmd), true)
	}

	return false
}

func (a *App) handleSetCommand(cmd string) {
	switch {
	case strings.HasPrefix(cmd, "set_timer"):
		a.updateSettingFromCommand(cmd, "set_timer", PomodoroSetting, true)

	case strings.HasPrefix(cmd, "set_pause"):
		a.updateSettingFromCommand(cmd, "set_pause", PauseSetting, false)

	case strings.HasPrefix(cmd, "set_longpause"):
		a.updateSettingFromCommand(cmd, "set_longpause", LongPauseSetting, false)

	case strings.HasPrefix(cmd, "set_interval"):
		a.updateSettingFromCommand(cmd, "set_interval", IntervalSetting, false)
	}
}

// Timer Handling

func (t *Timer) handleCommands(screen tcell.Screen) {
	select {
	case cmd := <-t.control:
		switch cmd {
		case CmdStop:
			t.status = Stopped
			clearBigTimerArea(screen)
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
		case CmdHelp:
			drawLongNotice(screen, "Управляй таймером через следующие команды: stop - для полной остановки таймера, reset - для перезапуска таймера, pause - для приостановки таймера, resume - для возобновления таймера, exit - для выхода из таймера, skip - пропустить таймер.")
		default:
			userError(screen, "⭔ Неизвестная команда: "+string(cmd), true)
		}
	default:
	}
}
