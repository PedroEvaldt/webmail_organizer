package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"webmail_organizer/internal/imapclient"
	"webmail_organizer/internal/notifier"
	"webmail_organizer/internal/storage"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	username := os.Getenv("IMAP_USERNAME")
	password := os.Getenv("PASSWORD")
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")

	seenUIDs, err := storage.LoadSeenUIDs("data/seen_uids.txt")
	if err != nil {
		log.Fatal(err)
	}

	config := imapclient.NewConfig("imap.inf.ufrgs.br", 993, username, password)

	emails, err := imapclient.FetchLatestEmails(config)
	if err != nil {
		log.Fatal(err)
	}

	newEmails := storage.FilterNewEmails(emails, seenUIDs)

	if len(newEmails) == 0 {
		log.Println("Nenhum email novo")
		return
	}

	err = notifier.SendNewEmailsNotification(newEmails, webhookURL)
	if err != nil {
		log.Fatal(err)
	}

	err = storage.AppendSeenUIDs(newEmails, "data/seen_uids.txt")
	if err != nil {
		log.Fatal(err)
	}

	err = storage.SaveEmailsToFile(newEmails, "data/emails.txt")
	if err != nil {
		log.Fatal(err)
	}
}
