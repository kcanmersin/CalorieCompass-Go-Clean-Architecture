package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"CalorieCompass/internal/entity"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user entity.User) (int64, error) {
	query := `
        INSERT INTO auth.users (email, password, name, created_at, updated_at)
        VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        RETURNING id
    `

	var id int64
	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Email,
		user.Password,
		user.Name,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("create user error: %w", err)
	}

	return id, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	query := `
        SELECT id, email, password, name, created_at, updated_at
        FROM auth.users
        WHERE email = $1
    `

	var user entity.User
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, fmt.Errorf("user not found")
		}
		return entity.User{}, fmt.Errorf("get user by email error: %w", err)
	}

	return user, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (entity.User, error) {
	query := `
        SELECT id, email, password, name, created_at, updated_at
        FROM auth.users
        WHERE id = $1
    `

	var user entity.User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, fmt.Errorf("user not found")
		}
		return entity.User{}, fmt.Errorf("get user by id error: %w", err)
	}

	return user, nil
}