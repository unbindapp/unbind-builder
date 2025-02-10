package github

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v69/github"
	"github.com/unbindapp/unbind-builder/config"
)

type GitHubHelper struct {
	config *config.Config
	client *github.Client
	token  *github.InstallationToken
}

func NewGithubClient(config *config.Config) (*GitHubHelper, error) {
	transport, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, config.GithubAppID, config.GithubInstallationID, config.GithubPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub installation transport: %v", err)
	}
	client := github.NewClient(&http.Client{Transport: transport})
	token, _, err := client.Apps.CreateInstallationToken(context.Background(), config.GithubInstallationID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create installation token: %v", err)
	}
	return &GitHubHelper{
		config: config,
		client: github.NewClient(&http.Client{Transport: transport, Timeout: 5}),
		token:  token,
	}, nil
}

func (gh *GitHubHelper) GetRepositoryCloneURL() (string, error) {
	u, err := url.Parse(gh.config.GitRepoURL)
	if err != nil {
		return "", err
	}
	u.User = url.UserPassword("x-access-token", gh.token.GetToken())
	return u.String(), nil
}
