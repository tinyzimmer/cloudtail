package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

type logEvent cloudwatchlogs.OutputLogEvent
type logGroup cloudwatchlogs.LogGroup
type logStream cloudwatchlogs.LogStream

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
					LogGroupName: group.LogGroupName,
				}
				logGroups = append(logGroups, tgroup)
			}
			return true
		})
	if err != nil {
		log.Fatal(err)
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
		log.Println("Multiple matching log groups")
		log.Fatal(results)
	} else if len(results) == 0 {
		log.Fatalf("No matching log groups found for: %s", searchGroup)
	} else {
		lgroup = results[0]
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
		log.Fatal(err)
	}
	for _, x := range resp.LogStreams {
		stream := logStream{
			LogStreamName: x.LogStreamName,
		}
		logStreams = append(logStreams, stream)
	}
	return
}

func (s LogSession) CollectEvents(group *logGroup, numEvents int) (events []logEvent) {
	for _, stream := range s.LogStreams {
		if len(events) >= numEvents {
			break
		}
		resp, err := s.LogService.GetLogEvents(&cloudwatchlogs.GetLogEventsInput{
			Limit:         aws.Int64(int64(numEvents)),
			LogGroupName:  group.LogGroupName,
			LogStreamName: stream.LogStreamName,
		})
		if err != nil {
			log.Fatal(err)
		}
		for _, event := range resp.Events {
			if len(events) < numEvents {
				tevent := logEvent{
					Message:   event.Message,
					Timestamp: event.Timestamp,
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
	events = s.CollectEvents(group, numEvents)
	// sort the events by timestamp
	sorted := sortEvents(events)
	// dump the events to stdout
	for _, event := range sorted {
		fmt.Println(formatEvent(event))
	}
}

func (s LogSession) FollowLogEvents(group *logGroup, interval int) {
	var oldEvents []logEvent
	var newEvents []logEvent
	s.RefreshLogStreams(group)
	newEvents = s.CollectEvents(group, DEFAULT_LOG_LINES)
	sorted := sortEvents(newEvents)
	for _, event := range sorted {
		fmt.Println(formatEvent(event))
		oldEvents = append(oldEvents, event)
	}
	for {
		newEvents = s.CollectEvents(group, DEFAULT_LOG_LINES)
		sorted := sortEvents(newEvents)
		for _, event := range sorted {
			if eventIsNew(event, oldEvents) {
				fmt.Println(formatEvent(event))
				oldEvents = append(oldEvents, event)
			}
		}
		time.Sleep(time.Duration(interval) * time.Second)
		go s.RefreshLogStreams(group)
	}
}

func formatEvent(event logEvent) (output string) {
	tm := time.Unix((*event.Timestamp / 1000), 0) // convert to UTC string with offset
	output = fmt.Sprintf("%s: %s", tm, strings.TrimSpace(*event.Message))
	return
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
