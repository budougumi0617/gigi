package gigi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/v35/github"
	"github.com/reviewdog/reviewdog/diff"
)

func GetDiffs(ctx context.Context) error {
	cli := github.NewClient(nil)
	pr, _, err := cli.PullRequests.Get(ctx, "budougumi0617", "nrseg", 14)
	if err != nil {
		return err
	}
	durl := pr.GetDiffURL()
	if durl == "" {
		return fmt.Errorf("cannot get diff source")
	}
	resp, err := http.Get(durl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fds, err := diff.ParseMultiFile(resp.Body)
	for _, fd := range fds {
		fmt.Printf("%+v", fd)
	}

	return nil
}
