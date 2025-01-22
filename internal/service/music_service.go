package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"music/internal/config"
	"music/internal/models"
	"music/internal/repository"
	"net/http"
	"time"
)

const songBaseUrl = "https://example.com"

type Music interface {
	GetMusics(ctx context.Context, group, song, releaseDate, text string, limit, offset int) ([]models.MusicInfo, error)
	GetSongLyricsByVerses(ctx context.Context, ID, couplet, size int) ([]string, error)
	AddMusic(ctx context.Context, music models.Music) error
	UpdateMusic(ctx context.Context, ID int, updates models.MusicUpdate) error
	DeleteMusic(ctx context.Context, ID int) error
}

type musicService struct {
	repos   repository.Music
	logger  *slog.Logger
	timeout time.Duration
}

func newMusicService(repos repository.Music, logger *slog.Logger) *musicService {
	return &musicService{
		repos:   repos,
		logger:  logger,
		timeout: 3 * time.Second,
	}
}

func (s *musicService) GetMusics(
	ctx context.Context, group, song, releaseDate, text string, limit, offset int,
) ([]models.MusicInfo, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return s.repos.GetMusics(c, group, song, releaseDate, text, limit, offset)
}

func (s *musicService) GetSongLyricsByVerses(ctx context.Context, ID, couplet, size int) ([]string, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return s.repos.GetSongLyricsByVerses(c, ID, couplet, size)
}

func (s *musicService) AddMusic(ctx context.Context, music models.Music) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	response, err := s.getMusic(music)
	if err != nil {
		return err
	}
	s.logger.DebugContext(ctx, "Got music response")

	return s.repos.AddMusic(c, response)
}

func (s *musicService) UpdateMusic(ctx context.Context, ID int, updates models.MusicUpdate) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return s.repos.UpdateMusic(c, ID, updates)
}

func (s *musicService) DeleteMusic(ctx context.Context, ID int) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	return s.repos.DeleteMusic(c, ID)
}

func (s *musicService) getMusic(music models.Music) (models.MusicInfo, error) {
	var result models.MusicInfo
	songUrl := config.Url + "/info"

	s.logger.Debug("URL to get music", slog.String("url", songUrl))

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, songUrl, nil)
	if err != nil {
		return result, err
	}

	query := req.URL.Query()
	query.Add("group", music.Group)
	query.Add("song", music.Song)

	req.URL.RawQuery = query.Encode()

	resp, err := client.Do(req)
	s.logger.Debug("sending request to get music")
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}
	s.logger.Debug("request status code: ", slog.String("status", resp.Status))

	var apiResponse struct {
		RelaseDate string `json:"relaseDate"`
		Text       string `json:"text"`
		Link       string `json:"link"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return result, err
	}

	result.Group = music.Group
	result.Song = music.Song
	result.RelaseDate = apiResponse.RelaseDate
	result.Text = apiResponse.Text
	result.Link = apiResponse.Link

	return result, nil
}
