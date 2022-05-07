package repository

import (
	"os"

	core "github.com/max-sanch/BotFreelancer-core"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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

func NewPostgresRepos() *Repository {
	db, err := newPostgresDB(Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("failed initialize postgres database: %s", err.Error())
	}

	return &Repository{
		Channel: NewChannelPostgres(db),
		User:    NewUserPostgres(db),
		Task:    NewTaskPostgres(db),
	}
}
