package message

import "time"

type Message struct {
	From    string    `json:"from" example:"John Doe"`
	Date    time.Time `json:"date"`
	Message string    `json:"message" example:"Hello world"`
}
