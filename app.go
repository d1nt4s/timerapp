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
	userNotice(a.screen, "⌨️  Введите команду (start / exit):")

	var timer Timer

	a.timer = &timer
	a.timer.control = make(chan string)	
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
		a.startTimer(1, 0)
	default:
		userError(a.screen, "🤷 Неизвестная команда")
	}

	return false
}

func (a *App) startTimer(min, sec int) {

	a.timer.setup(sec, min)
	a.timer.status = Continue
	a.acceptingTimerCommands = true

	go func() {
		defer func() {
			a.acceptingTimerCommands = false
			Debug("🟢 timer.run завершается")
			userNotice(a.screen, "⏱ Таймер завершён")
			userHint(a.screen, "✏️  Введите 'new' или 'exit'")

		}()
		a.timer.run(a.screen)

	}()
}
