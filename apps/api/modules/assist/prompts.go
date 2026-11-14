package assist

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/FacileStudio/tronc/bulle"
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

// buildPrompt is the Bulle call for one action: the task and the user's own
// note as the prompt, the thread and the text as untrusted input.
func buildPrompt(req Request) bulle.CompleteRequest {
	task := actionPrompts[req.Action]
	if req.Action == "translate" {
		task = fmt.Sprintf(task, req.Language)
	}
	var prompt strings.Builder
	prompt.WriteString(task)
	if req.Instructions != "" {
		prompt.WriteString("\n\nThe user adds: " + req.Instructions)
	}
	prompt.WriteString("\n\nThe input holds the email thread in <thread> and the user's text in <text>, when there are any.")
	var input strings.Builder
	if req.Context != "" {
		input.WriteString("<thread>\n" + req.Context + "\n</thread>")
	}
	if req.Text != "" {
		if input.Len() > 0 {
			input.WriteString("\n\n")
		}
		input.WriteString("<text>\n" + req.Text + "\n</text>")
	}
	return bulle.CompleteRequest{System: baseSystem, Prompt: prompt.String(), Input: input.String(), MaxTokens: answerTokens}
}

// tooLarge says whether Bulle would refuse call for its size: the input past
// its input cap, or the request past its body cap once encoded, where every
// quote and angle bracket in the text grows.
func tooLarge(call bulle.CompleteRequest) bool {
	encoded, err := json.Marshal(call)
	return err != nil || len(call.Input) > maxInput || len(encoded) > maxBody
}
