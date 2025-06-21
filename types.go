package main

type TimerStatus int
type TimerResult int
type TimerMode int
type Command string

const (
	Continued TimerStatus = iota
	Paused
	Finished
	Stopped
	ExitApp
)

const (
	TimerStopped TimerResult = iota
	TimerExitApp
	TimerFinished
)

const (
	Pomodoro TimerMode = iota
	Pause
	LongPause
)

const (
	CmdStop   Command = "stop"
	CmdReset  Command = "reset"
	CmdPause  Command = "pause"
	CmdResume Command = "resume"
	CmdExit   Command = "exit"
	CmdSkip   Command = "skip"
	CmdHelp   Command = "help"
)

var commandMap = map[string]Command{
	"stop":   CmdStop,
	"reset":  CmdReset,
	"pause":  CmdPause,
	"resume": CmdResume,
	"exit":   CmdExit,
	"skip":   CmdSkip,
	"help":   CmdHelp,
}
