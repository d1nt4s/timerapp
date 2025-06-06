package main

import (
	"context"
	"sync"
	"log"
	"github.com/gdamore/tcell/v2"
)

type App struct {
	screen tcell.Screen
	timer *Timer
	commandCh chan string
	quitCh chan struct{}
}

func NewApp() *App {
	screen, err := tcell.NewScreen()
	if err != nil {
        log.Fatalf("ошибка создания экрана: %v", err)
	}
	if err := screen.Init(); err != nil {
        log.Fatalf("ошибка инициализации экрана: %v", err)
	}
	return &App {
		screen: screen,
		commandCh: make(chan string),
		quitCh: make(chan struct{}),
	}	
}

func (a *App) Run() {
	drawMessage(a.screen, "⌨️  Введите команду (start / exit):", 0, tcell.StyleDefault)
	
	for {
		select {
		case cmd := <-a.commandCh:
			a.handleCommand(cmd)
		case <-a.quitCh:
			return
		}
	}
}

func (a *App) handleCommand(cmd string) {
	switch cmd {
	case "exit":
		drawMessage(a.screen, "👋 До свидания.", 1, tcell.StyleDefault.Foreground(tcell.ColorRed))
		close(a.quitCh)

	case "start":
		a.startTimer(0, 10)

	default:
		drawMessage(a.screen, "🤷 Неизвестная команда", 1, tcell.StyleDefault.Foreground(tcell.ColorOrange))
	}
}

func (a *App) startTimer(min, sec int) {
	drawMessage(a.screen, "⏳ Запуск таймера...", 1, tcell.StyleDefault.Foreground(tcell.ColorGreen))

	_, cancel := context.WithCancel(context.Background())

	timer := NewTimer(min, sec)
	timer.control = make(chan string)
	timer.status = Continue
	a.timer = timer

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		timer.Run(cancel, a.screen)
	}()
	wg.Wait()

	drawMessage(a.screen, "🟢 Таймер завершился", 2, tcell.StyleDefault.Foreground(tcell.ColorYellow))
}
