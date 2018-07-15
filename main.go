package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	VERSION                 = "0.1.6"
	DEFAULT_LOG_LINES       = 10
	DEFAULT_FOLLOW_STREAM   = false
	DEFAULT_FOLLOW_INTERVAL = 3
	DEFAULT_WAIT_PID        = -1
	DEFAULT_VERBOSE_OUTPUT  = false
	DEFAULT_HIDE_METADATA   = false
	DEFAULT_LIST_GROUPS     = false
)

var (
	versionDump    bool
	followEvents   bool
	verboseOutput  bool
	hideMetadata   bool
	listGroups     bool
	numLines       int
	followInterval int
	waitPid        int
)

func init() {
	flag.BoolVar(&versionDump, "version", false, "Display version and exit")
	flag.IntVar(&numLines, "n", DEFAULT_LOG_LINES, "Number of lines to dump")
	flag.BoolVar(&followEvents, "f", DEFAULT_FOLLOW_STREAM, "Follow the log group")
	flag.IntVar(&followInterval, "s", DEFAULT_FOLLOW_INTERVAL, "Interval (in seconds) to sleep during a follow")
	flag.IntVar(&waitPid, "p", DEFAULT_WAIT_PID, "with -f, terminate after process ID, PID dies")
	flag.BoolVar(&verboseOutput, "v", DEFAULT_VERBOSE_OUTPUT, "always output metadata for log events")
	flag.BoolVar(&hideMetadata, "q", DEFAULT_HIDE_METADATA, "never output metadata for log events")
	flag.BoolVar(&listGroups, "l", DEFAULT_LIST_GROUPS, "list available log groups and exit")
	flag.Parse()
	if verboseOutput && !hideMetadata {
		hideMetadata = false
	} else if !verboseOutput && !hideMetadata {
		hideMetadata = true
	}
}

func dumpVersion() {
	fmt.Printf("cloudtail: version %s\n", VERSION)
}

func main() {
	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(1)
	} else if versionDump {
		dumpVersion()
		os.Exit(0)
	}

	session, err := InitSession(verboseOutput, hideMetadata)
	if err != nil {
		os.Exit(1)
	}

	if listGroups {
		session.DumpLogGroups()
		os.Exit(0)
	}

	group, err := session.SearchLogGroups(os.Args[len(os.Args)-1])
	if err != nil {
		os.Exit(1)
	}
	if !followEvents {
		session.DumpLogEvents(&group, numLines)
	} else {
		session.FollowLogEvents(&group, followInterval, waitPid)
	}
	os.Exit(0)
}
