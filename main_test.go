package main

import (
	"os"
	"strings"
	"testing"
)

const appName = "jrnl"

func asArgs(command string) []string {
	return strings.Split(command, " ")
}

func TestMain(m *testing.M) {
	os.Setenv("DEV", "true")
	code := m.Run()
	os.Exit(code)
}

func TestShowPaths(t *testing.T) {
	command := appName + " showpaths"

	err := run(asArgs(command))

	if err != nil {
		t.Error("Error running showdbpath command: ", err)
	}
}
