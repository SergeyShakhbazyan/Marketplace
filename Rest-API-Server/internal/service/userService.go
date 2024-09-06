package service

import (
	"context"
	"errors"
	"github.com/gocql/gocql"
	"marketplace_project/internal/models"
	"marketplace_project/internal/repository"
	"marketplace_project/internal/utils"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(user *models.User) error {
	existingUser, err := s.repo.GetUserByEmail(context.Background(), user.Email)
	if err != nil && !errors.Is(err, utils.ErrNotFound) {
		return err
	}
	if existingUser != nil {
		return utils.ErrEmailExists
	}

	return s.repo.CreateUser(context.Background(), user)
}

func (s *UserService) SignIn(email, password string) (*models.User, error) {
	user, err := s.repo.GetUserByEmail(context.Background(), email)
	if err != nil {
		return nil, err
	}

	if user.Password != password {
		return nil, errors.New("invalid password")
	}

	return user, nil
}

func (s *UserService) GetUserByID(userID gocql.UUID) (*models.UserWrapContent, error) {
	user, err := s.repo.GetUser(context.Background(), userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
