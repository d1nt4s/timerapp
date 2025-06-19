package main

import (
	"strconv"
	"strings"
)

func parseTimeFromSetCommand(input string) (int, int, bool) {
	// Пример: "set 2:15"
	if !strings.HasPrefix(input, "set ") {
		return 0, 0, false
	}
	parts := strings.Split(strings.TrimPrefix(input, "set "), ":")
	if len(parts) != 2 {
		return 0, 0, false
	}
	min, err1 := strconv.Atoi(parts[0])
	sec, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil || min < 0 || sec < 0 || sec > 59 {
		return 0, 0, false
	}

	return min, sec, true
}
