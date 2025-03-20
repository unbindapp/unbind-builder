package registry

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/unbindapp/unbind-builder/config"
	"github.com/unbindapp/unbind-builder/internal/log"
)

// Executes docker push
func PushImageToRegistry(ctx context.Context, image string, cfg *config.Config) error {
	// Extract registry host from the image name.
	// This assumes that the image is prefixed with the registry host.
	// e.g., "unbind-registry.unbind.app/my-image:tag"
	parts := strings.Split(image, "/")
	if len(parts) < 2 {
		return fmt.Errorf("image name %q does not contain a registry hostname", image)
	}
	registry := parts[0]

	// Execute "docker login" with the provided credentials.
	loginCmd := exec.Command("docker", "login", registry, "-u", cfg.ContainerRegistryUser, "-p", cfg.ContainerRegistryPassword)
	loginCmd.Stdout = os.Stdout
	loginCmd.Stderr = os.Stderr

	log.Infof("Logging in to registry %s...\n", registry)
	if err := loginCmd.Run(); err != nil {
		return fmt.Errorf("docker login failed: %w", err)
	}

	// Execute "docker push" for the image.
	pushCmd := exec.Command("docker", "push", image)
	pushCmd.Stdout = os.Stdout
	pushCmd.Stderr = os.Stderr

	fmt.Printf("Pushing image %s...\n", image)
	if err := pushCmd.Run(); err != nil {
		return fmt.Errorf("docker push failed: %w", err)
	}

	return nil
}
