package repository

import (
	"github.com/jmoiron/sqlx"
	core "github.com/max-sanch/BotFreelancer-core"
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

type Repository struct {
	Channel
	User
	Parse
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Channel: NewChannelPostgres(db),
		User: NewUserPostgres(db),
	}
}
