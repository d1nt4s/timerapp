package main

func main() {

	app := NewApp()
	defer app.screen.Fini()

	app.Run()

	userNotice(app.screen, "👋 Программа завершена.")
}
