package responses

import "time"

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListUsersResponse struct {
	Users []UserResponse `json:"users"`
}

type CreateUserResponse struct {
	User    UserResponse `json:"user"`
	Message string       `json:"message"`
}

type GetUserResponse struct {
	User UserResponse `json:"user"`
}

type UpdateUserResponse struct {
	User    UserResponse `json:"user"`
	Message string       `json:"message"`
}
