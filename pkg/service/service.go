package service

import (
	core "github.com/max-sanch/BotFreelancer-core"
	"github.com/max-sanch/BotFreelancer-core/pkg/repository"
)

type Channel interface {
	GetByApiId(apiId int) (core.ChannelResponse, error)
	Create(channelInput core.ChannelInput) (int, error)
	Update(channelInput core.ChannelInput) (int, error)
	Delete(apiID int) error
}

type User interface {
	GetByTgId(tgId int) (core.UserResponse, error)
	Create(userInput core.UserInput) (int, error)
	Update(userInput core.UserInput) (int, error)
}

type Service struct {
	Channel
	User
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Channel: NewChannelService(repos),
		User:    NewUserService(repos),
	}
}
