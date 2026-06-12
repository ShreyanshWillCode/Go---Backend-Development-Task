package service

import (
	"context"
	"fmt"
	"time"

	db "github.com/shreyxnsh/anyx-user-api/db/sqlc"
	"github.com/shreyxnsh/anyx-user-api/internal/logger"
	"github.com/shreyxnsh/anyx-user-api/internal/models"
	"github.com/shreyxnsh/anyx-user-api/internal/repository"
	"go.uber.org/zap"
)

type UserService interface {
	CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.UserResponse, error)
	GetUser(ctx context.Context, id int32) (*models.UserWithAgeResponse, error)
	UpdateUser(ctx context.Context, id int32, req models.UpdateUserRequest) (*models.UserResponse, error)
	DeleteUser(ctx context.Context, id int32) error
	ListUsers(ctx context.Context, page, limit int) (*models.PaginatedUsersResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.UserResponse, error) {
	dob, err := models.ParseDob(req.Dob)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	user, err := s.repo.Create(ctx, db.CreateUserParams{
		Name:	req.Name,
		Dob:	dob,
	})
	if err != nil {
		logger.Error("service: failed to create user", zap.Error(err))
		return nil, fmt.Errorf("could not create user: %w", err)
	}

	logger.Info("service: user created", zap.Int32("id", user.ID), zap.String("name", user.Name))
	return toUserResponse(user), nil
}

func (s *userService) GetUser(ctx context.Context, id int32) (*models.UserWithAgeResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	logger.Info("service: user fetched", zap.Int32("id", id))
	return toUserWithAgeResponse(user), nil
}

func (s *userService) UpdateUser(ctx context.Context, id int32, req models.UpdateUserRequest) (*models.UserResponse, error) {
	dob, err := models.ParseDob(req.Dob)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	user, err := s.repo.Update(ctx, db.UpdateUserParams{
		ID:	id,
		Name:	req.Name,
		Dob:	dob,
	})
	if err != nil {
		logger.Error("service: failed to update user", zap.Int32("id", id), zap.Error(err))
		return nil, fmt.Errorf("could not update user: %w", err)
	}

	logger.Info("service: user updated", zap.Int32("id", id))
	return toUserResponse(user), nil
}

func (s *userService) DeleteUser(ctx context.Context, id int32) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		logger.Error("service: failed to delete user", zap.Int32("id", id), zap.Error(err))
		return fmt.Errorf("could not delete user: %w", err)
	}
	logger.Info("service: user deleted", zap.Int32("id", id))
	return nil
}

func (s *userService) ListUsers(ctx context.Context, page, limit int) (*models.PaginatedUsersResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	users, err := s.repo.List(ctx, db.ListUsersParams{
		Limit:	int32(limit),
		Offset:	int32(offset),
	})
	if err != nil {
		logger.Error("service: failed to list users", zap.Error(err))
		return nil, fmt.Errorf("could not list users: %w", err)
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		logger.Error("service: failed to count users", zap.Error(err))
		return nil, fmt.Errorf("could not count users: %w", err)
	}

	result := make([]models.UserWithAgeResponse, 0, len(users))
	for _, u := range users {
		result = append(result, *toUserWithAgeResponse(u))
	}

	logger.Info("service: users listed", zap.Int("count", len(result)), zap.Int("page", page))
	return &models.PaginatedUsersResponse{
		Data:	result,
		Total:	total,
		Page:	page,
		Limit:	limit,
	}, nil
}

func CalculateAge(dob, now time.Time) int {
	years := now.Year() - dob.Year()

	dobThisYear := time.Date(now.Year(), dob.Month(), dob.Day(), 0, 0, 0, 0, now.Location())
	if now.Before(dobThisYear) {
		years--
	}

	return years
}

func toUserResponse(u db.User) *models.UserResponse {
	return &models.UserResponse{
		ID:	u.ID,
		Name:	u.Name,
		Dob:	u.Dob.Format(models.DateLayout),
	}
}

func toUserWithAgeResponse(u db.User) *models.UserWithAgeResponse {
	return &models.UserWithAgeResponse{
		ID:	u.ID,
		Name:	u.Name,
		Dob:	u.Dob.Format(models.DateLayout),
		Age:	CalculateAge(u.Dob, time.Now()),
	}
}
