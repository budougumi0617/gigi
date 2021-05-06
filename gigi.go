package gigi

import (
	"bytes"
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
		if cfg.Filter != nil && cfg.Filter.MatchString(f.Name) {
			r.Filtered = append(r.Filtered, f)
		} else {
			r.TotalAddedCount += f.AddedCount
			r.Files = append(r.Files, f)
		}
	}

	return r, nil
}

func Report(ctx context.Context, cfg Config, result *Result) error {
	var bb bytes.Buffer
	if _, err := fmt.Fprintf(&bb, `## Pull Request is too large
allow added line is *%d*, but got *%d* lines.

`, cfg.MaxAddedCount, result.TotalAddedCount); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(&bb, `### Added files
|name|added lines|
|---|---|
`); err != nil {
		return err
	}
	for _, file := range result.Files {
		if _, err := fmt.Fprintf(&bb, "| %s | `%d` |\n", file.Name, file.AddedCount); err != nil {
			return err
		}
	}
	if len(result.Filtered) != 0 {
		if _, err := fmt.Fprintf(&bb, `
### Ignored files
|name|added lines|
|---|---|
`); err != nil {
			return err
		}
		for _, file := range result.Filtered {
			if _, err := fmt.Fprintf(&bb, "| %s | `%d` |\n", file.Name, file.AddedCount); err != nil {
				return err
			}
		}
	}
	if _, err := fmt.Fprintf(&bb, "\n\nreported by gigi `%s`(%s)", cfg.Version, cfg.Revision); err != nil {
		return err
	}
	body := bb.String()

	var hc *http.Client
	if len(cfg.GitHubToken) != 0 {
		hc = oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: cfg.GitHubToken},
		))
	}
	cli := github.NewClient(hc)
	// To add a regular comment to a pull request timeline, see "Create an issue comment."
	// https://docs.github.com/en/rest/reference/issues#create-an-issue-comment
	_, _, err := cli.Issues.CreateComment(ctx, cfg.Owner, cfg.Repository, cfg.PullRequestNumber, &github.IssueComment{
		Body: &body,
	})
	if err != nil {
		return err
	}
	return nil
}
