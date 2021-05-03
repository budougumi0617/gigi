package gigi

import (
	"context"
	"testing"
)

func TestGetDiffs(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetDiffs(context.TODO()); err != nil {
				t.Errorf("GetDiffs() error = %v", err)
			}
		})
	}
}
