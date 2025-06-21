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

	userHello(a.screen, Msg_введите_команду_start_exit_help)

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
				userError(a.screen, "⭔ Команда \""+cleaned+"\" не распознана", true)
			}
		} else {
			if a.handleCommand(cmd) {
				break Loop
			}
		}
	}
}

func (a *App) startTimer() {

	a.acceptingTimerCommands = true

	go func() {
		defer func() {
			a.acceptingTimerCommands = false
			Debug("🟢 timer.run завершается")
		}()

		for {
			result := a.timer.Run(a.screen)

			switch result {
			case TimerExitApp:
				a.uiCommandCh <- "exit"
				return
			case TimerStopped:
				userNotice(a.screen, "⏱ Таймер остановлен", true)
				userHint(a.screen, "🐲 Введите 'start' для повтора или 'exit'", false)
				return
			case TimerFinished:
			    continue
			}
		}


	}()
}
