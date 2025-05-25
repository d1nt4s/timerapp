package main

import (
	"fmt"
	"time"
)

func main() {
	var timer Timer
	timer.setup(0, 15)
	for timer.decrementSec() != End {
		fmt.Printf("Minutes: %d, Seconds: %d \r", timer.minutes, timer.seconds)
		time.Sleep(time.Second)
	}
}
