package imapclient

import (
	"fmt"
	"log"

	"github.com/emersion/go-imap/v2"
	imaplib "github.com/emersion/go-imap/v2/imapclient"
	"webmail_organizer/internal/model"
)


type Config struct {
	Host string
	Port int
	Username string
	Password string
}


func NewConfig(host string, port int, username, password string) Config {
	return Config{Host: host, Port: port, Username: username, Password: password}
}

func FetchLatestEmails(limit int, config Config) ([]model.Email, error) {
	// Conectar no servidor
	serverIP := fmt.Sprintf("%s:%d", config.Host, config.Port)
	c, err := connect(serverIP)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// Fazer login
	err = login(c, config)
	if err != nil {
		return nil, err
	}

	var emails []model.Email
	// Buscar emails
	messages, err := fetchEmails(c, "INBOX", limit)
	if err != nil {
		return nil, err
	}

	// Ler dados (remetente, assunto, data)
	for _, message := range messages { var from string

		if len(message.Envelope.From) > 0 {
			address := message.Envelope.From[0]
			if address.Name != "" {
				from = address.Name + " <" + address.Mailbox + "@" + address.Host + ">"
			} else {
				from = address.Mailbox + "@" + address.Host
			}
		}

		email := model.Email{
			UID:     uint32(message.UID),
			Subject: message.Envelope.Subject,
			From: from,
			Date: message.Envelope.Date,
		}
		emails = append(emails, email)
	}

	if err := c.Logout().Wait(); err != nil {
		return nil, fmt.Errorf("failed to logout: %w", err)
	}

	// Processar (Salvar em model.Email) e enviar
	return emails, nil
}

func connect(serverIP string) (*imaplib.Client, error) {
	c, err := imaplib.DialTLS(serverIP, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to dial IMAP server: %w", err)
	}
	log.Println("Connected to server")
	return c, nil
}

func login(c *imaplib.Client, config Config) (error) {
	if err := c.Login(config.Username, config.Password).Wait(); err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}
	log.Println("Login made succefuly")
	return nil
}

func fetchEmails(c *imaplib.Client, mbox string, emailLimit int) ([]*imaplib.FetchMessageBuffer, error) {
	var messages []*imaplib.FetchMessageBuffer

	selectedMbox, err := c.Select(mbox, nil).Wait()
	if err != nil {
		return nil, fmt.Errorf("failed to select %s: %w", mbox, err)
	}

	criteria := &imap.SearchCriteria{
		NotFlag: []imap.Flag{imap.FlagSeen},
	}
	if selectedMbox.NumMessages > 0 {

		searchData, err := c.UIDSearch(criteria, nil).Wait()
		if err != nil {
			return nil, fmt.Errorf("failed to catch uids: %w", err)
		}

		uidSet := imap.UIDSetNum(searchData.AllUIDs()...)

		fetchOptions := &imap.FetchOptions{Envelope: true, Flags: true}
		messages, err = c.Fetch(uidSet, fetchOptions).Collect()
		if err != nil {
			return nil, fmt.Errorf("failed fetching email: %w", err)
		}
	}
	return messages, nil
}
