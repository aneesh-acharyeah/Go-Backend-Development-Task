package repository

import (
	"context"
	"errors"
	"time"

	"github.com/anish/backend-development-task/db/sqlc"
	"github.com/jackc/pgx/v5"
)

var ErrNotFound = errors.New("user not found")

type User = sqlc.User

type UserRepository interface {
	CreateUser(ctx context.Context, name string, dob time.Time) (User, error)
	GetUserByID(ctx context.Context, id int32) (User, error)
	UpdateUser(ctx context.Context, id int32, name string, dob time.Time) (User, error)
	DeleteUser(ctx context.Context, id int32) error
	ListUsers(ctx context.Context) ([]User, error)
}

type userRepository struct {
	queries *sqlc.Queries
}

func NewUserRepository(queries *sqlc.Queries) UserRepository {
	return &userRepository{queries: queries}
}

func (r *userRepository) CreateUser(ctx context.Context, name string, dob time.Time) (User, error) {
	return r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Name: name,
		Dob:  dob,
	})
}

func (r *userRepository) GetUserByID(ctx context.Context, id int32) (User, error) {
	user, err := r.queries.GetUserByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}
	return user, err
}

func (r *userRepository) UpdateUser(ctx context.Context, id int32, name string, dob time.Time) (User, error) {
	user, err := r.queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:   id,
		Name: name,
		Dob:  dob,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}
	return user, err
}

func (r *userRepository) DeleteUser(ctx context.Context, id int32) error {
	_, err := r.queries.DeleteUser(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

func (r *userRepository) ListUsers(ctx context.Context) ([]User, error) {
	return r.queries.ListUsers(ctx)
}
