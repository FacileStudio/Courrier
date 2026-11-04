package assist

import (
	"context"
	"strings"
	"testing"
)

func TestValidateAndPrompt(t *testing.T) {
	cases := map[string]Request{
		"unknown action":       {Action: "delete", Text: "x"},
		"empty text":           {Action: "shorten"},
		"translate, no target": {Action: "translate", Text: "x"},
		"reply, no thread":     {Action: "reply"},
	}
	for name, req := range cases {
		if err := validate(&req); err == nil {
			t.Errorf("%s: want an error", name)
		}
	}
	req := Request{Action: "translate", Text: "Bonjour", Language: "English", Instructions: "keep it short"}
	if err := validate(&req); err != nil {
		t.Fatal(err)
	}
	system, user := buildPrompt(req, "You are Pixel.")
	if !strings.HasPrefix(system, "You are Pixel.") || !strings.Contains(user, "into English") || !strings.Contains(user, "<text>\nBonjour") {
		t.Fatalf("prompt:\n%s\n%s", system, user)
	}
}

func TestOffServiceSaysSo(t *testing.T) {
	s := NewService(nil, "")
	if st := s.Status(context.Background(), 1); st.Available || st.Reason != "off" {
		t.Fatalf("status = %+v", st)
	}
	if _, err := s.Run(context.Background(), 1, Request{Action: "fix", Text: "x"}); err == nil || !strings.Contains(err.Error(), "off") {
		t.Fatalf("run = %v", err)
	}
}
