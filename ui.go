package main

import (
	"fmt"
	"math"
	"strings"
	"time"

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
	// clearAllExceptMessagesAndInputLine(s)
	clearBigTimerArea(s)
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

func userHello(s tcell.Screen, msg string) {
	drawMessage(s, msg, 17, tcell.StyleDefault.Foreground(tcell.ColorAqua).Bold(true))
}

func userNotice(s tcell.Screen, msg string, toClear bool) {

	if toClear {
		drawMessageWithAutoClear(s, msg, 18, tcell.StyleDefault.Foreground(tcell.ColorAqua).Bold(true))
	} else {
		drawMessage(s, msg, 18, tcell.StyleDefault.Foreground(tcell.ColorAqua).Bold(true))
	}
}

func userHint(s tcell.Screen, msg string, toClear bool) {

	if toClear {
		drawMessageWithAutoClear(s, msg, 19, tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true))
	} else {
		drawMessage(s, msg, 19, tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true))
	}
}

func userError(s tcell.Screen, msg string, toClear bool) {

	if toClear {
		drawMessageWithAutoClear(s, msg, 20, tcell.StyleDefault.Foreground(tcell.ColorMaroon).Bold(true))
	} else {
		drawMessage(s, msg, 20, tcell.StyleDefault.Foreground(tcell.ColorMaroon).Bold(true))
	}
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

func clearAllExceptMessagesAndInputLine(s tcell.Screen) {
	width, height := s.Size()
	style := tcell.StyleDefault

	for y := 0; y < height-1; y++ {
		if y >= 18 && y <= 27 {
			continue // Не очищаем строки с сообщениями (18–27)
		}
		for x := 0; x < width; x++ {
			s.SetContent(x, y, ' ', nil, style)
		}
	}
}

func drawMessageWithAutoClear(s tcell.Screen, msg string, line int, style tcell.Style) {
	drawMessage(s, msg, line, style)

	go func() {
		time.Sleep(3 * time.Second)
		width, _ := s.Size()
		for x := 0; x < width; x++ {
			s.SetContent(x, line, ' ', nil, tcell.StyleDefault)
		}
		s.Show()
	}()
}

func drawLongNotice(s tcell.Screen, msg string) {

	style := tcell.StyleDefault.Foreground(tcell.ColorLightGrey).Bold(true)

	w, _ := s.Size()
	lines := wrapText(msg, w-10) // немного шире, но с полями

	startRow := 21
	for i, line := range lines {
		if i >= 6 {
			break
		}
		drawCenteredText(s, line, startRow+i, style)
	}
	s.Show()

	go func() {
		time.Sleep(15 * time.Second)
		for line := 21; line < 27; line++ {
			width, _ := s.Size()
			for x := range width {
				s.SetContent(x, line, ' ', nil, tcell.StyleDefault)
			}
		}

		s.Show()
	}()
}

func drawCenteredText(s tcell.Screen, text string, row int, style tcell.Style) {
	w, _ := s.Size()
	textLen := runewidth.StringWidth(text) // учитывает ширину символов
	startCol := (w - textLen) / 2
	col := startCol

	for _, r := range text {
		s.SetContent(col, row, r, nil, style)
		col += runewidth.RuneWidth(r)
	}
}

func wrapText(text string, maxWidth int) []string {
	var lines []string
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	current := words[0]
	for _, word := range words[1:] {
		if runewidth.StringWidth(current+" "+word) > maxWidth {
			lines = append(lines, current)
			current = word
		} else {
			current += " " + word
		}
	}
	lines = append(lines, current)
	return lines
}

func clearBigTimerArea(s tcell.Screen) {
	height := len(bigFont['0']) // 7 строк у цифры
	scrWidth, scrHeight := s.Size()
	startY := int(math.Max(float64((scrHeight-height)/2-2), 1)) // как в drawCenteredBigTimer

	for y := startY; y < startY+height; y++ {
		for x := 0; x < scrWidth; x++ {
			s.SetContent(x, y, ' ', nil, tcell.StyleDefault)
		}
	}

	s.Show()
}
