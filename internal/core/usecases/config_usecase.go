package usecases

import (
	"fmt"
	"relatorio/internal/core/domain"
	"relatorio/internal/core/ports"
)

// ConfigUseCase coordinates business logic related to app configuration
type ConfigUseCase struct {
	configRepo ports.ConfigRepository
	githubSvc  ports.GithubService
}

// NewConfigUseCase constructs a new ConfigUseCase
func NewConfigUseCase(configRepo ports.ConfigRepository, githubSvc ports.GithubService) *ConfigUseCase {
	return &ConfigUseCase{
		configRepo: configRepo,
		githubSvc:  githubSvc,
	}
}

// Load retrieves saved configuration
func (uc *ConfigUseCase) Load() (domain.Config, error) {
	return uc.configRepo.Load()
}

// Save validates and persists configuration, automatically resolving token expiration if needed
func (uc *ConfigUseCase) Save(cfg domain.Config) error {
	if err := cfg.Validate(); err != nil {
		return err
	}

	// Validate token with GitHub service
	exp, err := uc.githubSvc.GetTokenExpiration(cfg.Token)
	if err != nil {
		return fmt.Errorf("token inválido: %w", err)
	}
	cfg.TokenExpiration = exp

	return uc.configRepo.Save(cfg)
}

// GetTokenExpiration retrieves token expiration directly
func (uc *ConfigUseCase) GetTokenExpiration(token string) (string, error) {
	return uc.githubSvc.GetTokenExpiration(token)
}
