package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

var bigFont = map[rune][]string{
	'0': {" ██████ ", "██    ██", "██    ██", "██    ██", " ██████ ", "        ", "        "},
	'1': {"   ██   ", " ████   ", "   ██   ", "   ██   ", " ██████ ", "        ", "        "},
	'2': {" ██████ ", "     ██ ", " ██████ ", "██      ", "████████", "        ", "        "},
	'3': {" ██████ ", "     ██ ", " █████  ", "     ██ ", " ██████ ", "        ", "        "},
	'4': {"██   ██ ", "██   ██ ", "████████", "     ██ ", "     ██ ", "        ", "        "},
	'5': {"████████", "██      ", "██████  ", "     ██ ", "██████  ", "        ", "        "},
	'6': {" ██████ ", "██      ", "██████  ", "██   ██ ", " █████  ", "        ", "        "},
	'7': {"████████", "     ██ ", "    ██  ", "   ██   ", "  ██    ", "        ", "        "},
	'8': {" █████  ", "██   ██ ", " █████  ", "██   ██ ", " █████  ", "        ", "        "},
	'9': {" █████  ", "██   ██ ", " ██████ ", "     ██ ", " █████  ", "        ", "        "},
	':': {"        ", "   ██   ", "        ", "   ██   ", "        ", "        ", "        "},
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
	height := 7
	w, _ := s.Size()
	totalWidth := 0
	for _, ch := range msg {
		if lines, ok := bigFont[ch]; ok {
			totalWidth += runewidth.StringWidth(lines[0]) + 1
		} else {
			totalWidth += 8
		}
	}
	x := (w - totalWidth) / 2

	for y := 0; y < height; y++ {
		clearLine(s, startY+y, style)
	}

	for _, ch := range msg {
		lines, ok := bigFont[ch]
		if !ok {
			x += 8
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
	drawCenteredMessage(s, msg, 15, tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlue).Bold(true))
}

func userHint(s tcell.Screen, msg string) {
	clearUserLines(s)
	drawCenteredMessage(s, msg, 16, tcell.StyleDefault.Foreground(tcell.ColorYellow).Background(tcell.ColorBlack).Bold(true))
}

func userError(s tcell.Screen, msg string) {
	clearUserLines(s)
	drawCenteredMessage(s, msg, 17, tcell.StyleDefault.Foreground(tcell.ColorRed).Background(tcell.ColorBlue).Bold(true))
}

func drawCenteredMessage(s tcell.Screen, msg string, y int, style tcell.Style) {
	w, _ := s.Size()
	x := (w - runewidth.StringWidth(msg)) / 2
	clearLine(s, y, style)
	printAt(s, x, y, msg, style)
	s.Show()
}

func writeToInputLine(screen tcell.Screen, buffer []rune) {
	_, height := screen.Size()
	y := height - 3
	clearLine(screen, y, tcell.StyleDefault.Background(tcell.ColorBlack))
	for i, r := range buffer {
		screen.SetContent(i, y, r, nil, tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true))
	}
	screen.Show()
}

func clearInputLine(screen tcell.Screen) {
	width, height := screen.Size()
	y := height - 3
	for x := 0; x < width; x++ {
		screen.SetContent(x, y, ' ', nil, tcell.StyleDefault.Background(tcell.ColorBlack))
	}
	screen.Show()
}

func clearUserLines(s tcell.Screen) {
	for y := 15; y <= 17; y++ {
		clearLine(s, y, tcell.StyleDefault)
	}
	s.Show()
}
