package main

import (
	"strconv"
	"strings"
)

func parseTimeFromSetCommand(input, prefix string) (int, int, bool) {
	trimmed := strings.TrimPrefix(input, prefix+" ")
	parts := strings.Split(trimmed, ":")
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

