package main

import (
	"strings"

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

func scan_command(rl *readline.Instance, control chan string) {

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
