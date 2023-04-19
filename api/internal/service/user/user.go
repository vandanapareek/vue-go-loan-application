package user

import (
	"fmt"
	"loan-api/internal/models"
)

type UserService interface {
	Login(username, password string) (*models.User, error)
	GetUserInfo(username string) (*models.UserInfo, error)
}

type userService struct {
	repo models.UserRepo
}

func NewUserService(r models.UserRepo) UserService {
	return &userService{
		repo: r,
	}
}

func (s *userService) Login(username, password string) (*models.User, error) {
	user, err := s.repo.GetByUserName(username)
	if err != nil {
		return nil, fmt.Errorf("Inexistant user")
	}
	if user.Password != password {
		return nil, fmt.Errorf("Wrong password")
	}
	return user, nil
}

func (s *userService) GetUserInfo(username string) (*models.UserInfo, error) {
	user, err := s.repo.GetUserInfo(username)
	if err != nil {
		return nil, fmt.Errorf("Inexistant user")
	}
	return user, nil
}
