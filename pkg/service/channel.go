package service

import (
	core "github.com/max-sanch/BotFreelancer-core"
	"github.com/max-sanch/BotFreelancer-core/pkg/repository"
)

type ChannelService struct {
	repo repository.Channel
}

func NewChannelService(repo repository.Channel) *ChannelService {
	return &ChannelService{repo: repo}
}

func (s *ChannelService) GetChannel(apiID int) (core.ChannelResponse, error) {
	return s.repo.GetChannel(apiID)
}

func (s *ChannelService) CreateChannel(channelInput core.ChannelInput) (int, error) {
	return s.repo.CreateChannel(channelInput)
}

func (s *ChannelService) UpdateChannel(channelInput core.ChannelInput) (int, error) {
	return s.repo.UpdateChannel(channelInput)
}

func (s *ChannelService) DeleteChannel(apiID int) error {
	return s.repo.DeleteChannel(apiID)
}
