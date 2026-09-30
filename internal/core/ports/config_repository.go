package ports

import "relatorio/internal/core/domain"

// ConfigRepository defines output port for configuration storage
type ConfigRepository interface {
	Load() (domain.Config, error)
	Save(cfg domain.Config) error
}
