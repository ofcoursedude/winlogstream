package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	winlog "github.com/ofcoursedude/gowinlog"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load()
	/*
		Valid env:
		FORMAT=simple/rfc5424
		MSGOUT=full/singleLine/singleLineTrim
	*/
}

func main() {
	fmt.Println("Starting...")
	outputFormat := strings.ToLower(os.Getenv("FORMAT"))
	msgOut := strings.ToLower(os.Getenv("MSGOUT"))

	var outputFormatFunc func(evt *winlog.WinLogEvent, msgFormat func(msg string) string) string
	var msgOutFunc func(msg string) string

	switch msgOut {
	case "full":
		msgOutFunc = func(msg string) string { return msg }
	case "singleline":
		msgOutFunc = singleLine
	case "singlelinetrim":
		msgOutFunc = singleLineTrim
	default:
		log.Panic("Unrecognized msg output format")
	}

	switch outputFormat {
	case "simple":
		outputFormatFunc = ToSimple
	case "rfc5424":
		outputFormatFunc = ToRfc5424
	default:
		log.Panic("Unrecognized output format")
	}

	shutdowner := make(chan bool)
	go func(sig chan bool) {
		defer func() {
			sig <- true
		}()
		watcher, err := winlog.NewWinLogWatcher()
		if err != nil {
			fmt.Printf("Couldn't create watcher: %v\n", err)
			return
		}

		// Recieve any future messages on the Application channel
		// "*" doesn't filter by any fields of the event
		watcher.SubscribeFromNow("System", "*")
		defer watcher.Shutdown()
	EventLoop:
		for {
			select {
			case evt := <-watcher.Event():
				// Print the event struct
				// fmt.Printf("\nEvent: %v\n", evt)
				// or print basic output
				// fmt.Printf("\n%s: %s: %s\n", evt.LevelText, evt.ProviderName, evt.Msg)
				fmt.Println(outputFormatFunc(evt, msgOutFunc))
			case err := <-watcher.Error():
				fmt.Printf("\nError: %v\n\n", err)
				//Waiting for graceful shutdown signal is good enough to omit
				//the 'default' block
			case <-sig:
				break EventLoop
				/* default:
				// If no event is waiting, need to wait or do something else, otherwise
				// the the app fails on deadlock.
				<-time.After(1 * time.Millisecond) */
			}
		}
	}(shutdowner)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	fmt.Println("Attempting graceful shutdown")
	signal.Stop(ch)
	shutdowner <- true
	<-shutdowner
	fmt.Println("Finished")
}

func singleLine(msg string) string {
	return strings.ReplaceAll(msg, "\n", " ")
}

func singleLineTrim(msg string) string {
	return strings.Split(msg, "\n")[0]
}

func Parse(
	evt *winlog.WinLogEvent,
	formatFunc func(evt *winlog.WinLogEvent, msgFormat func(msg string) string) []string,
	msgFormat func(msg string) string) string {
	output := formatFunc(evt, msgFormat)
	return strings.Join(output, " ")
}

func ToSimple(evt *winlog.WinLogEvent, msgFormat func(msg string) string) string {
	output := []string{
		evt.Created.Format(time.RFC3339),
		fmt.Sprint("[", EventLevel(evt.Level).String(), "]"),
		evt.ProviderName,
		evt.Msg,
	}
	return strings.Join(output, " ")
	// return fmt.Sprint ("<34>1",  evt.Created.Format(time.RFC3339), evt.ComputerName, evt.ProviderName, evt.ProcessId, evt.EventId, evt.Msg)
}

func ToRfc5424(evt *winlog.WinLogEvent, msgFormat func(msg string) string) string {
	output := []string{
		"<34>1",
		evt.Created.Format(time.RFC3339),
		fmt.Sprint("[", EventLevel(evt.Level).String(), "]"),
		evt.ComputerName,
		evt.ProviderName,
		strconv.FormatInt(int64(evt.ProcessId), 10),
		strconv.FormatInt(int64(evt.EventId), 10),
		evt.Msg,
	}
	return strings.Join(output, " ")
	// return fmt.Sprint ("<34>1",  evt.Created.Format(time.RFC3339), evt.ComputerName, evt.ProviderName, evt.ProcessId, evt.EventId, evt.Msg)
}

type EventLevel int

const (
	Nuclear EventLevel = iota
	Critical
	Error
	Warning
	Information
	Verbose
)

func (e EventLevel) String() string {
	return [...]string{
		"DEAD",
		"CRIT",
		"ERROR",
		"WARN",
		"INFO",
		"VERB",
	}[e]
}
