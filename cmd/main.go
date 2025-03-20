package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
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

	// Parse secrets from env
	serializableSecrets := make(map[string]string)
	buildSecrets := make(map[string]string)
	if cfg.ServiceBuildSecrets != "" {
		if err := json.Unmarshal([]byte(cfg.ServiceBuildSecrets), &serializableSecrets); err != nil {
			log.Fatalf("Failed to parse secrets: %v", err)
		}

		// Convert back to map[string][]byte
		for k, v := range serializableSecrets {
			data, err := base64.StdEncoding.DecodeString(v)
			if err != nil {
				fmt.Printf("Error decoding secret %s: %v\n", k, err)
				continue
			}
			buildSecrets[k] = string(data)
		}
	}

	// Build with context
	dockerImg, repoName, err := builder.BuildWithRailpack(buildSecrets)
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
