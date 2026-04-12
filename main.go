package main

import (
	"log"
	"os"
	"github.com/joho/godotenv"

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

	config := imapclient.NewConfig("imap.inf.ufrgs.br", 993, username, password)

	emails, err := imapclient.FetchLatestEmails(10, config)
	if err != nil {
		log.Fatal(err)
	}

	seen_uids, err := storage.LoadSeenUIDS("data/seen_uids.txt")
	if err != nil {
		log.Fatal(err)
	}

	err = storage.SaveSeenUIDs(emails, seen_uids, "data/seen_uids.txt")
	if err != nil {
		log.Fatal(err)
	}

	err = storage.SaveEmailsToFile(emails, "data/emails.txt")
	if err != nil {
		log.Fatal(err)
	}
}
