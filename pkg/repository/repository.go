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

type Task interface {
	GetOrCreateCategoryByName(name string) (int, error)
	GetLastParseTime() (string, error)
	SetLastParseTime() error
	GetAllForChannels() ([]core.ChannelTaskResponse, error)
	GetAllForUsers() ([]core.UserTaskResponse, error)
	AddTasks(tasksInput core.TasksInput) error
	DeleteAll() error
}

type Repository struct {
	Channel
	User
	Task
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Channel:  NewChannelPostgres(db),
		User:     NewUserPostgres(db),
		Task:     NewTaskPostgres(db),
	}
}

type SettingObject struct {
	Id         int  `db:"id"`
	IsSafeDeal bool `db:"is_safe_deal"`
	IsBudget   bool `db:"is_budget"`
	IsTerm     bool `db:"is_term"`
}

type ChannelObject struct {
	Id      int    `db:"id"`
	ApiId   int    `db:"api_id"`
	ApiHash string `db:"api_hash"`
	Name    string `db:"name"`
}

type UserObject struct {
	Id       int    `db:"id"`
	TgId     int    `db:"tg_id"`
	Username string `db:"username"`
}
