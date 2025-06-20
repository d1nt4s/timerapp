package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Settings struct {
	PomodoroMinutes int `json:"pomodoro_minutes"`
	PomodoroSeconds int `json:"pomodoro_seconds"`
	PauseMinutes int `json:"pause_minutes"`
	PauseSeconds int `json:"pause_seconds"`
}

func getSettingsPath() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(base, "timerapp")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "settings.json"), nil
}

func SaveSettings(s Settings) error {
	path, err := getSettingsPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func LoadSettings() (Settings, error) {
	path, err := getSettingsPath()
	if err != nil {
		return Settings{PomodoroMinutes: 25, PomodoroSeconds: 0, PauseMinutes: 5, PauseSeconds: 1}, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		// Файл не найден — используем дефолт
		return Settings{PomodoroMinutes: 25, PomodoroSeconds: 0, PauseMinutes: 5, PauseSeconds: 1}, nil
	}
	var s Settings
	err = json.Unmarshal(data, &s)
	return s, err
}

func (a *App) applyNewSettings(min, sec int, isPause bool) {
	newSettings := Settings{}

	oldSettings, err := LoadSettings();
	if err != nil {
		userError(a.screen, "💥 Ошибка при загрузки настроек")
	}

	if isPause {
		newSettings.PauseMinutes = min
		newSettings.PauseSeconds = sec
		newSettings.PomodoroMinutes = oldSettings.PomodoroMinutes
		newSettings.PomodoroSeconds = oldSettings.PomodoroSeconds
	} else {
		newSettings.PomodoroMinutes = min
		newSettings.PomodoroSeconds = sec
		newSettings.PauseMinutes = oldSettings.PauseMinutes
		newSettings.PauseSeconds = oldSettings.PauseSeconds
	}

	if err := SaveSettings(newSettings); err != nil {
		userError(a.screen, "💥 Ошибка при сохранении настроек")
		return
	}

	userNotice(a.screen, "💾 Настройки по умолчанию сохранены!")

	a.timer = NewTimer(min, sec)
	a.startTimer()
}

