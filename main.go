package main

func main() {
	var timer Timer
	timer.control = make(chan string)
	timer.setup(0, 15)
	go scan_command(timer.control);
	timer.status = Continue  
	timer.run()
}
