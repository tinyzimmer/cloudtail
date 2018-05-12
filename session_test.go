package main

import (
	"os"
	"testing"
)

func TestVerboseSession(t *testing.T) {
	_, err := InitSession(true, false)
	if err != nil {
		t.Error(err)
	}
}

func TestBadCredentials(t *testing.T) {
	backupKey := os.Getenv("AWS_ACCESS_KEY_ID")
	backupSecret := os.Getenv("AWS_SECRET_ACCESS_KEY")
	profileBackup := os.Getenv("AWS_PROFILE")
	os.Setenv("AWS_ACCESS_KEY_ID", "")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "")
	os.Setenv("AWS_PROFILE", "NON_EXIST")
	_, err := InitSession(true, false)
	if err == nil {
		t.Error("No error on bad credentials")
	}
	os.Setenv("AWS_ACCESS_KEY_ID", backupKey)
	os.Setenv("AWS_SECRET_ACCESS_KEY", backupSecret)
	os.Setenv("AWS_PROFILE", profileBackup)
}
