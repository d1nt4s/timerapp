package main

import "fmt"

func main() {
	var timer Timer
	rl, err := createReadline()
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	timer.control = make(chan string)
	timer.done = make(chan bool)
	timer.setup(0, 15)
	go scan_command(rl, timer.control)
	timer.status = Continue
	timer.run(rl)

	<-timer.done
	fmt.Println("👋 Программа завершена.")
}
