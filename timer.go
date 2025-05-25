package main

type Timer struct {
	seconds int
	minutes int
}

type Status int

const (
	Continue Status = iota
	End
)

func (t *Timer) setup(sec int, min int) {
	t.seconds = sec
	t.minutes = min
}

func (t *Timer) decrementSec() Status {
	if t.seconds == 0 {
		if t.minutes == 0 {
			return End
		}
		t.minutes--
		t.seconds = 60
	}
	t.seconds--
	return Continue
}
