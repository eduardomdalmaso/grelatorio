package usecases

import (
	"errors"
	"fmt"
	"relatorio/internal/core/domain"
	"relatorio/internal/core/ports"
)

// ActivityUseCase coordinates fetching developer activities and repositories
type ActivityUseCase struct {
	configRepo ports.ConfigRepository
	githubSvc  ports.GithubService
}

// NewActivityUseCase constructs a new ActivityUseCase
func NewActivityUseCase(configRepo ports.ConfigRepository, githubSvc ports.GithubService) *ActivityUseCase {
	return &ActivityUseCase{
		configRepo: configRepo,
		githubSvc:  githubSvc,
	}
}

// FetchActivities fetches commits and pull requests for given dates and optional repository filter
func (uc *ActivityUseCase) FetchActivities(startDate, endDate, repoFilter string) ([]domain.GithubActivity, error) {
	cfg, err := uc.configRepo.Load()
	if err != nil {
		return nil, fmt.Errorf("falha ao carregar configurações: %w", err)
	}

	if cfg.Token == "" || cfg.Username == "" {
		return nil, errors.New("credenciais do GitHub não configuradas. Salve seu Token e Username nas Configurações")
	}

	filter := domain.ActivityFilter{
		Username:  cfg.Username,
		StartDate: startDate,
		EndDate:   endDate,
		Repo:      repoFilter,
	}

	return uc.githubSvc.FetchActivities(cfg.Token, filter)
}

// FetchUserRepositories returns the list of repositories accessible to the user
func (uc *ActivityUseCase) FetchUserRepositories() ([]domain.Repository, error) {
	cfg, err := uc.configRepo.Load()
	if err != nil {
		return nil, fmt.Errorf("falha ao carregar configurações: %w", err)
	}

	if cfg.Token == "" {
		return nil, errors.New("token do GitHub não configurado")
	}

	return uc.githubSvc.FetchUserRepositories(cfg.Token)
}
