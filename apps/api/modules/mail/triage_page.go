package mail

import (
	"net/url"
	"strconv"
)

const (
	maxPage      = 10000
	maxPageLimit = 100
	defaultLimit = 50
)

func parsePaging(v url.Values) (int, int) {
	page, limit := 1, defaultLimit
	if n, err := strconv.Atoi(v.Get("page")); err == nil && n > 1 {
		page = min(n, maxPage)
	}
	if n, err := strconv.Atoi(v.Get("limit")); err == nil && n >= 1 {
		limit = min(n, maxPageLimit)
	}
	return page, limit
}

func pageOf(items []TriageItem, page, limit int) []TriageItem {
	page = min(max(page, 1), maxPage)
	limit = min(max(limit, 1), maxPageLimit)
	start := (page - 1) * limit
	if start >= len(items) {
		return []TriageItem{}
	}
	return items[start:min(start+limit, len(items))]
}
