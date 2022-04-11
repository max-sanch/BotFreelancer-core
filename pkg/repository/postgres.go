package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable = "users"
	userSettingsTable = "user_settings"
	userCategoriesTable = "user_categories"
	channelsTable = "channels"
	channelSettingsTable = "channel_settings"
	channelCategoriesTable = "channel_categories"
	categoriesTable = "categories"
	freelanceDataTable = "freelance_data"
	freelanceSectionsTable = "freelance_sections"
)

type Config struct {
	Host string
	Port string
	Username string
	Password string
	DBName string
	SSLMode string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
