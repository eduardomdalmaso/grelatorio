package ports

import "relatorio/internal/core/domain"

// GithubService defines output port for interacting with the GitHub API
type GithubService interface {
	GetTokenExpiration(token string) (string, error)
	FetchUserRepositories(token string) ([]domain.Repository, error)
	FetchActivities(token string, filter domain.ActivityFilter) ([]domain.GithubActivity, error)
}
