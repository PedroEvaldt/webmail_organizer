package model

import (
	"time"
)

type Email struct {
	UID     uint32
	Subject string
	From    string
	Date    time.Time
}
