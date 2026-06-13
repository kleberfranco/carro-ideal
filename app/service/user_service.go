package service

import (
	"context"
	"errors"

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

func (s *UserService) Register(ctx context.Context, name, email, password, confirmPassword string) error {
	if name == "" || email == "" || password == "" || confirmPassword == "" {
		return errors.New("todos os campos são obrigatórios")
	}

	if password != confirmPassword {
		return errors.New("as senhas não conferem")
	}

	exists, err := s.repo.ExistsByEmail(ctx, email)
	if err != nil {
		return err
	}

	if exists {
		return ErrEmailAlreadyUsed
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hashed),
		Role:         "user",
		Active:       true,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (*models.User, error) {
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
