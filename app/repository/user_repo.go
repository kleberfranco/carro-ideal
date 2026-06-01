package repository

import (
	"context"
	"database/sql"

	"carro-ideal/app/models"
)

type UserRepository interface {
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, u *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", email).Scan(&exists)
	return exists, err
}

func (r *userRepository) Create(ctx context.Context, u *models.User) error {
	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO users (name, email, password) VALUES ($1,$2,$3)",
		u.Name, u.Email, u.Password,
	)
	return err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRowContext(
		ctx,
		"SELECT id, name, email, password, created_at FROM users WHERE email=$1",
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}
