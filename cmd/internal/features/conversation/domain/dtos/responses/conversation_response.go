package responses

import "time"

type ConversationResponseDTO struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ConversationListResponseDTO struct {
	Conversations []ConversationResponseDTO `json:"conversations"`
	Total         int                       `json:"total"`
	Page          int                       `json:"page"`
	PageSize      int                       `json:"page_size"`
}
