package main

type eventLevel int

const (
	deadEvent eventLevel = iota
	criticalEvent
	errorEvent
	warningEvent
	informationEvent
	verboseEvent
)

func (e eventLevel) String() string {
	return [...]string{
		"DEAD",
		"CRIT",
		"ERROR",
		"WARN",
		"INFO",
		"VERB",
	}[e]
}

func (e eventLevel) Color() string {
	return [...]string{
		"\033[31m",
		"\033[31m",
		"\033[31m",
		"\033[33m",
		"\033[32m",
		"\033[34m",
	}[e]
}

const resetColor = "\033[0m"
