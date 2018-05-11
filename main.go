package main

import (
	"flag"
	"os"
)

const (
	DEFAULT_LOG_LINES     = 10
	DEFAULT_STREAM_EVENTS = false
)

var (
	numLines     int
	streamEvents bool
)

func init() {
	flag.IntVar(&numLines, "n", DEFAULT_LOG_LINES, "Number of lines to dump")
	flag.BoolVar(&streamEvents, "f", DEFAULT_STREAM_EVENTS, "Stream the log group") // to do
	flag.Parse()
}

func main() {
	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(1)
	}
	session := InitSession()
	group := session.SearchLogGroups(os.Args[len(os.Args)-1])
	streams := session.GetLogStreams(group)
	session.DumpLogEvents(group, streams, numLines)
	os.Exit(0)
}
