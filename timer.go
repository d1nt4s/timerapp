package main

import (
	"context"
	"fmt"
	"time"
    "github.com/mattn/go-runewidth"
	"github.com/gdamore/tcell/v2"
)

type Timer struct {
	seconds int
	minutes int
	control chan string
	status  Status
}

type Status int

const (
	Continue Status = iota
	Pause
	End
)

func (t *Timer) setup(sec int, min int) {
	t.seconds = sec
	t.minutes = min
}

func (t *Timer) decrementSec() {
	if t.seconds == 0 {
		if t.minutes == 0 {
			t.status = End
			return
		}
		t.minutes--
		t.seconds = 60
	}
	t.seconds--
	t.status = Continue
}

func (t *Timer) run(cancel context.CancelFunc, s tcell.Screen) {

	for {
		t.manage(s)

		if t.status == End {
		Drain:
			for {
				select {
				case <-t.control:
				default:
					break Drain
				}
			}
			fmt.Println() 
			cancel()
			return
		}

		if t.status == Pause {
			time.Sleep(200 * time.Millisecond)
			continue
		}

		t.decrementSec()
		drawRemainingTime(s, t.minutes, t.seconds, 0, tcell.StyleDefault.Foreground(tcell.ColorWhite))

		time.Sleep(time.Second)
	}
}


func (t *Timer) manage(screen tcell.Screen) {
	select {
	case cmd := <-t.control:
		switch cmd {
		case "stop":
			t.setup(0, 0)
			t.status = End
			drawMessage(screen, "Таймер остановлен", 4, tcell.StyleDefault.Foreground(tcell.ColorRed))
		case "reset":
			t.setup(0, 15)
			drawMessage(screen, "🔁 Таймер сброшен", 4, tcell.StyleDefault.Foreground(tcell.ColorRed))
		case "pause":
			t.status = Pause
			drawMessage(screen, "⏸ Таймер на паузе", 4, tcell.StyleDefault.Foreground(tcell.ColorRed))
		case "resume":
			t.status = Continue
			drawMessage(screen, "▶️ Таймер продолжается", 4, tcell.StyleDefault.Foreground(tcell.ColorRed))
		case "exit":
			t.status = End
		default:
			drawFormattedMessage(screen, 4, tcell.StyleDefault.Foreground(tcell.ColorYellow), "🤷 Неизвестная команда: %s", cmd)

		}
	default:
	}
}

func drawRemainingTime(s tcell.Screen, tMin, tSec int, y int, style tcell.Style) {
    // 1. Формируем строку
    msg := fmt.Sprintf("⏳ Осталось: %d мин %02d сек", tMin, tSec)

    // 2. Очищаем старую строку (на всякий случай)
    w, _ := s.Size()
    for x := 0; x < w; x++ {
        s.SetContent(x, y, ' ', nil, style)
    }

    // 3. Выводим посимвольно
    x := 0
    for _, ch := range msg {
        s.SetContent(x, y, ch, nil, style)
        x += runewidth.RuneWidth(ch)
    }

    // 4. Отображаем на экране
    s.Show()
}

func drawMessage(s tcell.Screen, msg string, y int, style tcell.Style) {
    // Очистим всю строку перед выводом, чтобы не было "хвостов"
    w, _ := s.Size()
    for x := 0; x < w; x++ {
        s.SetContent(x, y, ' ', nil, style)
    }

    // Выводим сообщение посимвольно
    x := 0
    for _, ch := range msg {
        s.SetContent(x, y, ch, nil, style)
        x += runewidth.RuneWidth(ch)
    }

    s.Show()
}

func drawFormattedMessage(s tcell.Screen, y int, style tcell.Style, format string, args ...interface{}) {
    // 1. Формируем строку, как в fmt.Printf
    msg := fmt.Sprintf(format, args...)

    // 2. Очищаем строку от старого текста
    w, _ := s.Size()
    for x := 0; x < w; x++ {
        s.SetContent(x, y, ' ', nil, style)
    }

    // 3. Выводим символы с учетом Unicode ширины
    x := 0
    for _, ch := range msg {
        s.SetContent(x, y, ch, nil, style)
        x += runewidth.RuneWidth(ch)
    }

    // 4. Показываем обновление
    s.Show()
}
