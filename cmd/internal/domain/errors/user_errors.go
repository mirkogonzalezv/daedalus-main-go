package errors

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
)

func ErrUserNameRequired() *DomainError {
	return NewValidationError(UserNameRequired, "Nombre del usuario es requerido")
}

func ErrUserNotFound() *DomainError {
	return NewValidationError(UserNotFound, "Usuario no encontrado")
}

func ErrUserEmailRequired() *DomainError {
	return NewValidationError(UserEmailRequired, "Email del usuario es requerido")
}

func ErrUserPasswordRequired() *DomainError {
	return NewValidationError(UserPasswordRequired, "Password del usuario es requerido")
}

func ErrUserReadyExist() *DomainError {
	return NewValidationError(UserReadyExist, "Usuario ya existe")
}

func ErrUserRoleNotFound() *DomainError {
	return NewValidationError(UserRoleNotFound, "Role del usuario es requerido")
}

func ErrUserNotMatchRole() *DomainError {
	return NewValidationError(UserRoleNotFound, "Role no válido")
}

func ErrUserPasswordLengthError() *DomainError {
	return NewValidationError(UserPasswordLengthError, "Password debe tener al menos 8 carácteres")
}

func ErrUserHashPassword() *DomainError {
	return NewValidationError(UserHashPasswordError, "Error al hashear password")
}

func ErrUserIDRequired() *DomainError {
	return NewValidationError(UserIDRequiredError, "ID del usuario es requerido")
}

func ErrUserEmailReadyExist() *DomainError {
	return NewValidationError(UserEmailDuplicatedError, "Email ya se encuentra registrado")
}

func ErrUserDeleteError() *DomainError {
	return NewValidationError(UserErrorDelete, "Error al eliminar usuario")
}

func ErrListUserNotFound() *DomainError {
	return NewValidationError(UserErrorListNotFound, "No hay usuarios registrados")
}
