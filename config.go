package main

import (
	"flag"
	"strings"
)

type Config struct {
	UseColors     bool
	OutputFormat  string
	MessageOutput string
	LogName       string
}

func (cfg *Config) InitFromFlags() {
	if logName := flag.String("logName", "Application", "Event log to attach to"); *logName != "" {
		cfg.LogName = *logName
	}
	if outputFormat := flag.String("outputFormat", "simple", "Output format (simple/rfc5424)"); *outputFormat != "" {
		cfg.OutputFormat = strings.ToLower(*outputFormat)
	}
	if messageOutput := flag.String("messageOutput", "singleLine", "Message output format (singeLine/singleLineTrim/full"); *messageOutput != "" {
		cfg.MessageOutput = strings.ToLower(*messageOutput)
	}
	useColors := flag.String("useColors", "false", "Whether to use colors in simple output format (false/true)")
	cfg.UseColors = strings.ToLower(*useColors) == "true"
}
