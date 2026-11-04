package assist

import (
	"fmt"
	"strings"
)

const baseSystem = "You help the user write email inside their mail client. " +
	"Answer with the resulting text only: no preamble, no quotes around it, no explanation, no subject line. " +
	"Keep the user's language unless asked to translate, keep names, dates, numbers and links exactly, and never invent facts."

var actionPrompts = map[string]string{
	"rewrite":   "Rewrite the text so it reads clearly and naturally. Keep its meaning, tone and length.",
	"shorten":   "Make the text noticeably shorter. Keep every fact and request; drop filler.",
	"expand":    "Develop the text into a fuller version: add the transitions, politeness and detail a careful sender would write, without inventing facts.",
	"fix":       "Fix spelling, grammar and punctuation only. Change nothing else.",
	"formal":    "Rewrite the text in a formal, professional register.",
	"friendly":  "Rewrite the text in a warm, friendly register.",
	"translate": "Translate the text into %s.",
	"reply":     "Write the body of a reply to the latest email in the thread below, from the user's point of view. Sign it as the user would if their earlier messages show a signature.",
	"summarize": "Summarize the thread below in at most five short bullet points: who wants what, decisions, dates, and what the user still has to do.",
}

// Actions lists the helper's actions in menu order.
var Actions = []string{"rewrite", "shorten", "expand", "fix", "formal", "friendly", "translate", "reply", "summarize"}

func buildPrompt(req Request, persona string) (string, string) {
	system := baseSystem
	if persona != "" {
		system = persona + "\n\n" + system
	}
	task := actionPrompts[req.Action]
	if req.Action == "translate" {
		task = fmt.Sprintf(task, req.Language)
	}
	var user strings.Builder
	user.WriteString(task)
	if req.Instructions != "" {
		user.WriteString("\n\nThe user adds: " + req.Instructions)
	}
	if req.Context != "" {
		user.WriteString("\n\n<thread>\n" + req.Context + "\n</thread>")
	}
	if req.Text != "" {
		user.WriteString("\n\n<text>\n" + req.Text + "\n</text>")
	}
	return system, user.String()
}
