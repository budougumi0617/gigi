package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/budougumi0617/gigi"
)

var (
	Version  = "tip"
	Revision = "none"
)

func main() {
	os.Exit(run())
}
func run() int {
	flag.Parse()
	ctx := context.Background()
	cfg := gigi.Config{
		Owner:             "",
		Repository:        "",
		PullRequestNumber: 0,
		GitHubToken:       "",
		MaxAddedCount:     0,
		Filter:            nil,
		Version:           "",
		Revision:          "",
	}
	if len(flag.Args()) != 3 {
		fmt.Printf("usage: gigicli OWNER_NAME REPO_NAME PR_NUMBER\n")
		return 1
	}
	cfg.Version = Version
	cfg.Revision = Revision
	cfg.GitHubToken = os.Getenv("GITHUB_TOKEN")
	cfg.Owner = flag.Arg(0)
	cfg.Repository = flag.Arg(1)
	pr, err := strconv.Atoi(flag.Arg(2))
	if err != nil {
		fmt.Printf("cannot convert pr: %v\n", err)
		return 1
	}
	cfg.MaxAddedCount = 300 // 暫定
	cfg.PullRequestNumber = pr
	fmt.Printf("%s/%s pr number %d check!\n", cfg.Owner, cfg.Repository, cfg.PullRequestNumber)
	result, err := gigi.GetDiffs(ctx, cfg)
	if err != nil {
		fmt.Printf("failed to get result: %v\n", err)
		return 1
	}

	fmt.Printf("result added count %d\n", result.TotalAddedCount)
	if len(result.Files) != 0 {
		fmt.Printf("-----found files\n")
		for _, file := range result.Files {
			fmt.Printf("%q:%d\n", file.Name, file.AddedCount)
		}
	}
	if len(result.Filtered) != 0 {
		fmt.Printf("-----ignore files\n")
		for _, file := range result.Filtered {
			fmt.Printf("%q:%d\n", file.Name, file.AddedCount)
		}
	}

	return 0
}
