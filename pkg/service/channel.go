package service

import (
	core "github.com/max-sanch/BotFreelancer-core"
	"github.com/max-sanch/BotFreelancer-core/pkg/repository"
)

type ChannelService struct {
	repo *repository.Repository
}

func NewChannelService(repo *repository.Repository) *ChannelService {
	return &ChannelService{repo: repo}
}

func (s *ChannelService) GetByApiId(apiId int) (core.ChannelResponse, error) {
	return s.repo.Channel.GetByApiId(apiId)
}

func (s *ChannelService) Create(channelInput core.ChannelInput) (int, error) {
	return s.repo.Channel.Create(channelInput)
}

func (s *ChannelService) Update(channelInput core.ChannelInput) (int, error) {
	return s.repo.Channel.Update(channelInput)
}

func (s *ChannelService) Delete(apiId int) error {
	return s.repo.Channel.Delete(apiId)
}
