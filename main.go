package main

func main() {
	app := NewApp()
	defer app.screen.Fini()

	go scanCommand(app.screen, app.commandCh)

	app.Run()
}
