package service

import (
	core "github.com/max-sanch/BotFreelancer-core"
	"github.com/max-sanch/BotFreelancer-core/pkg/repository"
)

type Channel interface {
	GetChannel(apiID int) (core.ChannelResponse, error)
	CreateChannel(channelInput core.ChannelInput) (int, error)
	UpdateChannel(channelInput core.ChannelInput) (int, error)
	DeleteChannel(apiID int) error
}

type User interface {
	GetUser(tgID int) (core.UserResponse, error)
	CreateUser(userInput core.UserInput) (int, error)
	UpdateUser(userInput core.UserInput) (int, error)
}

type Parse interface {

}

type Service struct {
	Channel
	User
	Parse
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Channel: NewChannelService(repos.Channel),
		User: NewUserService(repos.User),
	}
}