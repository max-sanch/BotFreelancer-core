package service

import (
	core "github.com/max-sanch/BotFreelancer-core"
	"github.com/max-sanch/BotFreelancer-core/pkg/repository"
)

type UserService struct {
	repo *repository.Repository
}

func (s *UserService) GetTasks() ([]core.UserTaskResponse, error) {
	tasks, err := s.repo.Task.GetAllForUsers()
	if err != nil {
		return nil, err
	}

	return tasks, nil
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
