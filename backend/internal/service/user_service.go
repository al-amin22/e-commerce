package service

import (
	"context"
	"strings"

	"ecommerce-backend/internal/domain"
	"ecommerce-backend/internal/repository"
)

type UserService interface {
	Create(ctx context.Context, user *domain.User) error
	List(ctx context.Context) ([]domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(ctx context.Context, user *domain.User) error {
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	if strings.TrimSpace(user.Role) == "" {
		user.Role = "buyer"
	}
	return s.repo.Create(ctx, user)
}

func (s *userService) List(ctx context.Context) ([]domain.User, error) {
	return s.repo.List(ctx)
}
