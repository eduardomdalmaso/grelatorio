package main

import (
	"context"

	"relatorio/internal/adapters/secondary/github"
	"relatorio/internal/adapters/secondary/storage"
	"relatorio/internal/core/domain"
	"relatorio/internal/core/usecases"
)

// App is the primary Wails adapter and application facade
type App struct {
	ctx             context.Context
	configUseCase   *usecases.ConfigUseCase
	activityUseCase *usecases.ActivityUseCase
}

// NewApp creates a new App application struct wiring Hexagonal dependencies
func NewApp() *App {
	configRepo := storage.NewFileConfigRepository()
	githubService := github.NewGithubAPIClient()

	return &App{
		configUseCase:   usecases.NewConfigUseCase(configRepo, githubService),
		activityUseCase: usecases.NewActivityUseCase(configRepo, githubService),
	}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// LoadConfig retrieves saved configuration
func (a *App) LoadConfig() (domain.Config, error) {
	return a.configUseCase.Load()
}

// SaveConfig validates and saves configuration
func (a *App) SaveConfig(cfg domain.Config) error {
	return a.configUseCase.Save(cfg)
}

// GetTokenExpiration retrieves token expiration string from GitHub
func (a *App) GetTokenExpiration(token string) (string, error) {
	return a.configUseCase.GetTokenExpiration(token)
}

// FetchUserRepositories retrieves repositories accessible to the configured user
func (a *App) FetchUserRepositories() ([]domain.Repository, error) {
	return a.activityUseCase.FetchUserRepositories()
}

// FetchGithubActivity retrieves commits and pull requests for given date range and optional repo filter
func (a *App) FetchGithubActivity(startDate, endDate, repoFilter string) ([]domain.GithubActivity, error) {
	return a.activityUseCase.FetchActivities(startDate, endDate, repoFilter)
}
