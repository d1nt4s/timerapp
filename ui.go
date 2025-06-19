package main

import (
	"fmt"
	"math"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

var bigFont = map[rune][]string{
	'0': {" ██████ ", "██    ██", "██  ██ ██", "██ ██  ██", "██    ██", " ██████ ", "        "},
	'1': {"   ██   ", " ████   ", "   ██   ", "   ██   ", "   ██   ", " ██████ ", "        "},
	'2': {" ██████ ", "██    ██", "     ██ ", "  ████  ", " ██     ", "████████", "        "},
	'3': {" ██████ ", "     ██ ", "  ████  ", "     ██ ", "██   ██ ", " █████  ", "        "},
	'4': {"   ███  ", "  █ ██  ", " ██ ██  ", "████████", "    ██  ", "   ████ ", "        "},
	'5': {"███████ ", "██      ", "██████  ", "     ██ ", "██   ██ ", " █████  ", "        "},
	'6': {"  █████ ", " ██     ", "██████  ", "██   ██ ", "██   ██ ", " █████  ", "        "},
	'7': {"████████", "    ██  ", "   ██   ", "  ██    ", "  ██    ", "  ██    ", "        "},
	'8': {" █████  ", "██   ██ ", " █████  ", "██   ██ ", "██   ██ ", " █████  ", "        "},
	'9': {" █████  ", "██   ██ ", "██   ██ ", " ██████ ", "     ██ ", " █████  ", "        "},
	':': {"        ", "   ██   ", "        ", "        ", "   ██   ", "        ", "        "},
}

func drawCenteredBigTimer(s tcell.Screen, min, sec int, style tcell.Style) {
	clearAllExceptInputLine(s)
	msg := fmt.Sprintf("%02d:%02d", min, sec)
	height := len(bigFont['0'])
	width := 0
	for _, ch := range msg {
		if lines, ok := bigFont[ch]; ok {
			width += runewidth.StringWidth(lines[0]) + 2
		}
	}

	scrWidth, scrHeight := s.Size()
	startX := (scrWidth - width) / 2
	startY := int(math.Max(float64((scrHeight-height)/2-2), 1))
	x := startX
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
		x += runewidth.StringWidth(lines[0]) + 2
	}

	drawProgressBar(s, min, sec, scrHeight-4)

	s.Show()
}

func drawProgressBar(s tcell.Screen, min, sec, y int) {
	total := 60*min + sec
	if total == 0 {
		total = 1
	}
	scrWidth, _ := s.Size()
	filled := int(float64(scrWidth) * (float64(total) / float64(60*min+60)))
	style := tcell.StyleDefault.Foreground(tcell.ColorLightSkyBlue).Background(tcell.ColorBlack)
	for x := 0; x < scrWidth; x++ {
		char := '░'
		if x < filled {
			char = '█'
		}
		s.SetContent(x, y, char, nil, style)
	}
}

func drawMessage(s tcell.Screen, msg string, y int, style tcell.Style) {
	w, _ := s.Size()
	x := (w - runewidth.StringWidth(msg)) / 2
	for i := 0; i < w; i++ {
		s.SetContent(i, y, ' ', nil, tcell.StyleDefault.Background(tcell.ColorReset))
	}
	for _, r := range msg {
		s.SetContent(x, y, r, nil, style)
		x += runewidth.RuneWidth(r)
	}
	s.Show()
}

func userNotice(s tcell.Screen, msg string) {
	drawMessage(s, msg, 18, tcell.StyleDefault.Foreground(tcell.ColorAqua).Bold(true))
}

func userHint(s tcell.Screen, msg string) {
	drawMessage(s, msg, 19, tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true))
}

func userError(s tcell.Screen, msg string) {
	drawMessage(s, msg, 20, tcell.StyleDefault.Foreground(tcell.ColorMaroon).Bold(true))
}

func writeToInputLine(screen tcell.Screen, buffer []rune) {
	width, height := screen.Size()
	style := tcell.StyleDefault.Foreground(tcell.ColorHotPink).Background(tcell.ColorBlack).Bold(true)
	prompt := "> "
	for x := 0; x < width; x++ {
		screen.SetContent(x, height-1, ' ', nil, style)
	}
	x := 0
	for _, r := range prompt {
		screen.SetContent(x, height-1, r, nil, style)
		x += runewidth.RuneWidth(r)
	}
	for _, r := range buffer {
		screen.SetContent(x, height-1, r, nil, style)
		x += runewidth.RuneWidth(r)
	}
	screen.Show()
}

func clearInputLine(screen tcell.Screen) {
	width, height := screen.Size()
	for x := 0; x < width; x++ {
		screen.SetContent(x, height-1, ' ', nil, tcell.StyleDefault)
	}
	screen.Show()
}

func clearUserLines(s tcell.Screen) {
	for y := 18; y <= 20; y++ {
		w, _ := s.Size()
		for x := 0; x < w; x++ {
			s.SetContent(x, y, ' ', nil, tcell.StyleDefault)
		}
	}
	s.Show()
}

func clearAllExceptInputLine(s tcell.Screen) {
	width, height := s.Size()
	style := tcell.StyleDefault

	for y := 0; y < height-1; y++ { // height-1 = строка ввода
		for x := 0; x < width; x++ {
			s.SetContent(x, y, ' ', nil, style)
		}
	}
}

