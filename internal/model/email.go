package model

import (
	"time"
)

type Email struct {
	UID uint32
	Seen bool
	Subject string
	From string
	Date time.Time }
