package gigi

import "regexp"

type Config struct {
	Owner, Repository string
	PullRequestNumber int
	GitHubToken       string
	MaxAddedCount     int
	Filter            *regexp.Regexp
	Version           string
	Revision          string
}
