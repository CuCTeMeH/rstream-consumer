package consumer

import "time"

type Message struct {
	MessageID       string `json:"message_id"`
	RandomProperty  string `json:"random_property"`
	ConsumerID      string `json:"consumer_id"`
	ProcessingStart time.Time
	ProcessingEnd   time.Time
}
