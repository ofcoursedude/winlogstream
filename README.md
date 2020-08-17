# winlogstream

Winlog stream is a tool to stream Windows Event Log events to console.

## to run
Run the exe file to see usage. 

Please note log names are case sensitive and correspond to the field "Full Name" you can see when displaying the event log properties, not the descriptive name you see in the event log tree.

## from source

The tool is written in Go. Install go 1.14+ and simply run 

`go build .` 

to build the exe file, or

`go install`

to install the tool in your path. If you don't have go but have docker on your machine, you can also compile with the command:

`docker run --rm -v <path_to_the_source_files>:/usr/src/winlogstream -w /usr/src/winlogstream -e GOOS=windows -e GOARCH=amd64 golang:1.14 go build -v`