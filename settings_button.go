package main

import (
	"strconv"
	"strings"

	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func runWithTviewApp(app *App, tviewApp *tview.Application, content tview.Primitive) {
	app.screen.Fini()

	tviewApp.SetRoot(content, true)
	if err := tviewApp.Run(); err != nil {
		log.Printf("Ошибка при запуске tview: %v", err)
	}

	if err := app.screen.Init(); err != nil {
		log.Fatalf("Ошибка повторной инициализации tcell: %v", err)
	}

	clearAllExceptMessagesAndInputLine(app.screen)
	drawButtons(app.screen, getButtons(app.screen, app.mode), tcell.StyleDefault)
	userNotice(app.screen, "✅ Настройки применены", true)
}

func showSettingsModal(app *App) {
	tviewApp := tview.NewApplication()

	form := tview.NewForm().
		AddInputField("Pomodoro (mm:ss)", "", 7, nil, nil).
		AddInputField("Pause (mm:ss)", "", 7, nil, nil).
		AddInputField("Long Pause (mm:ss)", "", 7, nil, nil).
		AddInputField("Interval (N)", "", 5, nil, nil)

	form.AddButton("Сохранить", func() {
		timer := form.GetFormItemByLabel("Pomodoro (mm:ss)").(*tview.InputField).GetText()
		pause := form.GetFormItemByLabel("Pause (mm:ss)").(*tview.InputField).GetText()
		longPause := form.GetFormItemByLabel("Long Pause (mm:ss)").(*tview.InputField).GetText()
		interval := form.GetFormItemByLabel("Interval (N)").(*tview.InputField).GetText()

		if !isValidTimeFormat(timer) || !isValidTimeFormat(pause) || !isValidTimeFormat(longPause) {
			form.SetTitle("⛔ Формат времени должен быть mm:ss").SetTitleColor(tcell.ColorRed)
			return
		}
		if val, err := strconv.Atoi(interval); err != nil || val < 0 {
			form.SetTitle("⛔ Интервал должен быть числом ≥ 0").SetTitleColor(tcell.ColorRed)
			return
		}

		commands := []string{
			"set_timer " + timer,
			"set_pause " + pause,
			"set_longpause " + longPause,
			"set_interval " + interval,
		}

		for _, cmd := range commands {
			command, raw, ok := ParseCommand(cmd)
			if ok {
				app.handleCommand(command, raw)
			}
		}

		tviewApp.Stop()
	})

	form.AddButton("Отмена", func() {
		tviewApp.Stop()
	})

	form.SetBorder(true).
		SetTitle("Настройки таймера").
		SetTitleAlign(tview.AlignLeft).
		SetTitleColor(tcell.ColorWhite)

	runWithTviewApp(app, tviewApp, form)
}



func isValidTimeFormat(input string) bool {
	parts := strings.Split(input, ":")
	if len(parts) != 2 {
		return false
	}
	min, err1 := strconv.Atoi(parts[0])
	sec, err2 := strconv.Atoi(parts[1])
	return err1 == nil && err2 == nil && min >= 0 && sec >= 0 && sec < 60
}
