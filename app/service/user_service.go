package service

import (
	"context"
	"errors"
	"net/mail"
	"strings"

	"carro-ideal/app/models"
	"carro-ideal/app/repository"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailAlreadyUsed = errors.New("email já está em uso")
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Register(ctx context.Context, name, email, password, confirmPassword string) (*models.User, error) {
	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))

	if name == "" || email == "" || password == "" || confirmPassword == "" {
		return nil, errors.New("todos os campos são obrigatórios")
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return nil, errors.New("email inválido")
	}

	if password != confirmPassword {
		return nil, errors.New("as senhas não conferem")
	}

	if len(password) < 8 {
		return nil, errors.New("a senha deve ter pelo menos 8 caracteres")
	}

	exists, err := s.repo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, ErrEmailAlreadyUsed
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hashed),
		Role:         "user",
		Active:       true,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (*models.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || password == "" {
		return nil, errors.New("email e senha são obrigatórios")
	}

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("email ou senha incorretos")
	}

	if !user.Active {
		return nil, errors.New("usuário inativo")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("email ou senha incorretos")
	}

	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

func IsEmailAlreadyUsed(err error) bool {
	return errors.Is(err, ErrEmailAlreadyUsed)
}
