package entity

import "time"

type SMS struct {
	Phone    string
	Message  string
	SentDate time.Time
}
