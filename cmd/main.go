package main

import (
	"github.com/max-sanch/BotFreelancer-core"
	"github.com/max-sanch/BotFreelancer-core/pkg/handler"
	"github.com/max-sanch/BotFreelancer-core/pkg/repository"
	"github.com/max-sanch/BotFreelancer-core/pkg/service"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing config: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	repos := repository.NewPostgresRepos()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(core.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
