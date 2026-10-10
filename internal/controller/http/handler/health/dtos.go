package health

import "time"

type Response struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

type RequestHeaders struct {
	XCallerId string `header:"x-Caller-id" binding:"required"`
	XTestId   string `header:"X-Test-Id"`
}
