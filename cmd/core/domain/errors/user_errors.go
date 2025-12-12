package errors

// Definimos los errores normalizados para usuarios
const (
	UserNameRequired = "USER_001"
	UserNotFound     = "USER_002"
)

func ErrUserNameRequired() *DomainError {
	return NewValidationError(UserNameRequired, "Nombre del usuario es requerido")
}

func ErrUserNotFound() *DomainError {
	return NewValidationError(UserNotFound, "Usuario no encontrado")
}
