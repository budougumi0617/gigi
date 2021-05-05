package gigi

import (
	"context"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetDiffs(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		want Result
	}{
		{
			name: "test",
			// https://patch-diff.githubusercontent.com/raw/budougumi0617/nrseg/pull/14.diff
			cfg: Config{
				Owner:             "budougumi0617",
				Repository:        "nrseg",
				PullRequestNumber: 14,
			},
			want: Result{
				TotalAddedCount: 173,
				Files: []File{
					{Name: "inspect.go", AddedCount: 51},
					{Name: "nrseg.go", AddedCount: 104},
					{Name: "process.go", AddedCount: 18},
				},
			},
		},
		{
			name: "filter",
			// https://patch-diff.githubusercontent.com/raw/budougumi0617/nrseg/pull/16.diff
			cfg: Config{
				Owner:             "budougumi0617",
				Repository:        "nrseg",
				PullRequestNumber: 16,
				Filter:            regexp.MustCompile("go.sum|.*_test.go"),
			},
			want: Result{
				TotalAddedCount: 16,
				Files: []File{
					{Name: "inspect.go", AddedCount: 0},
					{Name: "nrseg.go", AddedCount: 16},
				},
				Filtered: []File{
					{Name: "go.sum", AddedCount: 0},
					{Name: "nrseg_test.go", AddedCount: 72},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := GetDiffs(context.TODO(), tt.cfg); err != nil {
				t.Errorf("GetDiffs() error = %v", err)
			} else if diff := cmp.Diff(*got, tt.want); len(diff) != 0 {
				t.Errorf("-got +want %v", diff)
			}
		})
	}
}
