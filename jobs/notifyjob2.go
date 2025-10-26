// Additional job type for live extensibility demo
package jobs

import (
	"context"
	"fmt"
	"time"
)

type NotifyJob2 struct{}
func (NotifyJob2) Run(ctx context.Context) {
	fmt.Printf("[NotifyJob2] 🚀 Fast notification at %v\n", time.Now().Format(time.RFC3339Nano))
	select {
	case <-time.After(50 * time.Millisecond):
	case <-ctx.Done():
	}
}
func (NotifyJob2) Name() string { return "notify-job-2" }
