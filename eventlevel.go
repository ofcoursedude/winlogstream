package main

import (
	"gitlab.certicon.cz/tools/winlogstream/colors"
)

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
		colors.Red,
		colors.Red,
		colors.Red,
		colors.Yellow,
		colors.Green,
		colors.Blue,
	}[e]
}
