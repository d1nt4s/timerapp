package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
)

type App struct {
	screen                 tcell.Screen
	timer                  *Timer
	uiCommandCh            chan string
	acceptingTimerCommands bool
}

func NewApp() *App {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("ошибка создания экрана: %v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("ошибка инициализации экрана: %v", err)
	}

	return &App{
		screen:      screen,
		uiCommandCh: make(chan string),
	}
}

func (a *App) Run() {
	defer func() {
		close(a.uiCommandCh) // Закрываем канал команд ТОЛЬКО здесь
		Debug("🟢 Основной канал сырых команд закрылся")
	}()

	userNotice(a.screen, "⌨️  Введите команду (new / exit):")

	a.timer = NewTimer(1, 0)
	a.acceptingTimerCommands = false

	go func() {
		defer func() {
			Debug("🟢 scanCommand завершается")
		}()
		scanCommand(a.screen, a.uiCommandCh)
	}()

Loop:

	for cmd := range a.uiCommandCh {
		if a.acceptingTimerCommands {
			if parsed, cleaned, ok := ParseCommand(cmd); ok {
				a.timer.control <- parsed
			} else {
				userError(a.screen, "⭔ Команда \""+cleaned+"\" не распознана")
			}
		} else {
			if a.handleCommand(cmd) {
				break Loop
			}
		}
	}
}

func (a *App) handleCommand(cmd string) bool {
	switch cmd {
	case "exit":
		return true
	case "new":
		a.timer = NewTimer(1, 0)
		a.startTimer()
	default:
		userError(a.screen, "⭔ Неизвестная команда "+cmd)
	}

	return false
}

func (a *App) startTimer() {

	a.acceptingTimerCommands = true

	go func() {
		defer func() {
			a.acceptingTimerCommands = false
			Debug("🟢 timer.run завершается")
			userNotice(a.screen, "⏱ Таймер завершён")
			userHint(a.screen, "🐲  Введите 'new' или 'exit'")

		}()
		a.timer.Set(1, 0)
		a.timer.Run(a.screen)

		exitStatus := a.timer.Run(a.screen)

		if exitStatus == TimerExitApp {
			a.uiCommandCh <- "exit"
		}

	}()
}
