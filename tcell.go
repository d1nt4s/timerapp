package main

import (
	"github.com/mattn/go-runewidth"
	"github.com/gdamore/tcell/v2"
	"fmt"
)

// func drawRemainingTime(s tcell.Screen, tMin, tSec int, y int, style tcell.Style) {
//     // 1. Формируем строку
//     msg := fmt.Sprintf("⏳ Осталось: %d мин %02d сек", tMin, tSec)

//     // 2. Очищаем старую строку (на всякий случай)
//     w, _ := s.Size()
//     for x := 0; x < w; x++ {
//         s.SetContent(x, y, ' ', nil, style)
//     }

//     // 3. Выводим посимвольно
//     x := 0
//     for _, ch := range msg {
//         s.SetContent(x, y, ch, nil, style)
//         x += runewidth.RuneWidth(ch)
//     }

//     // 4. Отображаем на экране
//     s.Show()
// }

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

var bigFont = map[rune][]string{
    '0': {
        " ███ ",
        "█   █",
        "█   █",
        "█   █",
        " ███ ",
    },
    '1': {
        "  █  ",
        " ██  ",
        "  █  ",
        "  █  ",
        "█████",
    },
    '2': {
        "████ ",
        "    █",
        " ███ ",
        "█    ",
        "█████",
    },
    '3': {
        "████ ",
        "    █",
        " ███ ",
        "    █",
        "████ ",
    },
    '4': {
        "█  █ ",
        "█  █ ",
        "█████",
        "   █ ",
        "   █ ",
    },
    '5': {
        "█████",
        "█    ",
        "████ ",
        "    █",
        "████ ",
    },
    '6': {
        " ███ ",
        "█    ",
        "████ ",
        "█   █",
        " ███ ",
    },
    '7': {
        "█████",
        "    █",
        "   █ ",
        "  █  ",
        "  █  ",
    },
    '8': {
        " ███ ",
        "█   █",
        " ███ ",
        "█   █",
        " ███ ",
    },
    '9': {
        " ███ ",
        "█   █",
        " ████",
        "    █",
        " ███ ",
    },
    ':': {
        "     ",
        "  █  ",
        "     ",
        "  █  ",
        "     ",
    },
    '⏳': {
        "⏳⏳⏳⏳⏳",
        "  ⏳  ⏳  ",
        "   ⏳⏳  ",
        "  ⏳  ⏳  ",
        "⏳⏳⏳⏳⏳",
    },
}

func drawBigTimer(s tcell.Screen, tMin, tSec int, startY int, style tcell.Style) {
    msg := fmt.Sprintf(" %02d:%02d", tMin, tSec)

    height := 5
    w, _ := s.Size()

    // Clear display area
    for y := 0; y < height; y++ {
        for x := 0; x < w; x++ {
            s.SetContent(x, startY+y, ' ', nil, style)
        }
    }

    // Draw message line by line
    x := 0
    for _, ch := range msg {
        lines, ok := bigFont[ch]
        if !ok {
            x += 6 // placeholder spacing
            continue
        }

        for dy, line := range lines {
			dx := 0
			for _, r := range line {
				s.SetContent(x+dx, startY+dy, r, nil, style)
				dx += runewidth.RuneWidth(r)
			}
        }
		x += runewidth.StringWidth(lines[0]) + 1
    }

    s.Show()
}