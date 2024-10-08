package gigi

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/go-github/v35/github"
)

func Load() (Config, error) {
	var cfg Config

	// https://docs.github.com/ja/actions/reference/environment-variables#default-environment-variables
	if ci := os.Getenv("CI"); ci != "true" {
		return cfg, fmt.Errorf("not on Actions")
	}
	fmt.Printf("event type: %q\n", os.Getenv("GITHUB_EVENT_NAME"))

	epath := os.Getenv("GITHUB_EVENT_PATH")
	var gpe github.PullRequestEvent
	f, err := os.Open(epath)
	if err != nil {
		return cfg, err
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&gpe); err != nil {
		return cfg, err
	}
	repo := gpe.GetRepo()
	n := strings.Split(repo.GetFullName(), "/")
	if len(n) != 2 {
		return cfg, fmt.Errorf("want \"owrner/repo\", but got %q", repo.GetName())
	}
	cfg.Owner = n[0]
	cfg.Repository = n[1]
	cfg.PullRequestNumber = gpe.GetPullRequest().GetNumber()
	max, err := strconv.Atoi(os.Getenv("GIGI_MAX_ADDED_COUNT"))
	if err != nil {
		return cfg, err
	}
	cfg.MaxAddedCount = max
	cfg.GitHubToken = os.Getenv("GIGI_GITHUB_TOKEN")
	fp := os.Getenv("GIGI_FILTER_PATTERN")
	if len(fp) != 0 {
		cp, err := regexp.Compile(fp)
		if err != nil {
			return cfg, err
		}
		cfg.Filter = cp
	}
	return cfg, nil
}
