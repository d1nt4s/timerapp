package main

type App struct {
	uiCommandCh chan string
}

func NewApp() *App {

	return &App {
		uiCommandCh: make(chan string),
	}	
}