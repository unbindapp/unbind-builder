package github

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/v69/github"
	"github.com/unbindapp/unbind-builder/config"
	"github.com/unbindapp/unbind-builder/internal/utils"
)

type GitHubHelper struct {
	config            *config.Config
	client            *github.Client
	InstallationToken *github.InstallationToken
}

func NewGithubClient(config *config.Config) (*GitHubHelper, error) {
	// Decode private key
	privateKey, err := utils.DecodePrivateKey(config.GithubAppPrivateKey)
	if err != nil {
		return nil, err
	}

	// Get token
	githubJwt, err := utils.GenerateJWT(config.GithubAppID, privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate github JWT: %v", err)
	}

	client := github.NewClient(&http.Client{}).WithAuthToken(githubJwt)
	token, _, err := client.Apps.CreateInstallationToken(context.Background(), config.GithubInstallationID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create installation token: %v", err)
	}
	return &GitHubHelper{
		config:            config,
		client:            client,
		InstallationToken: token,
	}, nil
}
