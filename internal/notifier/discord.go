package notifier

import (
	"bytes"
	"fmt"
	"encoding/json"
	"net/http"

	"webmail_organizer/internal/model"
)

type webhookPayload struct {
	Embeds []embed `json:"embeds"`
}

type embed struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Color       int     `json:"color"`
	Fields      []field `json:"fields"`
}

type field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

func SendNewEmailsNotification(emails []model.Email, webhookURL string) error {
	fields := make([]field, 0, len(emails))
	for _, email := range emails {
		fields = append(fields, field{
			Name:   email.Subject,
			Value:  "De: " + email.From,
			Inline: false,
		})
	}

	payload := webhookPayload{
		Embeds: []embed{
			{
				Title:       fmt.Sprintf("%d email(s) novo(s)", len(emails)),
				Description: "Você tem mensagens não lidas na INBOX",
				Color:       5814783,
				Fields:      fields,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Failed marshaling json: %w", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("Failed posting: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("discord returned status %d", resp.StatusCode)
	}

	return nil
}
