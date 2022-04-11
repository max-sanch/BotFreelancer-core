package service

import (
	core "github.com/max-sanch/BotFreelancer-core"
	"github.com/max-sanch/BotFreelancer-core/pkg/repository"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(tgID int) (core.UserResponse, error) {
	return s.repo.GetUser(tgID)
}

func (s *UserService) CreateUser(userInput core.UserInput) (int, error) {
	return s.repo.CreateUser(userInput)
}

func (s *UserService) UpdateUser(userInput core.UserInput) (int, error) {
	return s.repo.UpdateUser(userInput)
}
