package repository

import (
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable             = "users"
	userSettingsTable      = "user_settings"
	userCategoriesTable    = "user_categories"
	channelsTable          = "channels"
	channelSettingsTable   = "channel_settings"
	channelCategoriesTable = "channel_categories"
	categoriesTable        = "categories"
	freelanceTasksTable    = "freelance_tasks"
	lastParsedTasksTable   = "last_parsed_tasks"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

var postgresDB *sqlx.DB

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	if postgresDB != nil {
		return postgresDB, nil
	}

	var err error
	postgresDB, err = sqlx.Open("pgx", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	if err := postgresDB.Ping(); err != nil {
		return nil, err
	}

	return postgresDB, nil
}
