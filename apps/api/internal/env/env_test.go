package env

import (
	"slices"
	"testing"
)

func TestWithAppOriginsAddsEveryAppOnce(t *testing.T) {
	got := withAppOrigins([]string{"https://courrier.example", "https://localhost"}, DefaultMobileRedirect)
	want := []string{"https://courrier.example", "https://localhost", "tauri://localhost", "http://tauri.localhost"}
	if !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if got := withAppOrigins([]string{"https://courrier.example"}, ""); len(got) != 1 {
		t.Fatalf("app login off must leave the list alone, got %v", got)
	}
}
