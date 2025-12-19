package requests

type CreateConversationRequestDTO struct {
	Title string `json:"title" binding:"required,min=1,max=255"`
}

type UpdateConversationRequestDTO struct {
	Title string `json:"title" binding:"required,min=1,max=255"`
}
