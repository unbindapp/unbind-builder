package builder

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/unbindapp/unbind-builder/internal/github"
	"github.com/unbindapp/unbind-builder/internal/log"
	"github.com/unbindapp/unbind-builder/internal/utils"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

func (b *Builder) BuildWithNixpacks() (imageName, repoName string, err error) {
	// -- Generate image name
	repoName, err = utils.ExtractRepoName(b.config.GitRepoURL)
	if err != nil {
		log.Warnf("Failed to extract repository name: %v", err)
		repoName = fmt.Sprintf("unbind-build-%d", time.Now().Unix())
	}
	outputImage := fmt.Sprintf("%s/%s:%d", b.config.ContainerRegistryHost, repoName, time.Now().Unix())

	// -- Generate github tokens/create client
	ghHelper, err := github.NewGithubClient(b.config)
	if err != nil {
		return "", repoName, fmt.Errorf("failed to create GitHub client: %v", err)
	}

	// -- Clone repository
	// Create a temporary directory for cloning the repository.
	tmpDir, err := os.MkdirTemp("", "nixpacks-build-")
	if err != nil {
		return "", repoName, fmt.Errorf("failed to create temporary directory: %v", err)
	}
	log.Infof("Created temporary directory: %s", tmpDir)
	// Clean up the temporary directory when done.
	defer os.RemoveAll(tmpDir)

	log.Infof("Cloning ref '%s' from '%s'", b.config.GitRef, b.config.GitRepoURL)

	_, err = git.PlainClone(tmpDir, false, &git.CloneOptions{
		URL: b.config.GitRepoURL,
		Auth: &http.BasicAuth{
			Username: "x-access-token",
			Password: ghHelper.InstallationToken.GetToken(),
		},
		ReferenceName: plumbing.ReferenceName(b.config.GitRef),
	})

	if err != nil {
		return "", repoName, fmt.Errorf("failed to clone repository: %v", err)
	}

	// --- Nixpacks build
	log.Infof("Running nixpacks build in directory: %s", tmpDir)
	buildCmd := exec.Command("nixpacks", "build", tmpDir, "--tag", outputImage, "--docker-host", b.config.DockerHost)
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		log.Fatalf("nixpacks build failed: %v", err)
	}

	log.Infof("Built image %s", outputImage)
	return outputImage, repoName, nil
}
