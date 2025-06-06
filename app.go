package main

type App struct {
	timerDoneCh chan string
}

func NewApp() *App {

	return &App {
		timerDoneCh: make(chan string),
	}	
}