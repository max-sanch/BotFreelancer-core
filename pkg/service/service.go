package service

import (
	"github.com/max-sanch/BotFreelancer-core/pkg/repository"
)

type Authentication interface {

}

type User interface {

}

type Service struct {

}

func NewService(repos *repository.Repository) *Service {
	return &Service{}
}