package builder

import (
	"fmt"
	"os"
	"time"

	"github.com/railwayapp/railpack/buildkit"
	"github.com/railwayapp/railpack/core"
	a "github.com/railwayapp/railpack/core/app"
	"github.com/unbindapp/unbind-builder/internal/github"
	"github.com/unbindapp/unbind-builder/internal/log"
	"github.com/unbindapp/unbind-builder/internal/utils"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

func (self *Builder) BuildWithRailpack(buildSecrets map[string]string) (imageName, repoName string, err error) {
	// -- Generate image name
	repoName, err = utils.ExtractRepoName(self.config.GitRepoURL)
	if err != nil {
		log.Warnf("Failed to extract repository name: %v", err)
		repoName = fmt.Sprintf("unbind-build-%d", time.Now().Unix())
	}
	outputImage := fmt.Sprintf("%s/%s:%d", self.config.ContainerRegistryHost, repoName, time.Now().Unix())

	// -- Generate github tokens/create client
	ghHelper, err := github.NewGithubClient(self.config)
	if err != nil {
		return "", repoName, fmt.Errorf("failed to create GitHub client: %v", err)
	}

	// -- Clone repository
	// Create a temporary directory for cloning the repository.
	tmpDir, err := os.MkdirTemp("", "railpacks-build-")
	if err != nil {
		return "", repoName, fmt.Errorf("failed to create temporary directory: %v", err)
	}
	log.Infof("Created temporary directory: %s", tmpDir)
	// Clean up the temporary directory when done.
	defer os.RemoveAll(tmpDir)

	log.Infof("Cloning ref '%s' from '%s'", self.config.GitRef, self.config.GitRepoURL)

	_, err = git.PlainClone(tmpDir, false, &git.CloneOptions{
		URL: self.config.GitRepoURL,
		Auth: &http.BasicAuth{
			Username: "x-access-token",
			Password: ghHelper.InstallationToken.GetToken(),
		},
		ReferenceName: plumbing.ReferenceName(self.config.GitRef),
	})

	if err != nil {
		return "", repoName, fmt.Errorf("failed to clone repository: %v", err)
	}

	// --- Railpack build
	buildResult, app, _, err := GenerateBuildResult(tmpDir, buildSecrets)
	if err != nil {
		return "", repoName, fmt.Errorf("failed to generate build result: %v", err)
	}

	core.PrettyPrintBuildResult(buildResult, core.PrintOptions{Version: "unbind-builder"})

	if !buildResult.Success {
		return "", repoName, fmt.Errorf("build failed")
	}

	var platform buildkit.BuildPlatform
	platform = buildkit.DetermineBuildPlatformFromHost()

	err = buildkit.BuildWithBuildkitClient(app.Source, buildResult.Plan, buildkit.BuildWithBuildkitClientOptions{
		ImageName:    outputImage,
		DumpLLB:      false,
		OutputDir:    "",
		ProgressMode: "tty",
		CacheKey:     repoName,
		SecretsHash:  "",
		Secrets:      buildSecrets,
		Platform:     platform,
	})
	if err != nil {
		return "", repoName, fmt.Errorf("build failed: %v", err)
	}

	log.Infof("Built image %s", outputImage)
	return outputImage, repoName, nil
}

func GenerateBuildResult(directory string, buildSecrets map[string]string) (*core.BuildResult, *a.App, *a.Environment, error) {
	app, err := a.NewApp(directory)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error creating app: %w", err)
	}

	log.Infof("Building %s", app.Source)

	env := a.Environment{
		Variables: buildSecrets,
	}

	generateOptions := &core.GenerateBuildPlanOptions{
		RailpackVersion: "unbind-builder", // ! Add a version
	}

	buildResult := core.GenerateBuildPlan(app, &env, generateOptions)

	return buildResult, app, &env, nil
}
