package service

import (
	core "github.com/max-sanch/BotFreelancer-core"
	"github.com/max-sanch/BotFreelancer-core/pkg/repository"
)

type UserService struct {
	repo *repository.Repository
}

func NewUserService(repo *repository.Repository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetByTgId(tgId int) (core.UserResponse, error) {
	return s.repo.User.GetByTgId(tgId)
}

func (s *UserService) Create(userInput core.UserInput) (int, error) {
	return s.repo.User.Create(userInput)
}

func (s *UserService) Update(userInput core.UserInput) (int, error) {
	return s.repo.User.Update(userInput)
}
