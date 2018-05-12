package main

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
)

func TestLogError(t *testing.T) {
	LogError("test-error")
}

func TestLogInfo(t *testing.T) {
	LogInfo("test-info")
}

func TestLogWarn(t *testing.T) {
	LogWarn("test-warning")
}

func TestLogDebug(t *testing.T) {
	LogDebug("test-debug")
}

func TestOutputLogStream(t *testing.T) {
	testLogStream := logStream{
		Arn:                 aws.String("fake-arn"),
		CreationTime:        aws.Int64(0),
		FirstEventTimestamp: aws.Int64(0),
		LastEventTimestamp:  aws.Int64(10),
		LastIngestionTime:   aws.Int64(15),
		LogStreamName:       aws.String("test-name"),
		StoredBytes:         aws.Int64(100),
		UploadSequenceToken: aws.String("test-token"),
	}
	logLogStream(testLogStream)
}

func TestOutputLogGroup(t *testing.T) {
	testLogGroup := logGroup{
		Arn:               aws.String("test-arn"),
		CreationTime:      aws.Int64(0),
		LogGroupName:      aws.String("test-name"),
		MetricFilterCount: aws.Int64(10),
		StoredBytes:       aws.Int64(100),
	}
	logLogGroup(testLogGroup)
}

func TestLogLogEvent(t *testing.T) {
	testEvent := logEvent{
		IngestionTime: aws.Int64(0),
		Timestamp:     aws.Int64(0),
		Message:       aws.String("test-message"),
	}
	LogEvent(testEvent, false, false)
	LogEvent(testEvent, true, false)
	LogEvent(testEvent, false, true)
}

func TestLogFatal(t *testing.T) {
	LogFatal(errors.New("test-fatal"), 0)
}
