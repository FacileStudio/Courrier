package mail

import (
	"context"
	"io"
	"log/slog"
	"sync/atomic"
	"testing"
	"time"

	"github.com/FacileStudio/tronc/testdb"
)

func countingWatcher(polls *atomic.Int32) *Service {
	s := NewService(testORM, nil)
	s.poll = func(context.Context, *slog.Logger) { polls.Add(1) }
	return s
}

func TestOnlyTheLockHolderPolls(t *testing.T) {
	if testORM == nil {
		t.Skip(testdb.SkipReason("createdb courrier_test  # or point at any scratch database"))
	}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	var pollsA, pollsB atomic.Int32
	ctxA, cancelA := context.WithCancel(t.Context())
	defer cancelA()
	go countingWatcher(&pollsA).WatchInboxes(ctxA, 20*time.Millisecond, logger)
	go countingWatcher(&pollsB).WatchInboxes(t.Context(), 20*time.Millisecond, logger)

	waitForPolls(t, func() bool { return pollsA.Load()+pollsB.Load() >= 3 })
	first, second := &pollsA, &pollsB
	if pollsA.Load() == 0 {
		first, second = &pollsB, &pollsA
	}
	time.Sleep(200 * time.Millisecond)
	if got := second.Load(); got != 0 {
		t.Fatalf("the replica without the lock polled %d times", got)
	}
	if first == &pollsA {
		cancelA()
		waitForPolls(t, func() bool { return pollsB.Load() >= 2 })
	}
}

func waitForPolls(t *testing.T, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(20 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("timed out")
}
