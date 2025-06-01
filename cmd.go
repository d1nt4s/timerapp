package main

import (
	"strings"

	"github.com/chzyer/readline"
)

func scan_command(control chan string) {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "> ",
		HistoryFile:     "/tmp/readline.tmp",
		InterruptPrompt: "^C",
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			break // например, ^D
		}
		cmd := strings.ToLower(strings.TrimSpace(line))
		control <- cmd
		if cmd == "exit" {
			break
		}
	}

	// scanner := bufio.NewScanner(os.Stdin)
	// for scanner.Scan() {
	// 	cmd := strings.ToLower(scanner.Text())
	// 	control <- cmd
	// 	if cmd == "exit" {
	// 		break
	// 	}
	// }
}
