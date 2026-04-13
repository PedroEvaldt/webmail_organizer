package storage

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"webmail_organizer/internal/model"
)

func LoadSeenUIDs(filepath string) (map[uint32]bool, error) {
	file, err := os.Open(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return map[uint32]bool{}, nil
		}
		return nil, fmt.Errorf("failed opening %s: %w", filepath, err)
	}
	defer file.Close()

	seenUIDs := make(map[uint32]bool)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		val, err := strconv.ParseUint(line, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("malformed uid %q in %s: %w", line, filepath, err)
		}
		seenUIDs[uint32(val)] = true
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading lines: %w", err)
	}
	return seenUIDs, nil
}

func FilterNewEmails(emails []model.Email, seenUIDs map[uint32]bool) []model.Email {
	var newEmails []model.Email
	for _, email := range emails {
		if !seenUIDs[email.UID] {
			newEmails = append(newEmails, email)
		}
	}
	return newEmails
}

func AppendSeenUIDs(emails []model.Email, filepath string) error {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed opening %s: %w", filepath, err)
	}
	defer file.Close()

	for _, email := range emails {
		if _, err = fmt.Fprintf(file, "%v\n", email.UID); err != nil {
			return fmt.Errorf("error writing uid to file: %w", err)
		}
	}
	return nil
}
