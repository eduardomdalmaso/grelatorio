package domain

import (
	"errors"
	"strings"
)

// Config represents application configuration and client parameters
type Config struct {
	Token           string `json:"token"`
	Username        string `json:"username"`
	Client          string `json:"client"`
	Rate            string `json:"rate"`
	TokenExpiration string `json:"token_expiration"`
	CNPJ            string `json:"cnpj"`
	CompanyName     string `json:"company_name"`
}

// Validate checks if essential credentials are present
func (c *Config) Validate() error {
	if strings.TrimSpace(c.Token) == "" {
		return errors.New("token do GitHub é obrigatório")
	}
	if strings.TrimSpace(c.Username) == "" {
		return errors.New("nome de usuário do GitHub é obrigatório")
	}
	return nil
}
