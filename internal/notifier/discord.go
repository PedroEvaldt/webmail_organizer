package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"webmail_organizer/internal/model"
)

const (
	// Limites da API do Discord
	maxDescriptionLen = 4096
	maxTitleLen       = 256

	// Limites estéticos: assunto muito longo quebra o layout da lista
	maxSubjectLen   = 110
	maxSenderLen    = 60
	maxEmailsListed = 12

	footerText = "📥 INBOX · imap.inf.ufrgs.br"

	colorGreen  = 5763719  // 1-2 emails
	colorYellow = 16705372 // 3-5 emails
	colorRed    = 15548997 // 6+ emails
)

type webhookPayload struct {
	Embeds []embed `json:"embeds"`
}

type embed struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Color       int     `json:"color"`
	Timestamp   string  `json:"timestamp,omitempty"`
	Footer      *footer `json:"footer,omitempty"`
}

type footer struct {
	Text string `json:"text"`
}

func SendNewEmailsNotification(emails []model.Email, webhookURL string) error {
	if len(emails) == 0 {
		return nil
	}

	payload := webhookPayload{
		Embeds: []embed{
			{
				Title:       buildTitle(len(emails)),
				Description: buildDescription(emails),
				Color:       colorFor(len(emails)),
				Timestamp:   time.Now().Format(time.RFC3339),
				Footer:      &footer{Text: footerText},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed marshaling json: %w", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed posting: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("discord returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func buildTitle(count int) string {
	if count == 1 {
		return truncate("📬  1 novo email", maxTitleLen)
	}
	return truncate(fmt.Sprintf("📬  %d novos emails", count), maxTitleLen)
}

func colorFor(count int) int {
	switch {
	case count <= 2:
		return colorGreen
	case count <= 5:
		return colorYellow
	default:
		return colorRed
	}
}

// buildDescription monta a lista de emails, uma entrada por email:
//
//	✉ **Assunto**
//	└ *Remetente* · há 2 horas
func buildDescription(emails []model.Email) string {
	var b strings.Builder

	listed := emails
	if len(listed) > maxEmailsListed {
		listed = listed[:maxEmailsListed]
	}

	for i, email := range listed {
		entry := formatEntry(email)

		// Reserva espaço para a linha final "+ N outros"
		if b.Len()+len(entry)+80 > maxDescriptionLen {
			listed = emails[:i]
			break
		}

		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(entry)
	}

	if remaining := len(emails) - len(listed); remaining > 0 {
		fmt.Fprintf(&b, "\n*… e mais %d email(s) não listado(s)*", remaining)
	}

	return b.String()
}

func formatEntry(email model.Email) string {
	subject := strings.TrimSpace(email.Subject)
	if subject == "" {
		subject = "(sem assunto)"
	}
	subject = escapeMarkdown(truncate(subject, maxSubjectLen))

	sender := escapeMarkdown(truncate(senderName(email.From), maxSenderLen))

	line := fmt.Sprintf("✉ **%s**\n└ *%s*", subject, sender)
	if !email.Date.IsZero() {
		// Timestamp relativo nativo do Discord: renderiza "há 2 horas"
		// no fuso de quem está lendo.
		line += fmt.Sprintf(" · <t:%d:R>", email.Date.Unix())
	}

	return line + "\n"
}

// senderName extrai o nome amigável de "Nome <caixa@host>".
// Sem nome, devolve o próprio endereço.
func senderName(from string) string {
	from = strings.TrimSpace(from)
	if from == "" {
		return "remetente desconhecido"
	}

	if i := strings.LastIndex(from, " <"); i > 0 {
		return strings.Trim(strings.TrimSpace(from[:i]), `"`)
	}

	return from
}

// escapeMarkdown neutraliza formatação em assuntos/remetentes vindos de fora,
// para que um "*promoção*" no assunto não vire itálico no embed.
func escapeMarkdown(s string) string {
	replacer := strings.NewReplacer(
		`\`, `\\`,
		"*", `\*`,
		"_", `\_`,
		"~", `\~`,
		"`", "\\`",
		"|", `\|`,
		"#", `\#`,
		">", `\>`,
		"[", `\[`,
		"]", `\]`,
	)
	return replacer.Replace(s)
}

func truncate(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return strings.TrimSpace(string(runes[:max-1])) + "…"
}
