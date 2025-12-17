package errors

import (
	domainErrors "daedalus-engine-go/cmd/internal/pkg/errors"
)

// Definimos los errores normalizados para usuarios
const (
	UserNameRequired         = "USER_001"
	UserNotFound             = "USER_002"
	UserEmailRequired        = "USER_003"
	UserPasswordRequired     = "USER_004"
	UserReadyExist           = "USER_006"
	UserRoleNotFound         = "USER_007"
	UserRoleNotMatchError    = "USER_008"
	UserPasswordLengthError  = "USER_009"
	UserHashPasswordError    = "USER_010"
	UserIDRequiredError      = "USER_011"
	UserEmailDuplicatedError = "USER_012"
	UserErrorDelete          = "USER_013"
	UserErrorListNotFound    = "USER_014"
	UserErrorInactive        = "USER_015"
)

func ErrUserNameRequired() *domainErrors.DomainError {
	return domainErrors.NewValidationError(UserNameRequired, "Nombre del usuario es requerido")
}

func ErrUserNotFound() *domainErrors.DomainError {
	return domainErrors.NewValidationError(UserNotFound, "Usuario no encontrado")
}

func ErrUserEmailRequired() *domainErrors.DomainError {
	return domainErrors.NewValidationError(UserEmailRequired, "Email del usuario es requerido")
}

func ErrUserPasswordRequired() *domainErrors.DomainError {
	return domainErrors.NewValidationError(UserPasswordRequired, "Password del usuario es requerido")
}

func ErrUserReadyExist() *domainErrors.DomainError {
	return domainErrors.NewValidationError(UserReadyExist, "Usuario ya existe")
}

func ErrUserRoleNotFound() *domainErrors.DomainError {
	return domainErrors.NewValidationError(UserRoleNotFound, "Role del usuario es requerido")
}

func ErrUserNotMatchRole() *domainErrors.DomainError {
	return domainErrors.NewValidationError(UserRoleNotFound, "Role no válido")
}

func ErrUserPasswordLengthError() *domainErrors.DomainError {
	return domainErrors.NewValidationError(UserPasswordLengthError, "Password debe tener al menos 8 carácteres")
}

func ErrUserHashPassword() *domainErrors.DomainError {
	return domainErrors.NewValidationError(UserHashPasswordError, "Error al hashear password")
}

func ErrUserIDRequired() *domainErrors.DomainError {
	return domainErrors.NewValidationError(UserIDRequiredError, "ID del usuario es requerido")
}

func ErrUserEmailReadyExist() *domainErrors.DomainError {
	return domainErrors.NewValidationError(UserEmailDuplicatedError, "Email ya se encuentra registrado")
}

func ErrUserDeleteError() *domainErrors.DomainError {
	return domainErrors.NewValidationError(UserErrorDelete, "Error al eliminar usuario")
}

func ErrListUserNotFound() *domainErrors.DomainError {
	return domainErrors.NewValidationError(UserErrorListNotFound, "No hay usuarios registrados")
}

func ErrUserInactive() *domainErrors.DomainError {
	return domainErrors.NewValidationError(UserErrorInactive, "Usuario inactivo")
}
