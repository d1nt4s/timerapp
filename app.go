package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
)

type AppMode int

const (
	SetupMode AppMode = iota
	ActiveMode
)

type App struct {
	screen                 tcell.Screen
	timer                  *Timer
	uiCommandCh            chan string
	mode 				   AppMode					
}

func NewApp() *App {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("ошибка создания экрана: %v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("ошибка инициализации экрана: %v", err)
	}
	screen.EnableMouse()

	drawButtons(screen, getButtons(screen, SetupMode), tcell.StyleDefault.Foreground(tcell.ColorAqua).Bold(true))

	return &App{
		screen:      screen,
		timer: 	     NewTimer(1, 0),
		uiCommandCh: make(chan string),
		mode: 		 SetupMode,
	}
}

func (a *App) Run() {
	defer func() {
		close(a.uiCommandCh) // Закрываем канал команд ТОЛЬКО здесь
		Debug("🟢 Основной канал сырых команд закрылся")
	}()

	userHello(a.screen, Msg_введите_команду_start_exit_help)
	// a.changeMode()

	// a.timer = NewTimer(1, 0)
	// a.mode = SetupMode

	// drawButtons(a.screen, getButtons(a.screen, a.mode), tcell.StyleDefault.Foreground(tcell.ColorAqua).Bold(true))

	go func() {
		defer func() {
			Debug("🟢 scanCommand завершается")
		}()
		scanCommand(a)
	}()

Loop:

	for cmd := range a.uiCommandCh {
		if parsed, cleaned, ok := ParseCommand(cmd); ok {
			if a.mode == ActiveMode {
				// drawButtons(a.screen, getButtons(a.screen, a.mode), tcell.StyleDefault.Foreground(tcell.ColorAqua).Bold(true))
				a.timer.control <- parsed
			} else {
				if a.handleCommand(parsed, cleaned) {
					break Loop
				}
			}
		} else {
			userError(a.screen, "⭔ Команда \""+cleaned+"\" не распознана", true)
		}
	}
}

func (a *App) startTimer() {

	a.changeMode()

	go func() {
		defer func() {
			a.changeMode()
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

func (a *App) changeMode() {
	if (a.mode == SetupMode) {
		a.mode = ActiveMode
	} else {
		a.mode = SetupMode
	}

	clearButtonLine(a.screen)
	drawButtons(a.screen, getButtons(a.screen, a.mode), tcell.StyleDefault.Foreground(tcell.ColorAqua).Bold(true))
}
