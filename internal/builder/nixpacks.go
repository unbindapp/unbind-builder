package builder

import (
	"log"
	"os"
	"os/exec"

	"github.com/unbindapp/unbind-builder/internal/github"
)

func (b *Builder) BuildWithNixpacks() {
	ghHelper, err := github.NewGithubClient(b.config)
	if err != nil {
		log.Fatal("Failed to create GitHub client", "err", err)
	}

	cloneUrl, err := ghHelper.GetRepositoryCloneURL()
	if err != nil {
		log.Fatal("Failed to get repository clone URL", "err", err)
	}

	// Create a temporary directory for cloning the repository.
	tmpDir, err := os.MkdirTemp("", "nixpacks-build-")
	if err != nil {
		log.Fatalf("Failed to create temporary directory: %v", err)
	}
	log.Printf("Created temporary directory: %s", tmpDir)
	// Clean up the temporary directory when done.
	// defer os.RemoveAll(tmpDir)

	// Clone the repository into the temporary directory.
	log.Printf("Cloning repository %s", cloneUrl)
	cloneCmd := exec.Command("git", "clone", cloneUrl, tmpDir)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	if err := cloneCmd.Run(); err != nil {
		log.Fatalf("Failed to clone repository: %v", err)
	}

}
