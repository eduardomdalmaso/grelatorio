package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"relatorio/internal/core/domain"
)

// GithubAPIClient handles external communication with GitHub REST API
type GithubAPIClient struct {
	httpClient *http.Client
}

// NewGithubAPIClient constructs a new GithubAPIClient
func NewGithubAPIClient() *GithubAPIClient {
	return &GithubAPIClient{
		httpClient: &http.Client{Timeout: 20 * time.Second},
	}
}

func (c *GithubAPIClient) get(url, token, accept string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "github-report-generator-wails")
	if accept != "" {
		req.Header.Set("Accept", accept)
	} else {
		req.Header.Set("Accept", "application/vnd.github.v3+json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("github api returned status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// GetTokenExpiration checks the user token and retrieves expiration header
func (c *GithubAPIClient) GetTokenExpiration(token string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "github-report-generator-wails")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token inválido (status: %d)", resp.StatusCode)
	}

	return resp.Header.Get("GitHub-Authentication-Token-Expiration"), nil
}

type ghRepoItem struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Private     bool   `json:"private"`
	Description string `json:"description"`
	HTMLURL     string `json:"html_url"`
}

// FetchUserRepositories fetches repositories accessible to the user
func (c *GithubAPIClient) FetchUserRepositories(token string) ([]domain.Repository, error) {
	url := "https://api.github.com/user/repos?per_page=100&sort=updated&affiliation=owner,collaborator,organization_member"
	data, err := c.get(url, token, "")
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar repositórios: %w", err)
	}

	var items []ghRepoItem
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("erro ao decodificar repositórios: %w", err)
	}

	repos := make([]domain.Repository, 0, len(items))
	for _, item := range items {
		repos = append(repos, domain.Repository{
			ID:          item.ID,
			Name:        item.Name,
			FullName:    item.FullName,
			Private:     item.Private,
			Description: item.Description,
			HTMLURL:     item.HTMLURL,
		})
	}
	return repos, nil
}

type commitSearchItem struct {
	SHA    string `json:"sha"`
	Commit struct {
		Message string `json:"message"`
		Committer struct {
			Date string `json:"date"`
		} `json:"committer"`
	} `json:"commit"`
	HTMLURL    string `json:"html_url"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}

type commitSearchResponse struct {
	Items []commitSearchItem `json:"items"`
}

type issueSearchItem struct {
	Number        int    `json:"number"`
	Title         string `json:"title"`
	Body          string `json:"body"`
	HTMLURL       string `json:"html_url"`
	CreatedAt     string `json:"created_at"`
	RepositoryURL string `json:"repository_url"`
}

type issueSearchResponse struct {
	Items []issueSearchItem `json:"items"`
}

// FetchActivities retrieves commits and PRs with optional repo filtering
func (c *GithubAPIClient) FetchActivities(token string, filter domain.ActivityFilter) ([]domain.GithubActivity, error) {
	var activities []domain.GithubActivity

	repoQuery := ""
	if filter.Repo != "" && filter.Repo != "all" {
		repoQuery = "+repo:" + strings.TrimSpace(filter.Repo)
	}

	// 1. Fetch Commits
	commitUrl := fmt.Sprintf("https://api.github.com/search/commits?q=author:%s+committer-date:%s..%s%s&per_page=100",
		filter.Username, filter.StartDate, filter.EndDate, repoQuery)

	commitData, err := c.get(commitUrl, token, "application/vnd.github.cloak-preview+json")
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar commits: %w", err)
	}

	var commitResp commitSearchResponse
	if err := json.Unmarshal(commitData, &commitResp); err == nil {
		for _, item := range commitResp.Items {
			parts := strings.SplitN(item.Commit.Message, "\n", 2)
			title := parts[0]
			desc := ""
			if len(parts) > 1 {
				desc = strings.TrimSpace(parts[1])
			}

			ref := item.SHA
			if len(ref) > 7 {
				ref = ref[:7]
			}

			activities = append(activities, domain.GithubActivity{
				Type:        "commit",
				Repo:        item.Repository.FullName,
				Title:       title,
				Description: desc,
				URL:         item.HTMLURL,
				Date:        item.Commit.Committer.Date,
				Ref:         ref,
			})
		}
	}

	// 2. Fetch Pull Requests
	prUrl := fmt.Sprintf("https://api.github.com/search/issues?q=author:%s+type:pr+created:%s..%s%s&per_page=100",
		filter.Username, filter.StartDate, filter.EndDate, repoQuery)

	prData, err := c.get(prUrl, token, "")
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar pull requests: %w", err)
	}

	var prResp issueSearchResponse
	if err := json.Unmarshal(prData, &prResp); err == nil {
		for _, item := range prResp.Items {
			repoName := ""
			parts := strings.Split(item.RepositoryURL, "/repos/")
			if len(parts) > 1 {
				repoName = parts[1]
			}

			activities = append(activities, domain.GithubActivity{
				Type:        "pull_request",
				Repo:        repoName,
				Title:       item.Title,
				Description: item.Body,
				URL:         item.HTMLURL,
				Date:        item.CreatedAt,
				Ref:         fmt.Sprintf("#%d", item.Number),
			})
		}
	}

	return activities, nil
}
