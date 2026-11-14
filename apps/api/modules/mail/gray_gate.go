package mail

import (
	stderrors "errors"
	"net/http"
	"sync"
	"time"

	"github.com/FacileStudio/tronc/bulle"
)

const (
	grayDailyCalls = 40
	grayBackoffMin = 5 * time.Minute
	grayBackoffMax = 2 * time.Hour
)

// triageGate bounds the Bulle calls the inbox triage makes for each owner,
// so triage never spends the daily AI cap the person shares with every app:
// at most grayDailyCalls a UTC day, an exponential backoff after a failure,
// and a pause until the next UTC day after Bulle's 429. It lives in memory.
type triageGate struct {
	mu     sync.Mutex
	now    func() time.Time
	owners map[int64]*ownerTriage
}

type ownerTriage struct {
	day     string
	calls   int
	until   time.Time
	backoff time.Duration
}

// take reserves one call for owner, or refuses while owner backs off or has
// spent today's budget.
func (g *triageGate) take(owner int64) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	now := g.clock()
	state := g.state(owner, now)
	if now.Before(state.until) || state.calls >= grayDailyCalls {
		return false
	}
	state.calls++
	return true
}

// settle records how owner's call ended: a success clears the backoff, a 429
// pauses until the next UTC day, any other failure doubles the backoff from
// grayBackoffMin up to grayBackoffMax.
func (g *triageGate) settle(owner int64, err error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	now := g.clock()
	state := g.state(owner, now)
	var refused *bulle.StatusError
	switch {
	case err == nil:
		state.backoff = 0
	case stderrors.As(err, &refused) && refused.Status == http.StatusTooManyRequests:
		state.until = nextUTCDay(now)
	default:
		state.backoff = min(max(2*state.backoff, grayBackoffMin), grayBackoffMax)
		state.until = now.Add(state.backoff)
	}
}

func (g *triageGate) clock() time.Time {
	if g.now != nil {
		return g.now()
	}
	return time.Now()
}

func (g *triageGate) state(owner int64, now time.Time) *ownerTriage {
	if g.owners == nil {
		g.owners = map[int64]*ownerTriage{}
	}
	state, ok := g.owners[owner]
	if !ok {
		state = &ownerTriage{}
		g.owners[owner] = state
	}
	if day := now.UTC().Format(time.DateOnly); state.day != day {
		state.day, state.calls = day, 0
	}
	return state
}

func nextUTCDay(now time.Time) time.Time {
	return now.UTC().Truncate(24 * time.Hour).Add(24 * time.Hour)
}
