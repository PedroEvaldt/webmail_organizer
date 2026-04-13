package main

import (
	"log"
	"os"
	"github.com/joho/godotenv"

	"webmail_organizer/internal/notifier"
	"webmail_organizer/internal/imapclient"
	"webmail_organizer/internal/storage"
)


func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loanding .env file")
	}

	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	webhookURL := os.Getenv("DISCORD_WEBHOOK_URL")

	config := imapclient.NewConfig("imap.inf.ufrgs.br", 993, username, password)

	emails, err := imapclient.FetchLatestEmails(10, config)
	if err != nil {
		log.Fatal(err)
	}

	seen_uids, err := storage.LoadSeenUIDS("data/seen_uids.txt")
	if err != nil {
		log.Fatal(err)
	}

	unseen_emails, err := storage.SaveSeenUIDs(emails, seen_uids, "data/seen_uids.txt")
	if err != nil {
		log.Fatal(err)
	}


	if len(unseen_emails) > 0{
		err = notifier.SendNewEmailsNotification(unseen_emails, webhookURL)
		if err != nil {
			log.Fatal(err)
		}

		err = storage.SaveEmailsToFile(unseen_emails, "data/emails.txt")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		print("Nenhum email novo")
	}

}
