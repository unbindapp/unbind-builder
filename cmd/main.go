package main

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/unbindapp/unbind-builder/config"
	"github.com/unbindapp/unbind-builder/internal/builder"
	"github.com/unbindapp/unbind-builder/internal/kubernetes"
	"github.com/unbindapp/unbind-builder/internal/log"
	"github.com/unbindapp/unbind-builder/internal/registry"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := godotenv.Load(); err != nil {
		log.Warnf("Failed to load .env file: %v", err)
	}

	cfg := config.NewConfig()
	os.Setenv("BUILDKIT_HOST", cfg.BuildkitHost)

	builder := builder.NewBuilder(cfg)
	kubernetesUtil := kubernetes.NewKubernetesUtil(cfg)

	// Build with context
	dockerImg, repoName, err := builder.BuildWithRailpack()
	if err != nil {
		log.Fatalf("Failed to build with railpack: %v", err)
	}

	// Push to registry with context
	if err := registry.PushImageToRegistry(ctx, dockerImg, cfg); err != nil {
		log.Fatalf("Failed to push to registry: %v", err)
	}

	// Deploy to kubernetes with context
	createdCRD, err := kubernetesUtil.DeployImage(ctx, repoName, dockerImg)
	if err != nil {
		log.Fatalf("Failed to deploy image: %v", err)
	}

	log.Infof("Created CRD: %v", createdCRD)
}
