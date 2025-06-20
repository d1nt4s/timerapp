package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type SettingType int

const (
	PomodoroSetting SettingType = iota
	PauseSetting
	LongPauseSetting
	IntervalSetting
)

type Settings struct {
	PomodoroMinutes   int `json:"pomodoro_minutes"`
	PomodoroSeconds   int `json:"pomodoro_seconds"`
	PauseMinutes      int `json:"pause_minutes"`
	PauseSeconds      int `json:"pause_seconds"`
	LongPauseMinutes  int `json:"long_pause_minutes"`
	LongPauseSeconds  int `json:"long_pause_seconds"`
	LongBreakInterval int `json:"long_break_interval"` 
}

func DefaultSettingsCopy() Settings {
	return Settings{
		PomodoroMinutes:   25,
		PomodoroSeconds:   0,
		PauseMinutes:      5,
		PauseSeconds:      0,
		LongPauseMinutes:  15,
		LongPauseSeconds:  0,
		LongBreakInterval: 4,
	}
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
		return DefaultSettingsCopy(), nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		// Файл не найден — используем дефолт
		return DefaultSettingsCopy(), nil
	}
	var s Settings
	err = json.Unmarshal(data, &s)
	return s, err
}

