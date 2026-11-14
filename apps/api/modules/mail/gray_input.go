package mail

import (
	"encoding/json"
	"strings"
	"unicode/utf8"

	"github.com/FacileStudio/Courrier/apps/api/schemas"
)

// grayInput is the JSON list Bulle reads, dropping the oldest emails that
// would push it past Bulle's input cap, or push the request past Bulle's
// body cap once the list is encoded again as a JSON string.
func grayInput(emails []schemas.Email) (string, error) {
	list := make([]grayEmail, 0, len(emails))
	size, wire := 2, 2
	for i, email := range emails {
		item := grayEmail{N: i, From: email.FromAddress, Name: email.FromName, Subject: clampText(email.Subject, 300), Body: clampText(email.BodyText, grayBodyBytes)}
		encoded, err := json.Marshal(item)
		if err != nil {
			return "", err
		}
		quoted, err := json.Marshal(string(encoded))
		if err != nil {
			return "", err
		}
		if size+len(encoded)+1 > grayInputBytes || wire+len(quoted)-1 > grayWireBytes {
			break
		}
		size += len(encoded) + 1
		wire += len(quoted) - 1
		list = append(list, item)
	}
	data, err := json.Marshal(list)
	return string(data), err
}

// sameSenderDomain keeps, in order and up to limit, the emails from the
// sender domain of the first one.
func sameSenderDomain(emails []schemas.Email, limit int) []schemas.Email {
	if len(emails) == 0 {
		return nil
	}
	domain := senderDomain(emails[0].FromAddress)
	kept := make([]schemas.Email, 0, min(limit, len(emails)))
	for _, email := range emails {
		if len(kept) < limit && senderDomain(email.FromAddress) == domain {
			kept = append(kept, email)
		}
	}
	return kept
}

func senderDomain(address string) string {
	address = strings.ToLower(strings.TrimSpace(address))
	return address[strings.LastIndexByte(address, '@')+1:]
}

func clampText(text string, limit int) string {
	text = strings.TrimSpace(text)
	if len(text) <= limit {
		return text
	}
	cut := limit
	for cut > 0 && !utf8.RuneStart(text[cut]) {
		cut--
	}
	return text[:cut]
}
