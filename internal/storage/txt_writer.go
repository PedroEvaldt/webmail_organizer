package storage

import (
	"fmt"
	"os"

	"webmail_organizer/internal/model"
)

func SaveEmailsToFile(emails []model.Email, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", filename, err)
	}
	defer file.Close()

	for _, email := range emails {
		_, err := fmt.Fprintf(file,
		"UID: %v\nFrom: %s\nSubject: %s\nDate: %v\n--------------------\n",
			email.UID,
			email.From,
			email.Subject,
			email.Date,
		)
		if err != nil {
			return fmt.Errorf("error writing to file: %w", err)
		}
	}
	return nil
}
