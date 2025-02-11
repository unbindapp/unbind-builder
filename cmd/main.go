package main

import (
	"github.com/joho/godotenv"
	"github.com/unbindapp/unbind-builder/config"
	"github.com/unbindapp/unbind-builder/internal/builder"
	"github.com/unbindapp/unbind-builder/internal/kubernetes"
	"github.com/unbindapp/unbind-builder/internal/log"
	"github.com/unbindapp/unbind-builder/internal/registry"
)

func main() {
	godotenv.Load()

	cfg := config.NewConfig()

	builder := builder.NewBuilder(
		cfg,
	)

	kubernetesUtil := kubernetes.NewKubernetesUtil(cfg)

	dockerImg, repoName, err := builder.BuildWithNixpacks()
	if err != nil {
		log.Fatalf("Failed to build with nixpacks: %v", err)
	}

	// Push to registry
	if err := registry.PushImageToRegistry(dockerImg, cfg); err != nil {
		log.Fatalf("Failed to push to registry: %v", err)
	}

	// Deploy to kubernetes
	createdCRD, err := kubernetesUtil.DeployImage(repoName, dockerImg)
	if err != nil {
		log.Fatalf("Failed to deploy image: %v", err)
	}
	log.Infof("Created CRD: %v", createdCRD)
}
