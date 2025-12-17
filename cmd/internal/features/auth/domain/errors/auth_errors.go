package errors

import (
	domainErrors "daedalus-engine-go/cmd/internal/pkg/errors"
)

// Definimos los errores normalizados para usuarios
const (
	AuthInvalidCredentials = "AUTH_001"
	AuthUserInactive       = "AUTH_002"
	AuthInternalServerErr  = "AUTH_003"
	AuthTokenExpired       = "AUTH_004"
	AuthTokenRevoked       = "AUTH_005"
)

func ErrInvalidCredentials() *domainErrors.DomainError {
	return domainErrors.NewValidationError(AuthInvalidCredentials, "Email o Password incorrectos")
}

func ErrUserInactive() *domainErrors.DomainError {
	return domainErrors.NewValidationError(AuthUserInactive, "Usuario inactivo")
}

func ErrAuthInternalServerError() *domainErrors.DomainError {
	return domainErrors.NewValidationError(AuthInternalServerErr, "Internal server error")
}

func ErrTokenExpired() *domainErrors.DomainError {
	return domainErrors.NewValidationError(AuthTokenExpired, "Refresh token expirado")
}

func ErrTokenRevoked() *domainErrors.DomainError {
	return domainErrors.NewValidationError(AuthTokenRevoked, "Token refresh revocado")
}
