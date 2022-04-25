package service

import (
	"bytes"
	"encoding/json"
	core "github.com/max-sanch/BotFreelancer-core"
	"github.com/max-sanch/BotFreelancer-core/pkg/repository"
	"github.com/spf13/viper"
	"net/http"
)

type ChannelService struct {
	repo *repository.Repository
}

func NewChannelService(repo *repository.Repository) *ChannelService {
	return &ChannelService{repo: repo}
}

func (s *ChannelService) GetTasks() ([]core.ChannelTaskResponse, error) {
	var emptyTasks []core.ChannelTaskResponse

	lastParseTime, err := s.repo.Task.GetLastParseTime()
	if err != nil {
		return nil, err
	}

	parseTasks, err := getParseTasks(lastParseTime)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Task.SetLastParseTime(); err != nil {
		return nil, err
	}

	if len(parseTasks.Tasks) == 0 {
		return emptyTasks, nil
	}

	if err := s.repo.Task.DeleteAll(); err != nil {
		return nil, err
	}

	if err := s.repo.Task.AddTasks(parseTasks); err != nil {
		return nil, err
	}

	tasks, err := s.repo.Task.GetAllForChannels()
	if err != nil {
		return nil, err
	}

	return tasks, nil
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

func getParseTasks(datetime string) (core.TasksInput, error) {
	var tasks, emptyTasks core.TasksInput

	jsonRequest, err := json.Marshal(map[string]string{
		"datetime": datetime,
	})

	resp, err := http.Post(viper.GetString("url_parse_tasks"), "application/json", bytes.NewBuffer(jsonRequest))
	if err != nil {
		return emptyTasks, err
	}

	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		return emptyTasks, err
	}

	return tasks, nil
}
