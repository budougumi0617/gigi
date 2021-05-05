package gigi

import (
	"context"
	"fmt"
	"os"

	"github.com/budougumi0617/gigi"
)

func main() {
	os.Exit(run())
}
func run() int {
	ctx := context.Background()
	cfg, err := gigi.Load()
	if err != nil {
		fmt.Printf("cannot load setting: %v", err)
		return 1
	}
	result, err := gigi.GetDiffs(ctx, cfg)
	if err != nil {
		fmt.Printf("failed to get result: %v", err)
		return 1
	}
	if cfg.MaxAddedCount < result.TotalAddedCount {
		if err := gigi.Report(ctx, cfg, result); err != nil {
			fmt.Printf("failed to report: %v", err)
		}
		// alert unexpected result.
		return 1
	}
	return 0
}
