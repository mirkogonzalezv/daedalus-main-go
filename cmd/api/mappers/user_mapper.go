package mappers

import (
	"daedalus-engine-go/cmd/api/dtos/responses"
	domain "daedalus-engine-go/cmd/core/domain/entities"
)

func ToUserResponse(user *domain.User) responses.UserResponse {
	return responses.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToCreateUserResponse(user *domain.User) responses.CreateUserResponse {
	return responses.CreateUserResponse{
		User:    ToUserResponse(user),
		Message: "Usuario creado exitosamente",
	}
}

func ToGetUserResponse(user *domain.User) responses.GetUserResponse {
	return responses.GetUserResponse{
		User: ToUserResponse(user),
	}
}

func ToUpdateUserResponse(user *domain.User) responses.UpdateUserResponse {
	return responses.UpdateUserResponse{
		User:    ToUserResponse(user),
		Message: "Usuario actualizado exitosamente",
	}
}
