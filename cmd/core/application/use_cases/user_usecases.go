package usecases

import (
	"context"
	domain "daedalus-engine-go/cmd/core/domain/entities"
	domainErrors "daedalus-engine-go/cmd/core/domain/errors"
	"daedalus-engine-go/cmd/core/domain/repository"
	"errors"
	"unicode"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	repo repository.UserRepository
	log  *zap.Logger
}

// Inyectamos repo y logger para tener visibilidad dentro del caso de uso
func NewUserUseCase(repo repository.UserRepository, log *zap.Logger) *UserUseCase {
	return &UserUseCase{repo: repo, log: log}
}

// Definidmos las funciones del caso de uso
func (uc *UserUseCase) CreateUser(ctx context.Context, name, tenant_id, email, password, role string) (*domain.User, error) {
	if name == "" {
		return nil, domainErrors.ErrUserNameRequired()
	}

	if email == "" {
		return nil, domainErrors.ErrUserEmailRequired()
	}

	if password == "" {
		return nil, domainErrors.ErrUserPasswordRequired()
	}

	if role == "" {
		return nil, domainErrors.ErrUserRoleNotFound()
	}

	if role != "owner" && role != "admin" && role != "user" {
		return nil, domainErrors.ErrUserNotMatchRole()
	}

	existe, err := uc.repo.GetByEmailAndTenant(ctx, tenant_id, email)

	if err != nil {
		return nil, domainErrors.ErrUserNotFound()
	}

	uc.log.Info("Nuevo usuario: Verificando email único")
	if existe != nil {
		return nil, domainErrors.ErrUserReadyExist()
	}

	// Hash del password
	uc.log.Info("Nuevo usuario: Verificando calidad de la password")
	if err := validarPassword(password); err != nil {
		return nil, domainErrors.ErrUserPasswordLengthError()
	}

	hashed, err := hashPassword(password)
	if err != nil {
		return nil, domainErrors.ErrUserHashPassword()
	}

	// creamos la instancia domai de user
	newUser := &domain.User{
		ID:           uuid.NewString(),
		TenantID:     &tenant_id,
		Name:         name,
		Email:        email,
		PasswordHash: hashed,
		Role:         role,
		Status:       "active",
	}

	err = uc.repo.Create(ctx, newUser)

	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (uc *UserUseCase) ObtenerUsuarioPorId(ctx context.Context, id string) (*domain.User, error) {
	usuario, err := uc.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if usuario == nil {
		return nil, domainErrors.ErrUserNotFound()
	}

	return usuario, nil
}

func (uc *UserUseCase) ObtenerUsuarioPorEmailYTenant(ctx context.Context, tenant_id, email string) (*domain.User, error) {
	usuario, err := uc.repo.GetByEmailAndTenant(ctx, tenant_id, email)
	if err != nil {
		return nil, err
	}

	if usuario == nil {
		return nil, domainErrors.ErrUserNotFound()
	}

	return usuario, nil
}

func (uc *UserUseCase) ActualizarUsuario(ctx context.Context, u *domain.User) error {

	if u.Name == "" {
		return domainErrors.ErrUserNameRequired()
	}

	if u.ID == "" {
		return domainErrors.ErrUserIDRequired()
	}

	// Cambiamos email, pero si existe el mismo en db no se puede
	if u.Email != "" {
		existe, err := uc.repo.GetByEmailAndTenant(ctx, *u.TenantID, u.Email)
		if err != nil {
			return err
		}
		if existe != nil && existe.Email != u.Email {
			return domainErrors.ErrUserEmailReadyExist()
		}
	}

	return uc.repo.Update(ctx, u)
}

func (uc *UserUseCase) EliminarUsuario(ctx context.Context, id string) error {
	if id == "" {
		return domainErrors.ErrUserIDRequired()
	}
	err := uc.repo.Delete(ctx, id)
	if err != nil {
		return domainErrors.ErrUserDeleteError()
	}
	return nil
}

func validarPassword(pw string) error {
	var (
		tieneMayuscula bool
		tieneSimbolo   bool
	)

	if len(pw) < 8 {
		return errors.New("password debe tener 8 caracteres como mínimo")
	}

	for _, c := range pw {
		switch {
		case unicode.IsUpper(c):
			tieneMayuscula = true
		case unicode.IsSymbol(c) || unicode.IsPunct(c):
			tieneSimbolo = true
		}
	}

	if !tieneMayuscula {
		return errors.New("password debe tener una letra mayúscula")
	}

	if !tieneSimbolo {
		return errors.New("password debe tener al menos un carácter especial")
	}

	return nil
}

func hashPassword(pw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
