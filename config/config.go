package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/unbindapp/unbind-builder/internal/log"
)

type Config struct {
	GithubAppID int64 `env:"GITHUB_APP_ID,required"`
	// Installation ID of the app
	GithubInstallationID int64 `env:"GITHUB_INSTALLATION_ID,required"`
	// Repository to clone (github, https)
	GitRepoURL string `env:"GITHUB_REPO_URL,required"`
	// Branch to checkout and build
	GitBranch string `env:"GIT_BRANCH,required" envDefault:"master"`
	// GITHUB private key, needed to authenticate
	GithubPrivateKey string `env:"GITHUB_PRIVATE_KEY_PATH,required" envDefault:"/etc/github/private-key"`
}

// Parse environment variables into a Config struct
func NewConfig() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("Error parsing environment", "err", err)
	}
	return &cfg
}
