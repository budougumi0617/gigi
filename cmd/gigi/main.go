package main

import (
	"context"
	"errors"
	"fmt"
	"os"

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
	ctx := context.Background()
	cfg, err := gigi.Load()
	if err != nil {
		if errors.Is(err, gigi.ErrNoEventTypePullRequest) {
			fmt.Println("not support event type")
			return 0
		}
		fmt.Printf("cannot load setting: %v\n", err)
		return 1
	}
	cfg.Version = Version
	cfg.Revision = Revision
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

	if cfg.MaxAddedCount < result.TotalAddedCount {
		if err := gigi.Report(ctx, cfg, result); err != nil {
			fmt.Printf("failed to report: %v\n", err)
		}
		// alert unexpected result.
		return 1
	}
	return 0
}
