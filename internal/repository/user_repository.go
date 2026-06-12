package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/shreyxnsh/anyx-user-api/db/sqlc"
)

type UserRepository interface {
	Create(ctx context.Context, params db.CreateUserParams) (db.User, error)
	GetByID(ctx context.Context, id int32) (db.User, error)
	Update(ctx context.Context, params db.UpdateUserParams) (db.User, error)
	Delete(ctx context.Context, id int32) error
	List(ctx context.Context, params db.ListUsersParams) ([]db.User, error)
	Count(ctx context.Context) (int64, error)
}

type pgUserRepository struct {
	queries *db.Queries
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &pgUserRepository{
		queries: db.New(pool),
	}
}

func (r *pgUserRepository) Create(ctx context.Context, params db.CreateUserParams) (db.User, error) {
	return r.queries.CreateUser(ctx, params)
}

func (r *pgUserRepository) GetByID(ctx context.Context, id int32) (db.User, error) {
	return r.queries.GetUser(ctx, id)
}

func (r *pgUserRepository) Update(ctx context.Context, params db.UpdateUserParams) (db.User, error) {
	return r.queries.UpdateUser(ctx, params)
}

func (r *pgUserRepository) Delete(ctx context.Context, id int32) error {
	return r.queries.DeleteUser(ctx, id)
}

func (r *pgUserRepository) List(ctx context.Context, params db.ListUsersParams) ([]db.User, error) {
	return r.queries.ListUsers(ctx, params)
}

func (r *pgUserRepository) Count(ctx context.Context) (int64, error) {
	return r.queries.CountUsers(ctx)
}
