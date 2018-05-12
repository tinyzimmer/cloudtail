package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type OutputLogEventLiteral struct {
	IngestionTime string
	Message       string
	Timestamp     string
}

type OutputLogGroupLiteral struct {
	Arn               string
	CreationTime      string
	LogGroupName      string
	MetricFilterCount string
	StoredBytes       string
}

type OutputLogStreamLiteral struct {
	Arn                 string
	CreationTime        string
	FirstEventTimestamp string
	LastEventTimestamp  string
	LastIngestionTime   string
	LogStreamName       string
	StoredBytes         string
	UploadSequenceToken string
}

func LogError(msg string) {
	line := fmt.Sprintf("%s %s", ColorRed("cloudtail-error"), msg)
	log.Println(line)
}

func LogInfo(msg string) {
	line := fmt.Sprintf("%s %s", ColorGreen("cloudtail-info"), msg)
	log.Println(line)
}

func LogWarn(msg string) {
	line := fmt.Sprintf("%s %s", ColorYellow("cloudtail-warning"), msg)
	log.Println(line)
}

func LogDebug(msg string) {
	line := fmt.Sprintf("%s %s", ColorBlue("cloudtail-debug"), msg)
	log.Println(line)
}

func logLogStream(stream logStream) {
	output := OutputLogStreamLiteral{
		Arn:                 fmt.Sprintf("%+v\n", *stream.Arn),
		CreationTime:        fmt.Sprintf("%+v\n", *stream.CreationTime),
		FirstEventTimestamp: fmt.Sprintf("%+v\n", *stream.FirstEventTimestamp),
		LastEventTimestamp:  fmt.Sprintf("%+v\n", *stream.LastEventTimestamp),
		LastIngestionTime:   fmt.Sprintf("%+v\n", *stream.LastIngestionTime),
		LogStreamName:       fmt.Sprintf("%+v\n", *stream.LogStreamName),
		StoredBytes:         fmt.Sprintf("%+v\n", *stream.StoredBytes),
		UploadSequenceToken: fmt.Sprintf("%+v\n", *stream.UploadSequenceToken),
	}
	fmt.Println(fmt.Sprintf("%s\n%+v", ColorWhite("cloudtail-log-stream"), output))
}

func logLogGroup(group logGroup) {
	output := OutputLogGroupLiteral{
		Arn:               fmt.Sprintf("%+v\n", *group.Arn),
		CreationTime:      fmt.Sprintf("%+v\n", *group.CreationTime),
		LogGroupName:      fmt.Sprintf("%+v\n", *group.LogGroupName),
		MetricFilterCount: fmt.Sprintf("%+v\n", *group.MetricFilterCount),
		StoredBytes:       fmt.Sprintf("%+v\n", *group.StoredBytes),
	}
	fmt.Println(fmt.Sprintf("%s\n%+v", ColorCyan("cloudtail-log-group"), output))
}

func logLogEvent(event logEvent) (output OutputLogEventLiteral) {
	output = OutputLogEventLiteral{
		IngestionTime: fmt.Sprintf("%+v\n", *event.IngestionTime),
		Message:       *event.Message,
		Timestamp:     fmt.Sprintf("%+v\n", convertTimestamp(*event.Timestamp)),
	}
	return
}

func LogEvent(event logEvent, verbose bool, hideMetadata bool) {
	var text interface{}
	var line string
	if verbose && !hideMetadata {
		text = logLogEvent(event)
		line = fmt.Sprintf("%s\n%+v", ColorPurple("cloudtail-log-event"), text)
	} else {
		tm := convertTimestamp(*event.Timestamp)
		text = fmt.Sprintf("%s: %s", tm, strings.TrimSpace(*event.Message))
		line = fmt.Sprintf("%s %+v", ColorPurple("cloudtail-log-event"), text)
	}
	fmt.Println(line)
}

func LogFatal(err error) {
	LogError(err.Error())
}

func ColorRed(value string) string {
	return fmt.Sprintf("\033[0;31m%s\033[0m", value)
}

func ColorGreen(value string) string {
	return fmt.Sprintf("\033[0;32m%s\033[0m", value)
}

func ColorYellow(value string) string {
	return fmt.Sprintf("\033[0;33m%s\033[0m", value)
}

func ColorBlue(value string) string {
	return fmt.Sprintf("\033[0;34m%s\033[0m", value)
}

func ColorPurple(value string) string {
	return fmt.Sprintf("\033[0;35m%s\033[0m", value)
}

func ColorCyan(value string) string {
	return fmt.Sprintf("\033[0;36m%s\033[0m", value)
}

func ColorWhite(value string) string {
	return fmt.Sprintf("\033[0;37m%s\033[0m", value)
}

func convertTimestamp(timestamp int64) (tm time.Time) {
	tm = time.Unix((timestamp / 1000), 0) // convert to UTC string with offset
	return
}
