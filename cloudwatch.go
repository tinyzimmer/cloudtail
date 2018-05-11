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

func (s LogSession) GetLogGroups() (logGroups []*cloudwatchlogs.LogGroup) {
	// Retrieve a list of all known log groups
	err := s.LogService.DescribeLogGroupsPages(&cloudwatchlogs.DescribeLogGroupsInput{},
		func(page *cloudwatchlogs.DescribeLogGroupsOutput, lastPage bool) bool {
			for _, group := range page.LogGroups {
				logGroups = append(logGroups, group)
			}
			return true
		})
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (s LogSession) SearchLogGroups(searchGroup string) (logGroup *cloudwatchlogs.LogGroup) {
	// Look for log groups that match user input
	results := make([]*cloudwatchlogs.LogGroup, 0)
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
		logGroup = results[0]
	}
	return
}

func (s LogSession) GetLogStreams(logGroup *cloudwatchlogs.LogGroup) (logStreams []*cloudwatchlogs.LogStream) {
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
	logStreams = resp.LogStreams
	return
}

func (s LogSession) DumpLogEvents(group *cloudwatchlogs.LogGroup, streams []*cloudwatchlogs.LogStream, numEvents int) {
	var events []*cloudwatchlogs.OutputLogEvent
	// iterate the streams and create a slice of events
	for _, stream := range streams {
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
				events = append(events, event)
			} else {
				break
			}
		}
	}
	// sort the events by timestamp
	sort.Slice(events, func(i, j int) bool {
		return *events[i].Timestamp < *events[j].Timestamp
	})
	// dump the events to stdout
	for _, event := range events {
		tm := time.Unix((*event.Timestamp / 1000), 0)
		fmt.Println(fmt.Sprintf("%s: %s", tm, strings.TrimSpace(*event.Message)))
	}
}
