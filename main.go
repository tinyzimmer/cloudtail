package main

import (
	"flag"
	"os"
)

const (
	DEFAULT_LOG_LINES       = 10
	DEFAULT_FOLLOW_STREAM   = false
	DEFAULT_FOLLOW_INTERVAL = 3
)

var (
	numLines       int
	followEvents   bool
	followInterval int
)

func init() {
	flag.IntVar(&numLines, "n", DEFAULT_LOG_LINES, "Number of lines to dump")
	flag.BoolVar(&followEvents, "f", DEFAULT_FOLLOW_STREAM, "Follow the log group")
	flag.IntVar(&followInterval, "s", DEFAULT_FOLLOW_INTERVAL, "Interval (in seconds) to sleep during a follow")
	flag.Parse()
}

func main() {
	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(1)
	}
	session := InitSession()
	group := session.SearchLogGroups(os.Args[len(os.Args)-1])
	if !followEvents {
		session.DumpLogEvents(&group, numLines)
	} else {
		session.FollowLogEvents(&group, followInterval)
	}
	os.Exit(0)
}
