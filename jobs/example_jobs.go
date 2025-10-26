// Package jobs contains example job implementations for the scheduler.
package jobs

import (
	"context"
	"fmt"
	"time"
)

type ExportJob struct{}
func (ExportJob) Run(ctx context.Context) {
	fmt.Printf("[ExportJob] Exporting data at %v\n", time.Now().Format(time.RFC3339Nano))
	// Simulate doing work (should listen on ctx)
	select {
	case <-time.After(500 * time.Millisecond):
	case <-ctx.Done():
	}
}
func (ExportJob) Name() string { return "export-job" }

type LeaderboardJob struct{}
func (LeaderboardJob) Run(ctx context.Context) {
	fmt.Printf("[LeaderboardJob] Calculating leaderboard at %v\n", time.Now().Format(time.RFC3339Nano))
	select {
	case <-time.After(300 * time.Millisecond):
	case <-ctx.Done():
	}
}
func (LeaderboardJob) Name() string { return "leaderboard-job" }

type NotifyJob struct{}
func (NotifyJob) Run(ctx context.Context) {
	fmt.Printf("[NotifyJob] Sending notifications at %v\n", time.Now().Format(time.RFC3339Nano))
	select {
	case <-time.After(150 * time.Millisecond):
	case <-ctx.Done():
	}
}
func (NotifyJob) Name() string { return "notify-job" }
