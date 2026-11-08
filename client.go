package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/google/go-github/v92/github"
)

// NewGitHubClient creates GitHub client.
//
// When GITHUB_TOKEN is present in environment variable, the function will use it for authentication.
// When GitHub CLI is present in path, the function will get token from it.
//
// When both are present, environment variable will take precedence.
func NewGitHubClient() (*github.Client, error) {
	envTok := os.Getenv("GITHUB_TOKEN")
	ghTok, err := exec.Command("gh", "auth", "token").Output()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		fmt.Fprintf(os.Stderr, "failed to execute 'gh auth token': %v", err)
	}

	tok := strings.TrimSpace(string(ghTok))
	if envTok != "" {
		tok = envTok
	}

	var opts []github.ClientOptionsFunc
	if tok != "" {
		opts = append(opts, github.WithAuthToken(tok))
	}
	client, err := github.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create github client: %w", err)
	}
	return client, nil
}
