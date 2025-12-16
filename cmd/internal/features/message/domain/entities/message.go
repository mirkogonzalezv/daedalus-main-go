package domain

import "time"

type Message struct {
	ID             string
	ConversationId string
	UserId         *string // nullable: nil if sender is AI/system
	Sender         string  // "user" | "ai"
	Content        string
	CreatedAt      time.Time
}
