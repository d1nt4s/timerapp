package main

import (
	"strconv"
	"strings"
	"fmt"
)

func (a *App) updateSettingFromCommand(cmd, prefix string, settingType SettingType, startAfter bool) {
	settings, err := LoadSettings()
	if err != nil {
		userError(a.screen, "💥 Ошибка при загрузке настроек", false)
		return
	}

	switch settingType {
	case PomodoroSetting, PauseSetting, LongPauseSetting:
		min, sec, ok := parseTimeFromSetCommand(cmd, prefix)
		if !ok {
			userError(a.screen, fmt.Sprintf("Введите в формате %s mm:ss", prefix), true)
			return
		}

		switch settingType {
		case PomodoroSetting:
			settings.PomodoroMinutes = min
			settings.PomodoroSeconds = sec
		case LongPauseSetting:
			settings.LongPauseMinutes = min
			settings.LongPauseSeconds = sec
		case PauseSetting:
			settings.PauseMinutes = min
			settings.PauseSeconds = sec
		}

		if startAfter && settingType == PomodoroSetting {
			a.timer = NewTimer(min, sec, Pomodoro)
			a.startTimer()
		}

	case IntervalSetting:
		val, ok := parseIntFromCommand(cmd, prefix)
		if !ok {
			userError(a.screen, "Введите в формате set_interval {число}, при 0 длинные паузы отключаются", true)
			return
		}
		settings.LongBreakInterval = val
	}

	if err := SaveSettings(settings); err != nil {
		userError(a.screen, "💥 Ошибка при сохранении настроек", false)
		return
	}

	userNotice(a.screen, "💾 Настройки по умолчанию сохранены!", true)
}



func parseTimeFromSetCommand(input, prefix string) (int, int, bool) {
	trimmed := strings.TrimPrefix(input, prefix+" ")
	parts := strings.Split(trimmed, ":")
	if len(parts) != 2 {
		return 0, 0, false
	}
	min, err1 := strconv.Atoi(parts[0])
	sec, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil || min < 0 || sec < 0 || sec > 59 {
		return 0, 0, false
	}
	return min, sec, true
}

func parseIntFromCommand(input, prefix string) (int, bool) {
	trimmed := strings.TrimPrefix(input, prefix+" ")
	val, err := strconv.Atoi(trimmed)
	if err != nil || val < 0 {
		return 0, false
	}
	return val, true
}
