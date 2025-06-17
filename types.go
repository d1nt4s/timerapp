package main

type TimerStatus int
type TimerResult int
type Command string

const (
	Continued TimerStatus = iota
	Paused
	Stopped
	ExitApp
)

const (
	TimerStopped TimerResult = iota
	TimerExitApp
)

const (
	CmdStop   Command = "stop"
	CmdReset  Command = "reset"
	CmdPause  Command = "pause"
	CmdResume Command = "resume"
	CmdExit   Command = "exit"
)

var commandMap = map[string]Command{
	"stop":   CmdStop,
	"reset":  CmdReset,
	"pause":  CmdPause,
	"resume": CmdResume,
	"exit":   CmdExit,
}
