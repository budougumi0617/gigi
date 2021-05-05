package gigi

import (
	"fmt"
	"log"
	"os"
)

type EventType string

const EventTypePullRequest EventType = "pull_request"

func load() error {
	if ci := os.Getenv("CI"); ci != "true" {
		log.Println("")
		return fmt.Errorf("format is ")
	}
	return nil
}
