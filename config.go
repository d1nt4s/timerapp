package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Settings struct {
	DefaultMinutes int `json:"default_minutes"`
	DefaultSeconds int `json:"default_seconds"`
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
		return Settings{DefaultMinutes: 1, DefaultSeconds: 0}, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		// Файл не найден — используем дефолт
		return Settings{DefaultMinutes: 1, DefaultSeconds: 0}, nil
	}
	var s Settings
	err = json.Unmarshal(data, &s)
	return s, err
}
