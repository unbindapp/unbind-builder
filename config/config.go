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
	// Github app private key
	GithubAppPrivateKey string `env:"GITHUB_APP_PRIVATE_KEY,required"`
	// Registry specific
	ContainerRegistryHost     string `env:"CONTAINER_REGISTRY_HOST,required" envDefault:"docker-registry.unbind-system:5000"`
	ContainerRegistryUser     string `env:"CONTAINER_REGISTRY_USER,required" envDefault:"admin"`
	ContainerRegistryPassword string `env:"CONTAINER_REGISTRY_PASSWORD,required"`
	// Docker host because nixpacks ignores the variable https://github.com/railwayapp/nixpacks/issues/1194
	BuildkitHost string `env:"BUILDKIT_HOST" envDefault:"docker-container://buildkit"`
	// Deployment namespace (kubernetes)
	DeploymentNamespace string `env:"DEPLOYMENT_NAMESPACE" envDefault:"unbind-user"`
	// Service specific
	ServiceRuntime    string  `env:"SERVICE_RUNTIME"`
	ServiceFramework  string  `env:"SERVICE_FRAMEWORK"`
	ServicePublic     *bool   `env:"SERVICE_PUBLIC"`
	ServicePort       *int32  `env:"SERVICE_PORT"`
	ServiceHost       *string `env:"SERVICE_HOST"`
	ServiceReplicas   *int32  `env:"SERVICE_REPLICAS"`
	ServiceSecretName string  `env:"SERVICE_SECRET_NAME,required"`
	// Kubeconfig for local testing
	KubeConfig string `env:"KUBECONFIG"`
}

// Parse environment variables into a Config struct
func NewConfig() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("Error parsing environment", "err", err)
	}
	return &cfg
}
