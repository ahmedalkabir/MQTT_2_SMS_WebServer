package entity

import "time"

type User struct {
	ID          string
	Email       string
	Password    string
	Created_At  time.Time
	Updated_At  time.Time
	Name        string
	Phone       string
	Is_Disabled bool
}
