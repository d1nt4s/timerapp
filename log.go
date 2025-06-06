package main

import (
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

var debugLog *log.Logger

func init() {
	_ = os.Mkdir("logs", 0755)

	file, err := os.Create("logs/debug.log") // Перезаписываем каждый раз
	if err != nil {
		panic("не удалось создать debug.log: " + err.Error())
	}

	debugLog = log.New(file, "", 0) // Отключаем стандартные флаги, кастом ниже
}

// Debug пишет лог с короткой датой, миллисекундами и именем функции
func Debug(msg string) {
	pc, file, line, ok := runtime.Caller(1)
	funcName := "unknown"
	if ok {
		if fn := runtime.FuncForPC(pc); fn != nil {
			funcName = shortFuncName(fn.Name())
		}
	}

	now := time.Now().Format("01-02 15:04:05.000") // MM-DD HH:MM:SS.mmm
	shortFile := shortFileName(file)

	debugLog.Printf("[%s %s:%d %s] %s", now, shortFile, line, funcName, msg)
}

func shortFuncName(full string) string {
	parts := strings.Split(full, "/")
	last := parts[len(parts)-1]
	return strings.ReplaceAll(last, "main.", "")
}

func shortFileName(full string) string {
	parts := strings.Split(full, "/")
	return parts[len(parts)-1]
}
