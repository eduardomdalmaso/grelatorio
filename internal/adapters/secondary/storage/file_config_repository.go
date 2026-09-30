package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"relatorio/internal/core/domain"
)

// FileConfigRepository persists configuration into a JSON file in user app data
type FileConfigRepository struct{}

// NewFileConfigRepository constructs a new FileConfigRepository
func NewFileConfigRepository() *FileConfigRepository {
	return &FileConfigRepository{}
}

func (r *FileConfigRepository) getFilePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(dir, "github-report-generator")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(appDir, "config.json"), nil
}

// Load reads and parses config.json
func (r *FileConfigRepository) Load() (domain.Config, error) {
	filePath, err := r.getFilePath()
	if err != nil {
		return domain.Config{}, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return domain.Config{}, nil
		}
		return domain.Config{}, err
	}
	defer file.Close()

	var cfg domain.Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return domain.Config{}, err
	}
	return cfg, nil
}

// Save writes config into config.json
func (r *FileConfigRepository) Save(cfg domain.Config) error {
	filePath, err := r.getFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(cfg)
}
