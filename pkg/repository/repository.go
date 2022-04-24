package repository

import (
	"github.com/jmoiron/sqlx"
	core "github.com/max-sanch/BotFreelancer-core"
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

type Repository struct {
	Channel
	User
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Channel: NewChannelPostgres(db),
		User:    NewUserPostgres(db),
	}
}

type SettingObject struct {
	Id         int  `json:"id"`
	IsSafeDeal bool `json:"is_safe_deal"`
	IsBudget   bool `json:"is_budget"`
	IsTerm     bool `json:"is_term"`
}

type ChannelObject struct {
	Id      int    `json:"id"`
	ApiId   int    `json:"api_id"`
	ApiHash string `json:"api_hash"`
	Name    string `json:"name"`
}

type UserObject struct {
	Id       int    `json:"id"`
	TgId     int    `json:"tg_id"`
	Username string `json:"username"`
}
