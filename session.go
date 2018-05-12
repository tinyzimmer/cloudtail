package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

type LogSession struct {
	LogService   *cloudwatchlogs.CloudWatchLogs
	LogGroups    []logGroup
	LogStreams   []logStream
	Verbose      bool
	HideMetadata bool
}

func InitSession(verbose bool, hideMetadata bool) (logSession LogSession) {
	// Get credentials and a session
	logSession.Verbose = verbose
	logSession.HideMetadata = hideMetadata
	if verbose {
		LogInfo("Retrieving and testing AWS Credentials")
	}
	creds := getCreds()
	if verbose {
		LogInfo("Validated AWS Credentials")
	}
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String("us-west-2"),
	}))
	logSession.LogService = cloudwatchlogs.New(sess)
	if verbose {
		LogInfo("Created CloudWatch Session. Gathering Log Groups.")
	}
	// initialize log group list
	logSession.RefreshLogGroups()
	return logSession
}

func getCreds() (creds *credentials.Credentials) {
	// Check for credentials in following order
	//    1. Environment Variables
	//    2. EC2 IAM Role
	//    3. Shared Credentials File
	sess := session.Must(session.NewSession())
	creds = credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.EnvProvider{},
			&ec2rolecreds.EC2RoleProvider{
				Client: ec2metadata.New(sess),
			},
			&credentials.SharedCredentialsProvider{},
		})
	_, err := creds.Get()
	if err != nil {
		LogFatal(err)
	}
	return
}
