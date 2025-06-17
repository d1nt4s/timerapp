package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

var bigFont = map[rune][]string{
	'0': {" ███ ", "█   █", "█   █", "█   █", " ███ "},
	'1': {"  █  ", " ██  ", "  █  ", "  █  ", "█████"},
	'2': {"████ ", "    █", " ███ ", "█    ", "█████"},
	'3': {"████ ", "    █", " ███ ", "    █", "████ "},
	'4': {"█  █ ", "█  █ ", "█████", "   █ ", "   █ "},
	'5': {"█████", "█    ", "████ ", "    █", "████ "},
	'6': {" ███ ", "█    ", "████ ", "█   █", " ███ "},
	'7': {"█████", "    █", "   █ ", "  █  ", "  █  "},
	'8': {" ███ ", "█   █", " ███ ", "█   █", " ███ "},
	'9': {" ███ ", "█   █", " ████", "    █", " ███ "},
	':': {"     ", "  █  ", "     ", "  █  ", "     "},
}

// 🧼 Утилита: очистка строки
func clearLine(s tcell.Screen, y int, style tcell.Style) {
	w, _ := s.Size()
	for x := 0; x < w; x++ {
		s.SetContent(x, y, ' ', nil, style)
	}
}

// 🧼 Утилита: вывод текста по координатам
func printAt(s tcell.Screen, x, y int, msg string, style tcell.Style) {
	for _, r := range msg {
		s.SetContent(x, y, r, nil, style)
		x += runewidth.RuneWidth(r)
	}
}

// ✅ Простой вывод
func drawMessage(s tcell.Screen, msg string, y int, style tcell.Style) {
	clearLine(s, y, style)
	printAt(s, 0, y, msg, style)
	s.Show()
}

// ✅ Форматированный вывод
func drawFormattedMessage(s tcell.Screen, y int, style tcell.Style, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	drawMessage(s, msg, y, style)
}

// ✅ Вывод больших цифр
func drawBigTimer(s tcell.Screen, min, sec int, startY int, style tcell.Style) {
	msg := fmt.Sprintf("%02d:%02d", min, sec)
	height := 5
	w, _ := s.Size()

	// очистка всей области вывода
	for y := 0; y < height; y++ {
		for x := 0; x < w; x++ {
			s.SetContent(x, startY+y, ' ', nil, style)
		}
	}

	x := 0
	for _, ch := range msg {
		lines, ok := bigFont[ch]
		if !ok {
			x += 6 // просто пропускаем
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

func userNotice(s tcell.Screen, msg string) {
	clearUserLines(s)
	drawMessage(s, msg, 7, tcell.StyleDefault.Foreground(tcell.ColorWhite))
}

func userHint(s tcell.Screen, msg string) {
	clearUserLines(s)
	drawMessage(s, msg, 8, tcell.StyleDefault.Foreground(tcell.ColorYellow))
}

func userError(s tcell.Screen, msg string) {
	clearUserLines(s)
	drawMessage(s, msg, 9, tcell.StyleDefault.Foreground(tcell.ColorRed))
}

func writeToInputLine(screen tcell.Screen, buffer []rune) {
	width, height := screen.Size()
	for x := 0; x < width; x++ {
		screen.SetContent(x, height-1, ' ', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
	}
	for i, r := range buffer {
		screen.SetContent(i, height-1, r, nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
	}
	screen.Show()
}

func clearInputLine(screen tcell.Screen) {
	width, height := screen.Size()
	for x := 0; x < width; x++ {
		screen.SetContent(x, height-1, ' ', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
	}
	screen.Show()
}

func clearUserLines(s tcell.Screen) {
	for y := 7; y <= 9; y++ {
		clearLine(s, y, tcell.StyleDefault)
	}
	s.Show()
}
