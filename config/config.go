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
	GitRef string `env:"GIT_REF,required"`
	// GITHUB private key, needed to authenticate
	GithubPrivateKey string `env:"GITHUB_PRIVATE_KEY_PATH,required" envDefault:"/etc/github/private-key"`
	// Registry specific
	ContainerRegistryHost     string `env:"CONTAINER_REGISTRY_HOST,required" envDefault:"docker-registry.unbind-system:5000"`
	ContainerRegistryUser     string `env:"CONTAINER_REGISTRY_USER,required" envDefault:"admin"`
	ContainerRegistryPassword string `env:"CONTAINER_REGISTRY_PASSWORD,required"`
	// Docker host because nixpacks ignores the variable https://github.com/railwayapp/nixpacks/issues/1194
	DockerHost string `env:"DOCKER_HOST" envDefault:"unix:///var/run/docker.sock"`
	// Deployment namespace (kubernetes)
	DeploymentNamespace string `env:"DEPLOYMENT_NAMESPACE" envDefault:"unbind-user"`
}

// Parse environment variables into a Config struct
func NewConfig() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("Error parsing environment", "err", err)
	}
	return &cfg
}
