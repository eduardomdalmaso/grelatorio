package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

type Config struct {
	Token           string `json:"token"`
	Username        string `json:"username"`
	Client          string `json:"client"`
	Rate            string `json:"rate"`
	TokenExpiration string `json:"token_expiration"`
	CNPJ            string `json:"cnpj"`
	CompanyName     string `json:"company_name"`
}

type GithubActivity struct {
	Type        string `json:"type"`        // "commit" or "pull_request"
	Repo        string `json:"repo"`        // "owner/repo"
	Title       string `json:"title"`       // commit message or PR title
	Description string `json:"description"` // body of PR, commit description, or empty
	URL         string `json:"url"`
	Date        string `json:"date"`        // ISO timestamp
	Ref         string `json:"ref"`         // Commit SHA or PR number
}

// Github API Structs
type CommitSearchResponse struct {
	Items []CommitSearchItem `json:"items"`
}

type CommitSearchItem struct {
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

type IssueSearchResponse struct {
	Items []IssueSearchItem `json:"items"`
}

type IssueSearchItem struct {
	Number        int    `json:"number"`
	Title         string `json:"title"`
	Body          string `json:"body"`
	HTMLURL       string `json:"html_url"`
	CreatedAt     string `json:"created_at"`
	RepositoryURL string `json:"repository_url"`
}

func getConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(dir, "github-report-generator")
	err = os.MkdirAll(appDir, 0755)
	if err != nil {
		return "", err
	}
	return filepath.Join(appDir, "config.json"), nil
}

func (a *App) LoadConfig() (Config, error) {
	path, err := getConfigPath()
	if err != nil {
		return Config{}, err
	}
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, nil // Return empty config if file doesn't exist
		}
		return Config{}, err
	}
	defer file.Close()

	var cfg Config
	err = json.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (a *App) SaveConfig(cfg Config) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(cfg)
}

func (a *App) GetTokenExpiration(token string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "github-report-generator-wails")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token inválido (status: %d)", resp.StatusCode)
	}

	expiration := resp.Header.Get("GitHub-Authentication-Token-Expiration")
	return expiration, nil
}

func githubGet(url string, token string, accept string) ([]byte, error) {
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

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github api returned status %d: %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

func (a *App) FetchGithubActivity(startDate string, endDate string) ([]GithubActivity, error) {
	cfg, err := a.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	if cfg.Token == "" || cfg.Username == "" {
		return nil, fmt.Errorf("GitHub credentials not configured. Please save your Token and Username in Settings.")
	}

	activities := []GithubActivity{}

	// 1. Fetch Commits
	commitUrl := fmt.Sprintf("https://api.github.com/search/commits?q=author:%s+committer-date:%s..%s&per_page=100",
		cfg.Username, startDate, endDate)

	commitData, err := githubGet(commitUrl, cfg.Token, "application/vnd.github.cloak-preview+json")
	if err == nil {
		var commitResp CommitSearchResponse
		if err := json.Unmarshal(commitData, &commitResp); err == nil {
			for _, item := range commitResp.Items {
				parts := strings.SplitN(item.Commit.Message, "\n", 2)
				title := parts[0]
				desc := ""
				if len(parts) > 1 {
					desc = strings.TrimSpace(parts[1])
				}

				activities = append(activities, GithubActivity{
					Type:        "commit",
					Repo:        item.Repository.FullName,
					Title:       title,
					Description: desc,
					URL:         item.HTMLURL,
					Date:        item.Commit.Committer.Date,
					Ref:         item.SHA[:7],
				})
			}
		}
	} else {
		return nil, fmt.Errorf("failed to fetch commits: %w", err)
	}

	// 2. Fetch Pull Requests
	prUrl := fmt.Sprintf("https://api.github.com/search/issues?q=author:%s+type:pr+created:%s..%s&per_page=100",
		cfg.Username, startDate, endDate)

	prData, err := githubGet(prUrl, cfg.Token, "")
	if err == nil {
		var prResp IssueSearchResponse
		if err := json.Unmarshal(prData, &prResp); err == nil {
			for _, item := range prResp.Items {
				repoName := ""
				parts := strings.Split(item.RepositoryURL, "/repos/")
				if len(parts) > 1 {
					repoName = parts[1]
				}

				activities = append(activities, GithubActivity{
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
	} else {
		return nil, fmt.Errorf("failed to fetch pull requests: %w", err)
	}

	return activities, nil
}
