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
	GetByID(ctx context.Context, id int64) (*models.User, error)
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
	row := r.db.QueryRowContext(
		ctx,
		"INSERT INTO users (name, email, password_hash, role, active) VALUES ($1,$2,$3,$4,$5) RETURNING id, created_at, updated_at",
		u.Name, u.Email, u.PasswordHash, u.Role, u.Active,
	)
	return row.Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRowContext(
		ctx,
		"SELECT id, name, email, password_hash, role, active, created_at, updated_at FROM users WHERE email=$1",
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.Active, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRowContext(
		ctx,
		"SELECT id, name, email, password_hash, role, active, created_at, updated_at FROM users WHERE id=$1",
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.Active, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}
