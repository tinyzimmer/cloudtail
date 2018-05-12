package main

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

type logEvent cloudwatchlogs.OutputLogEvent
type logGroup cloudwatchlogs.LogGroup
type logStream cloudwatchlogs.LogStream

var oldLogStreams []logStream

func (s *LogSession) RefreshLogGroups() {
	s.LogGroups = make([]logGroup, 0)
	groups := s.GetLogGroups()
	for _, group := range groups {
		s.LogGroups = append(s.LogGroups, group)
	}
}

func (s *LogSession) RefreshLogStreams(group *logGroup) {
	s.LogStreams = make([]logStream, 0)
	streams := s.GetLogStreams(group)
	for _, stream := range streams {
		s.LogStreams = append(s.LogStreams, stream)
	}
}

func (s LogSession) GetLogGroups() (logGroups []logGroup) {
	// Retrieve a list of all known log groups
	err := s.LogService.DescribeLogGroupsPages(&cloudwatchlogs.DescribeLogGroupsInput{},
		func(page *cloudwatchlogs.DescribeLogGroupsOutput, lastPage bool) bool {
			for _, group := range page.LogGroups {
				tgroup := logGroup{
					Arn:               group.Arn,
					CreationTime:      group.CreationTime,
					KmsKeyId:          group.KmsKeyId,
					LogGroupName:      group.LogGroupName,
					MetricFilterCount: group.MetricFilterCount,
					RetentionInDays:   group.RetentionInDays,
					StoredBytes:       group.StoredBytes,
				}
				logGroups = append(logGroups, tgroup)
			}
			return true
		})
	if err != nil {
		LogFatal(err, 1)
	}
	return
}

func (s LogSession) SearchLogGroups(searchGroup string) (lgroup logGroup) {
	// Look for log groups that match user input
	results := make([]logGroup, 0)
	for _, group := range s.LogGroups {
		if strings.Contains(*group.LogGroupName, searchGroup) {
			results = append(results, group)
		}
	}
	if len(results) > 1 {
		err := errors.New("Multiple matching log groups. Try narrowing down the search.")
		LogFatal(err, 1)
	} else if len(results) == 0 {
		err := errors.New(fmt.Sprintf("No matching log groups found for: %s", searchGroup))
		LogFatal(err, 1)
	} else {
		lgroup = results[0]
		if s.Verbose && !s.HideMetadata {
			logLogGroup(lgroup)
		}
	}
	return
}

func (s LogSession) GetLogStreams(logGroup *logGroup) (logStreams []logStream) {
	// Retrieve log streams associated with a log group
	// Sort by descending timestamp and only give us the last 10 streams
	resp, err := s.LogService.DescribeLogStreams(&cloudwatchlogs.DescribeLogStreamsInput{
		Descending:   aws.Bool(true),
		Limit:        aws.Int64(10),
		LogGroupName: logGroup.LogGroupName,
		OrderBy:      aws.String("LastEventTime"),
	})
	if err != nil {
		LogFatal(err, 1)
	}
	for _, x := range resp.LogStreams {
		stream := logStream{
			Arn:                 x.Arn,
			CreationTime:        x.CreationTime,
			FirstEventTimestamp: x.FirstEventTimestamp,
			LastEventTimestamp:  x.LastEventTimestamp,
			LastIngestionTime:   x.LastIngestionTime,
			LogStreamName:       x.LogStreamName,
			StoredBytes:         x.StoredBytes,
			UploadSequenceToken: x.UploadSequenceToken,
		}
		logStreams = append(logStreams, stream)
	}
	return
}

func (s LogSession) CollectEvents(group *logGroup, numEvents int, waitPid int) (events []logEvent) {
	for _, stream := range s.LogStreams {
		checkPid(waitPid)
		if s.Verbose && !s.HideMetadata {
			if streamIsNew(stream) {
				logLogStream(stream)
			}
			oldLogStreams = append(oldLogStreams, stream)
		}
		if len(events) >= numEvents {
			break
		}
		resp, err := s.LogService.GetLogEvents(&cloudwatchlogs.GetLogEventsInput{
			Limit:         aws.Int64(int64(numEvents)),
			LogGroupName:  group.LogGroupName,
			LogStreamName: stream.LogStreamName,
		})
		if err != nil {
			LogFatal(err, 1)
		}
		for _, event := range resp.Events {
			if len(events) < numEvents {
				tevent := logEvent{
					Message:       event.Message,
					Timestamp:     event.Timestamp,
					IngestionTime: event.IngestionTime,
				}
				events = append(events, tevent)
			} else {
				break
			}
		}
	}
	return
}

func (s LogSession) DumpLogEvents(group *logGroup, numEvents int) {
	var events []logEvent
	// iterate the streams and create a slice of events
	s.RefreshLogStreams(group)
	events = s.CollectEvents(group, numEvents, -1)
	// sort the events by timestamp
	sorted := sortEvents(events)
	// dump the events to stdout
	for _, event := range sorted {
		LogEvent(event, s.Verbose, s.HideMetadata)
	}
}

func (s LogSession) FollowLogEvents(group *logGroup, interval int, waitPid int) {
	checkPid(waitPid)
	var oldEvents []logEvent
	var newEvents []logEvent
	s.RefreshLogStreams(group)
	newEvents = s.CollectEvents(group, DEFAULT_LOG_LINES, waitPid)
	sorted := sortEvents(newEvents)
	for _, event := range sorted {
		LogEvent(event, s.Verbose, s.HideMetadata)
		oldEvents = append(oldEvents, event)
	}
	for {
		checkPid(waitPid)
		newEvents = s.CollectEvents(group, DEFAULT_LOG_LINES, waitPid)
		sorted := sortEvents(newEvents)
		for _, event := range sorted {
			if eventIsNew(event, oldEvents) {
				LogEvent(event, s.Verbose, s.HideMetadata)
				oldEvents = append(oldEvents, event)
			}
		}
		time.Sleep(time.Duration(interval) * time.Second)
		go s.RefreshLogStreams(group)
	}
}

func sortEvents(events []logEvent) (sorted []logEvent) {
	sorted = events
	sort.Slice(sorted, func(i, j int) bool {
		return *events[i].Timestamp < *events[j].Timestamp
	})
	return
}

func eventIsNew(newEvent logEvent, events []logEvent) bool {
	for _, event := range events {
		if *newEvent.Timestamp == *event.Timestamp && *newEvent.Message == *event.Message {
			return false
		}
	}
	return true
}

func streamIsNew(newStream logStream) bool {
	for _, stream := range oldLogStreams {
		if *newStream.LogStreamName == *stream.LogStreamName && *newStream.LastEventTimestamp == *stream.LastEventTimestamp {
			return false
		}
	}
	return true
}

func pidRunning(pid int) bool {
	process, err := os.FindProcess(int(pid))
	if err != nil {
		LogWarn(fmt.Sprintf("Process %v exited\n", pid))
		return false
	} else {
		err := process.Signal(syscall.Signal(0))
		if err.Error() == "no such process" {
			LogWarn(fmt.Sprintf("Process %v exited\n", pid))
			return false
		}
	}
	return true
}

func checkPid(waitPid int) {
	if waitPid != -1 {
		if !pidRunning(waitPid) {
			os.Exit(0)
		}
	}
}
