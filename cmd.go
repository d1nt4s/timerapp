package main

import (
	"context"
	"fmt"
	"strings"

	// "sync"

	"github.com/chzyer/readline"
)

func createReadline() (*readline.Instance, error) {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "> ",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		return nil, err
	}
	return rl, nil
}

// func scan_command(rl *readline.Instance, control chan string, done <-chan bool, wg *sync.WaitGroup) {
func scan_command(ctx context.Context, rl *readline.Instance, control chan string) {
	// defer wg.Done()

	for {
		// Создаём select, чтобы слушать и ввод, и завершение
		select {
		case <-ctx.Done():
			fmt.Println("scancommand gets")
			return // безопасный выход
		default:
			line, err := rl.Readline()
			if err != nil {
				return // ^D или закрытие
			}
			cmd := strings.ToLower(strings.TrimSpace(line))
			fmt.Println("⏳ Перед отправкой команды в канал")
			control <- cmd
			fmt.Println("⏳ Перед отправкой команды в канал")
			if cmd == "exit" {
				return
			}
		}
	}
}

