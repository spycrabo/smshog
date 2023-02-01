package http

import "time"

type Message struct {
	Sender    string    `json:"sender"`
	Phone     string    `json:"phone"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
}
