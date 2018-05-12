package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	VERSION                 = "0.1.4"
	DEFAULT_LOG_LINES       = 10
	DEFAULT_FOLLOW_STREAM   = false
	DEFAULT_FOLLOW_INTERVAL = 3
	DEFAULT_WAIT_PID        = -1
)

var (
	versionDump    bool
	numLines       int
	followEvents   bool
	followInterval int
	waitPid        int
)

func init() {
	flag.BoolVar(&versionDump, "version", false, "Display version and exit")
	flag.IntVar(&numLines, "n", DEFAULT_LOG_LINES, "Number of lines to dump")
	flag.BoolVar(&followEvents, "f", DEFAULT_FOLLOW_STREAM, "Follow the log group")
	flag.IntVar(&followInterval, "s", DEFAULT_FOLLOW_INTERVAL, "Interval (in seconds) to sleep during a follow")
	flag.IntVar(&waitPid, "p", DEFAULT_WAIT_PID, "with -f, terminate after process ID, PID dies")
	flag.Parse()
}

func dumpVersion() {
	fmt.Printf("cloudtail: version %s\n", VERSION)
	os.Exit(0)
}

func main() {
	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(1)
	} else if versionDump {
		dumpVersion()
	}

	session := InitSession()
	group := session.SearchLogGroups(os.Args[len(os.Args)-1])
	if !followEvents {
		session.DumpLogEvents(&group, numLines)
	} else {
		session.FollowLogEvents(&group, followInterval, waitPid)
	}
	os.Exit(0)
}
