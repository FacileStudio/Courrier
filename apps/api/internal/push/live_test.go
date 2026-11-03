package push

import (
	"context"
	"slices"
	"testing"
)

func TestHubPublishesToOpenStreamsOnly(t *testing.T) {
	hub := NewHub(HubConfig{})
	notes, stop := hub.Subscribe(7)
	if got := hub.LiveUsers(); !slices.Equal(got, []int64{7}) {
		t.Fatalf("live users %v", got)
	}
	hub.Notify(context.Background(), 7, Notification{Title: "New mail from Ada"})
	hub.Notify(context.Background(), 8, Notification{Title: "not for 7"})
	if n := <-notes; n.Title != "New mail from Ada" {
		t.Fatalf("got %+v", n)
	}
	select {
	case n := <-notes:
		t.Fatalf("another user's notification leaked: %+v", n)
	default:
	}
	stop()
	if got := hub.LiveUsers(); len(got) != 0 {
		t.Fatalf("closed stream still listed: %v", got)
	}
	var nilHub *Hub
	nilHub.Notify(context.Background(), 7, Notification{})
	if nilHub.LiveUsers() != nil {
		t.Fatal("a nil hub has no live users")
	}
}
