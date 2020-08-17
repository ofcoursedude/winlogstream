//Command line tool for hooking into the Windows Event Log and streaming messages as they come in
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	winlog "github.com/ofcoursedude/gowinlog"

	"gitlab.certicon.cz/tools/winlogstream/colors"
)

//winlogstream is a line tool for hooking into the Windows Event Log and streaming messages as they come in

//In-app configuration object
var config Config

//Version string
var Version string = "dev build"

func main() {
	fmt.Println("Welcome to winlogstream")
	fmt.Println("Version:", Version)
	fmt.Println("Usage:")
	config = Config{}
	config.InitFromFlags()
	flag.PrintDefaults()

	fmt.Println("Starting...")

	var outputFormatFunc func(evt *winlog.WinLogEvent, msgFormat func(msg string) string) string
	var msgOutFunc func(msg string) string

	switch config.MessageOutput {
	case "full":
		msgOutFunc = func(msg string) string {
			return msg
		}
	case "singleline":
		msgOutFunc = singleLine
	case "singlelinetrim":
		msgOutFunc = singleLineTrim
	default:
		log.Fatal("Invalid Message Format")
	}

	switch config.OutputFormat {
	case "simple":
		outputFormatFunc = toSimple
	case "rfc5424":
		outputFormatFunc = toRfc5424
	default:
		log.Fatal("Invalid output format")
	}

	shutdowner := make(chan bool)
	go func(sig chan bool) {
		// when we exit, signal it's done
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
		err = watcher.SubscribeFromNow(config.LogName, "*")
		if err != nil {
			log.Fatal(fmt.Sprint("Can not subscribe to log ", config.LogName))
		}
		defer watcher.Shutdown()
	EventCollectionLoop:
		for {
			select {
			case evt := <-watcher.Event():
				if evt.Level <= config.Severity {
					fmt.Println(outputFormatFunc(evt, msgOutFunc))
				}
			case err := <-watcher.Error():
				fmt.Printf("\nError: %v\n\n", err)
				// Waiting for graceful shutdown signal is good enough to omit
				// the 'default' block
			case <-sig:
				break EventCollectionLoop
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
	return replaceMulti(msg, []string{"\r", "\n"}, " ")
}

func singleLineTrim(msg string) string {
	return strings.Split(strings.Replace(msg, "\r", "", 1), "\r\n")[0]
}

func toSimple(evt *winlog.WinLogEvent, msgFormat func(msg string) string) string {
	level := eventLevel(evt.Level)
	var levelMsg string
	if config.UseColors {
		levelMsg = fmt.Sprint(level.Color(), "[", level.String(), "]", colors.Reset)
	} else {
		levelMsg = fmt.Sprint("[", eventLevel(evt.Level).String(), "]")
	}
	output := []string{
		evt.Created.Format(time.RFC3339),
		levelMsg,
		strings.ReplaceAll(evt.ProviderName, " ", "_"),
		msgFormat(evt.Msg),
	}
	return strings.Join(output, " ")
}

func toRfc5424(evt *winlog.WinLogEvent, msgFormat func(msg string) string) string {
	output := []string{
		"<34>1",
		evt.Created.Format(time.RFC3339),
		fmt.Sprint("[", eventLevel(evt.Level).String(), "]"),
		evt.ComputerName,
		strings.ReplaceAll(evt.ProviderName, " ", "_"),
		strconv.FormatInt(int64(evt.ProcessId), 10),
		strconv.FormatInt(int64(evt.EventId), 10),
		msgFormat(evt.Msg),
	}
	return strings.Join(output, " ")
}

func replaceMulti(source string, toReplace []string, replacement string) string {
	toReturn := source
	for _, item := range toReplace {
		toReturn = strings.ReplaceAll(toReturn, item, replacement)
	}
	return toReturn
}
