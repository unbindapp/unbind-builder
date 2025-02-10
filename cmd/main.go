package main

import (
	"github.com/joho/godotenv"
	"github.com/unbindapp/unbind-builder/config"
	"github.com/unbindapp/unbind-builder/internal/builder"
	"github.com/unbindapp/unbind-builder/internal/log"
	"github.com/unbindapp/unbind-builder/internal/registry"
)

func main() {
	godotenv.Load()

	builder := builder.NewBuilder(
		config.NewConfig(),
	)

	dockerImg, err := builder.BuildWithNixpacks()
	if err != nil {
		log.Fatalf("Failed to build with nixpacks: %v", err)
	}

	// Push to registry
	if err := registry.PushImageToRegistry(dockerImg); err != nil {
		log.Fatalf("Failed to push to registry: %v", err)
	}
}
