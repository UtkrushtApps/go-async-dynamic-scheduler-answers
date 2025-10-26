package scheduler

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

type testJob struct {
	name string
	calls *int32
}
func (t testJob) Run(ctx context.Context) {
	atomic.AddInt32(t.calls, 1)
	select {
	case <-time.After(20 * time.Millisecond):
	case <-ctx.Done():
	}
}
func (t testJob) Name() string { return t.name }

func TestScheduler_AddRemoveAndShutdown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	calls := int32(0)
	job := testJob{name: "test-job", calls: &calls}

	sched := NewScheduler(ctx)
	sched.AddJob(job, 40*time.Millisecond)

	time.Sleep(115 * time.Millisecond) // should run 2-3 times
	sched.RemoveJob("test-job")
	oldCalls := atomic.LoadInt32(&calls)

	if oldCalls < 2 {
		t.Fatalf("expected job to run at least 2 times, got %d", oldCalls)
	}
	// Wait more to ensure the job doesn't run further.
	time.Sleep(90 * time.Millisecond)
	if atomic.LoadInt32(&calls) != oldCalls {
		t.Fatal("job was not properly removed (still running after removal)")
	}

	// Test shutdown with in-flight jobs
	calls2 := int32(0)
	job2 := testJob{name: "shutdown-job", calls: &calls2}
	sched.AddJob(job2, 20*time.Millisecond)
	time.Sleep(45 * time.Millisecond)
	// Should be running 2+ times, then shutdown

	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()
	
	sched.Shutdown()
	final := atomic.LoadInt32(&calls2)
	if final == 0 {
		t.Fatal("job2 did not run before shutdown")
	}

	// After Shutdown, job2 should not run more.
	time.Sleep(30 * time.Millisecond)
	if atomic.LoadInt32(&calls2) != final {
		t.Fatal("shutdown did not stop in-flight jobs properly")
	}
}
