package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v92/github"
)

type LatestCmd struct {
	Repo string `arg:"" help:"format: <owner>/<repo>"`
}

func (l *LatestCmd) Run(ctx *Context) error {
	sep := strings.Split(l.Repo, "/")
	if len(sep) != 2 {
		return fmt.Errorf("repo name should be like <owner>/<repo>")
	}
	owner := sep[0]
	repoName := sep[1]

	client, err := github.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create github client: %w", err)
	}

	repo, _, err := client.Repositories.Get(context.TODO(), owner, repoName)
	if err != nil {
		return fmt.Errorf("failed to fetch repository info: %w", err)
	}
	defaultBranch := repo.DefaultBranch
	if defaultBranch == nil {
		tmp := "main"
		defaultBranch = &tmp
	}

	branch, _, err := client.Repositories.GetBranch(context.TODO(), owner, repoName, *defaultBranch, 8)
	if err != nil {
		return fmt.Errorf("failed to get branch info: %w", err)
	}
	latest := branch.Commit.SHA
	fmt.Println(*latest)

	return nil
}
