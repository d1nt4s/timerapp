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
	for {

		select {
		case cmd := <-a.uiCommandCh:
			if a.acceptingTimerCommands {
				a.timer.control <- cmd
			} else {
				if a.handleCommand(cmd) {
					break Loop
				}
			}
		}
	}
}

func (a *App) handleCommand(cmd string) bool {
	switch cmd {
	case "exit":
		return true
	case "new":
		a.startTimer()
	default:
		userError(a.screen, "🤷 Неизвестная команда")
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
			userHint(a.screen, "✏️  Введите 'new' или 'exit'")

		}()
		a.timer.setTimer(1, 0)
		a.timer.run(a.screen)

	}()
}
