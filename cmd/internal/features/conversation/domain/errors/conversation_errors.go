package errors

import (
	domainErrors "daedalus-engine-go/cmd/internal/pkg/errors"
)

const (
	ConversarionInternalServerError = "CONVERSATION_001"
	ConversationUnauthorizedError   = "CONVERSATION_002"
	ConversationTitleRequired       = "CONVERSATION_003"
	ConversationForbidden           = "CONVERSATION_004"
	ConversationInvalidParameters   = "CONVERSATION_005"
	ConversationNotFound            = "CONVERSATION_006"
)

func InternalServerError() *domainErrors.DomainError {
	return domainErrors.NewValidationError(ConversarionInternalServerError, "Internal server error")
}

func UnauthorizedError() *domainErrors.DomainError {
	return domainErrors.NewValidationError(ConversationUnauthorizedError, "Unauthorized")
}

func TitleRequiredError() *domainErrors.DomainError {
	return domainErrors.NewValidationError(ConversationTitleRequired, "Título es requerido")
}

func ForbiddenError() *domainErrors.DomainError {
	return domainErrors.NewValidationError(ConversationForbidden, "Forbidden")
}

func InvalidParametersError() *domainErrors.DomainError {
	return domainErrors.NewValidationError(ConversationInvalidParameters, "Invalid parameters")
}

func NotFoundError() *domainErrors.DomainError {
	return domainErrors.NewValidationError(ConversationNotFound, "Not Found")
}
