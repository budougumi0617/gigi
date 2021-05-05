package gigi

import (
	"context"
	"fmt"
	"net/http"
	"path"

	"github.com/google/go-github/v35/github"
	"github.com/reviewdog/reviewdog/diff"
	"golang.org/x/oauth2"
)

func GetDiffs(ctx context.Context, cfg Config) (*Result, error) {
	var hc *http.Client
	if len(cfg.GitHubToken) != 0 {
		hc = oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: cfg.GitHubToken},
		))
	}
	cli := github.NewClient(hc)
	pr, _, err := cli.PullRequests.Get(ctx, cfg.Owner, cfg.Repository, cfg.PullRequestNumber)
	if err != nil {
		return nil, err
	}
	durl := pr.GetDiffURL()
	if durl == "" {
		return nil, fmt.Errorf("cannot get diff source")
	}
	resp, err := http.Get(durl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	fds, err := diff.ParseMultiFile(resp.Body)
	if err != nil {
		return nil, err
	}

	r := &Result{}
	for _, fd := range fds {
		f := File{Name: path.Base(fd.PathNew)}
		for _, h := range fd.Hunks {
			for _, line := range h.Lines {
				if line.Type == diff.LineAdded {
					f.AddedCount++
				}
			}
		}
		r.TotalAddedCount += f.AddedCount
		r.Files = append(r.Files, f)
	}

	return r, nil
}
