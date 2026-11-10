package middleware

import (
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/httpjson"
)

type visitor struct {
	mu        sync.Mutex
	count     int
	windowEnd time.Time
}

// RateLimit returns middleware that allows at most requests per window per
// remote address, replying 429 with a Retry-After header once the budget is
// spent. It is middleware in the chi sense, not an http.Handler.
func RateLimit(requests int, window time.Duration) func(http.Handler) http.Handler {
	return RateLimitBy(requests, window, remoteIP)
}

// RateLimitBy is RateLimit with the budget keyed by key(r) instead of the
// remote address, such as the person an access token speaks for. An empty
// key falls back to the remote address.
func RateLimitBy(requests int, window time.Duration, key func(*http.Request) string) func(http.Handler) http.Handler {
	var visitors sync.Map
	go sweep(&visitors, window)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := key(r)
			if id == "" {
				id = remoteIP(r)
			}
			now := time.Now()
			val, _ := visitors.LoadOrStore(id, &visitor{windowEnd: now.Add(window)})
			if retryAfter, ok := val.(*visitor).take(now, requests, window); !ok {
				w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
				httpjson.WriteError(w, errors.RateLimited("too many requests, try again later"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// take spends one request from the visitor's window, or reports how many
// seconds remain until the window resets.
func (v *visitor) take(now time.Time, requests int, window time.Duration) (int, bool) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if now.After(v.windowEnd) {
		v.count = 0
		v.windowEnd = now.Add(window)
	}
	if v.count >= requests {
		return int(time.Until(v.windowEnd).Seconds()) + 1, false
	}
	v.count++
	return 0, true
}

// sweep drops visitors whose window has ended, once per window.
func sweep(visitors *sync.Map, window time.Duration) {
	ticker := time.NewTicker(window)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		visitors.Range(func(key, value any) bool {
			v := value.(*visitor)
			v.mu.Lock()
			expired := now.After(v.windowEnd)
			v.mu.Unlock()
			if expired {
				visitors.Delete(key)
			}
			return true
		})
	}
}

// remoteIP is the address part of RemoteAddr, which RealIP has already set
// to the visitor behind trusted proxies.
func remoteIP(r *http.Request) string {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	if ip == "" {
		return r.RemoteAddr
	}
	return ip
}
