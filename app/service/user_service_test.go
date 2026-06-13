package service

import (
	"context"
	"errors"
	"testing"

	"carro-ideal/app/models"

	"golang.org/x/crypto/bcrypt"
)

type fakeUserRepository struct {
	user   *models.User
	exists bool
}

func (r *fakeUserRepository) ExistsByEmail(_ context.Context, _ string) (bool, error) {
	return r.exists, nil
}

func (r *fakeUserRepository) Create(_ context.Context, user *models.User) error {
	user.ID = 1
	r.user = user
	return nil
}

func (r *fakeUserRepository) GetByEmail(_ context.Context, email string) (*models.User, error) {
	if r.user == nil || r.user.Email != email {
		return nil, errors.New("not found")
	}
	return r.user, nil
}

func (r *fakeUserRepository) GetByID(_ context.Context, id int64) (*models.User, error) {
	if r.user == nil || r.user.ID != id {
		return nil, errors.New("not found")
	}
	return r.user, nil
}

func (r *fakeUserRepository) Update(_ context.Context, user *models.User) error {
	r.user = user
	return nil
}

func (r *fakeUserRepository) Deactivate(_ context.Context, id int64) error {
	if r.user != nil && r.user.ID == id {
		r.user.Active = false
	}
	return nil
}

func TestUserServiceRegisterAndLogin(t *testing.T) {
	repo := &fakeUserRepository{}
	service := NewUserService(repo)

	user, err := service.Register(
		context.Background(),
		"  Maria Silva  ",
		"  MARIA@EXAMPLE.COM ",
		"strong-password",
		"strong-password",
	)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if user.Name != "Maria Silva" {
		t.Fatalf("Register() name = %q, want normalized name", user.Name)
	}
	if user.Email != "maria@example.com" {
		t.Fatalf("Register() email = %q, want normalized email", user.Email)
	}
	if user.PasswordHash == "strong-password" {
		t.Fatal("Register() stored a plaintext password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("strong-password")); err != nil {
		t.Fatalf("Register() password hash is invalid: %v", err)
	}

	loggedIn, err := service.Login(context.Background(), "MARIA@example.com", "strong-password")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if loggedIn.ID != user.ID {
		t.Fatalf("Login() user ID = %d, want %d", loggedIn.ID, user.ID)
	}
}

func TestUserServiceRejectsDuplicateEmail(t *testing.T) {
	service := NewUserService(&fakeUserRepository{exists: true})

	_, err := service.Register(
		context.Background(),
		"Maria Silva",
		"maria@example.com",
		"strong-password",
		"strong-password",
	)
	if !errors.Is(err, ErrEmailAlreadyUsed) {
		t.Fatalf("Register() error = %v, want ErrEmailAlreadyUsed", err)
	}
}

func TestUserServiceRejectsWeakPassword(t *testing.T) {
	service := NewUserService(&fakeUserRepository{})

	_, err := service.Register(
		context.Background(),
		"Maria Silva",
		"maria@example.com",
		"short",
		"short",
	)
	if err == nil {
		t.Fatal("Register() accepted a password shorter than 8 characters")
	}
}

func TestUserServiceRejectsMismatchedPasswords(t *testing.T) {
	service := NewUserService(&fakeUserRepository{})

	_, err := service.Register(
		context.Background(),
		"Maria Silva",
		"maria@example.com",
		"senha12345",
		"senha99999",
	)
	if err == nil {
		t.Fatal("Register() accepted mismatched passwords")
	}
}

func TestUserServiceRejectsInvalidEmail(t *testing.T) {
	service := NewUserService(&fakeUserRepository{})

	_, err := service.Register(
		context.Background(),
		"Maria Silva",
		"not-an-email",
		"senha12345",
		"senha12345",
	)
	if err == nil {
		t.Fatal("Register() accepted invalid email format")
	}
}

func TestUserServiceLoginWrongPassword(t *testing.T) {
	repo := &fakeUserRepository{}
	svc := NewUserService(repo)

	_, err := svc.Register(context.Background(), "Maria", "maria@example.com", "correctpass", "correctpass")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	_, err = svc.Login(context.Background(), "maria@example.com", "wrongpass")
	if err == nil {
		t.Fatal("Login() accepted wrong password")
	}
}

func TestUserServiceLoginInactiveUser(t *testing.T) {
	repo := &fakeUserRepository{
		user: &models.User{
			ID:           1,
			Email:        "inactive@example.com",
			PasswordHash: "$2a$10$invalid",
			Active:       false,
		},
	}
	service := NewUserService(repo)

	_, err := service.Login(context.Background(), "inactive@example.com", "anypassword")
	if err == nil {
		t.Fatal("Login() accepted inactive user")
	}
}

func TestUserServiceLoginEmptyCredentials(t *testing.T) {
	service := NewUserService(&fakeUserRepository{})

	_, err := service.Login(context.Background(), "", "")
	if err == nil {
		t.Fatal("Login() accepted empty credentials")
	}
}

func TestUserServiceGetByID(t *testing.T) {
	repo := &fakeUserRepository{}
	svc := NewUserService(repo)

	user, err := svc.Register(context.Background(), "João", "joao@example.com", "senha12345", "senha12345")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	got, err := svc.GetByID(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if got.Email != "joao@example.com" {
		t.Fatalf("GetByID() email = %q, want joao@example.com", got.Email)
	}
}

func TestIsEmailAlreadyUsed(t *testing.T) {
	if !IsEmailAlreadyUsed(ErrEmailAlreadyUsed) {
		t.Fatal("IsEmailAlreadyUsed() returned false for ErrEmailAlreadyUsed")
	}
	if IsEmailAlreadyUsed(errors.New("other error")) {
		t.Fatal("IsEmailAlreadyUsed() returned true for unrelated error")
	}
}
