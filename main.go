package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-async-dynamic-scheduler/jobs"
	"go-async-dynamic-scheduler/scheduler"
)

func main() {
	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	fmt.Println("Starting Dynamic Asynchronous Scheduler...")
	sched := scheduler.NewScheduler(rootCtx)

	// Add three jobs with different intervals
	sched.AddJob(jobs.ExportJob{}, 2*time.Second)
	sched.AddJob(jobs.LeaderboardJob{}, 5*time.Second)
	sched.AddJob(jobs.NotifyJob{}, 1*time.Second)

	fmt.Println("Jobs registered. Press Ctrl+C to exit...")

	// Demo: after 7 sec, remove one job and add a different one
	go func() {
		time.Sleep(7 * time.Second)
		fmt.Println("\n[Demo] Removing 'notify-job' and adding a fast 'notify-job-2'")
		sched.RemoveJob("notify-job")
		// New variant of notify job, shows extensibility
		sched.AddJob(jobs.NotifyJob2{}, 500*time.Millisecond)
	}()

	<-rootCtx.Done() // Wait for termination

	fmt.Println("\nShutting down scheduler...")
	sched.Shutdown()
	fmt.Println("All jobs complete. Exiting!")
}
