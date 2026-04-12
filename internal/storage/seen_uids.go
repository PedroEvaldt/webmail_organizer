package storage

import (
	"os"
	"fmt"
	"bufio"
	"strconv"

	"webmail_organizer/internal/model"
)

func LoadSeenUIDS(filepath string) (map[uint32]bool, error){
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed opening %s: %w", filepath, err)
	}
	defer file.Close()

	seen_uids := make(map[uint32]bool)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		val, _ := strconv.ParseUint(line, 10, 32)
		uid := uint32(val)
		seen_uids[uid] = true
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error reading lines: %w", err)
	}
	return seen_uids, nil
}

func SaveSeenUIDs(emails []model.Email, seen_uids map[uint32]bool, filepath string) error {
	//Salvar seen_uids no filepath
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Failed opening %s: %w", filepath, err)
	}
	defer file.Close()

	for i, email := range emails {
		if !(seen_uids[email.UID]){
			_, err = fmt.Fprintf(file, "%v\n", email.UID)
			if err != nil {
				return fmt.Errorf("Erro trying to write uid in file: %w", err)
			}
		} else {
			emails[i].Seen = true
		}
	}
	return nil
}

