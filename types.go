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
	CmdStop   		Command = "stop"
	CmdReset  		Command = "reset"
	CmdPause  		Command = "pause"
	CmdResume 		Command = "resume"
	CmdExit   		Command = "exit"
	CmdSkip   		Command = "skip"
	CmdHelp   		Command = "help"
	CmdSetTimer     Command = "set_timer"
	CmdSetPause     Command = "set_pause"
	CmdSetInterval  Command = "set_interval"
	CmdSetLongPause Command = "set_longpause"
	CmdStart  		Command = "start"
	CmdSettings 	Command = "settings"
	CmdSnooze5Minutes	Command = "snooze5m"
	CmdSnooze10Minutes	Command = "snooze10m"
)

var commandMap = map[string]Command{
	"stop":   	CmdStop,
	"reset":  	CmdReset,
	"pause":  	CmdPause,
	"resume": 	CmdResume,
	"exit":   	CmdExit,
	"skip":   	CmdSkip,
	"help":   	CmdHelp,
	"start":  	CmdStart,
	"settings": CmdSettings,
	"s": 		CmdStart,
	"st": 		CmdStop,
	"res":		CmdReset,
	"p": 		CmdPause,
	"r": 		CmdResume,
	"e":		CmdExit,
	"sk":		CmdSkip,
	"h":		CmdHelp,
	"snooze5m": CmdSnooze5Minutes,
	"snooze10m": CmdSnooze10Minutes,
}
