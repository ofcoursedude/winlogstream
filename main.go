package main

import (
	"fmt"
	winlog "github.com/ofcoursedude/gowinlog"
	"strconv"
	"strings"
	"time"
)

func main() {
	fmt.Println("Starting...")
	watcher, err := winlog.NewWinLogWatcher()
	if err != nil {
		fmt.Printf("Couldn't create watcher: %v\n", err)
		return
	}

	// Recieve any future messages on the Application channel
	// "*" doesn't filter by any fields of the event
	watcher.SubscribeFromNow("Application", "*")
	for {
		select {
		case evt := <-watcher.Event():
			// Print the event struct
			// fmt.Printf("\nEvent: %v\n", evt)
			// or print basic output
			// fmt.Printf("\n%s: %s: %s\n", evt.LevelText, evt.ProviderName, evt.Msg)
			fmt.Println(ToRfc5424(evt))
		case err := <-watcher.Error():
			fmt.Printf("\nError: %v\n\n", err)
		default:
			// If no event is waiting, need to wait or do something else, otherwise
			// the the app fails on deadlock.
			<-time.After(1 * time.Millisecond)
		}
	}
}

func ToRfc5424(evt *winlog.WinLogEvent) string{
	output:= []string{
		"<34>1",
		evt.Created.Format(time.RFC3339),
		evt.ComputerName,
		evt.ProviderName,
		strconv.FormatInt(int64(evt.ProcessId), 10),
		strconv.FormatInt(int64(evt.EventId), 10),
		evt.Msg,
	}
	return strings.Join(output, " ")
	// return fmt.Sprint ("<34>1",  evt.Created.Format(time.RFC3339), evt.ComputerName, evt.ProviderName, evt.ProcessId, evt.EventId, evt.Msg)
}
