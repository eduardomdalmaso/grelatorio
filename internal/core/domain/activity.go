package domain

// GithubActivity represents a tracked commit or pull request item
type GithubActivity struct {
	Type        string `json:"type"`        // "commit" or "pull_request"
	Repo        string `json:"repo"`        // "owner/repo"
	Title       string `json:"title"`       // commit message summary or PR title
	Description string `json:"description"` // commit body or PR body
	URL         string `json:"url"`         // link on GitHub
	Date        string `json:"date"`        // ISO timestamp
	Ref         string `json:"ref"`         // Commit SHA or PR number (e.g. #12)
}

// ActivityFilter holds search parameters for fetching developer activities
type ActivityFilter struct {
	Username  string
	StartDate string
	EndDate   string
	Repo      string // Optional: filter by specific "owner/repo" or empty for all
}
