# Solution Steps

1. Create a scheduler package with a Scheduler struct that tracks all registered jobs and their goroutines, using a mutex for concurrency management.

2. Define a Job interface (with Run(ctx) and Name()) and implement registration (AddJob), removal (RemoveJob), and clean goroutine/ticker management inside the scheduler.

3. Use a context in the scheduler to coordinate clean shutdown (via Shutdown()), signaling goroutines to gracefully terminate either on context cancellation or explicit job removal.

4. Ensure each registered job runs in its own goroutine, scheduling actual task execution independently and safely using time.Ticker and select to multiplex between ticker, stop channel, and context.

5. Instrument the scheduler and jobs to use sync.WaitGroup for tracking in-flight job executions, so graceful shutdown waits until all jobs are complete.

6. Implement three example jobs in a 'jobs' package (ExportJob, LeaderboardJob, NotifyJob), each with different schedules (intervals) and logging to demonstrate execution.

7. In main.go, instantiate the scheduler, add the three jobs, and demonstrate dynamic removal and addition of another job at runtime. Handle SIGINT/SIGTERM for graceful system-wide shutdown.

8. Add an additional example job (NotifyJob2) to demonstrate extensibility and hot registration of new job types at runtime.

9. Write a scheduler_test.go file with tests that prove jobs run at correct intervals, can be dynamically registered/unregistered, and that shutdown waits for in-flight jobs and cancels future ones cleanly.

10. Ensure all parts are goroutine-leak safe: tickers are stopped and goroutines exited on removal, replacement, or scheduler shutdown.

