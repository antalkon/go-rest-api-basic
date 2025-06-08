package models

import "time"

type Ping struct {
	ID        int
	Timestamp time.Time
	IP        string
}
