package main

import (
	"flag"
	"strings"
)

//All functionality related to configuration

//Configuration struct holding command line parameters
type Config struct {
	UseColors     bool
	OutputFormat  string
	MessageOutput string
	LogName       string
}

//Initialize config object from command line parameters
func (cfg *Config) InitFromFlags() {
	logName := flag.String("logname", "Application", "Event log to attach to")
	outputFormat := flag.String("outfmt", "simple", "Output format (simple/rfc5424)")
	messageOutput := flag.String("msgout", "singleLine", "Message output format (singeLine/singleLineTrim/full")
	useColors := flag.String("colors", "false", "Whether to use colors in simple output format (false/true)")
	flag.Parse()

	if *logName != "" {
		cfg.LogName = *logName
	}
	if *outputFormat != "" {
		cfg.OutputFormat = strings.ToLower(*outputFormat)
	}
	if *messageOutput != "" {
		cfg.MessageOutput = strings.ToLower(*messageOutput)
	}

	cfg.UseColors = strings.ToLower(*useColors) == "true"
}
