package domain

import "time"

type Status string

const (
	StatusUP Status = "UP"
)

type Health struct {
	Status    Status
	Timestamp time.Time
	Version   string
}
