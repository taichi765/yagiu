package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type RepositoryInfo struct {
	ClonedAt time.Time `json:"clonedAt"`
}

type RepositoryInfoList struct {
	Repositories []RepositoryInfo `json:"repositories"`
}

type CloneCmd struct {
	rawOptions []string `arg:"" passthrough:"all"`
}

func (c *CloneCmd) Run(ctx *Context) error {
	/*path := RepositoryInfoFile(ctx.ConfigDir)
	old, err := readRepositoriesInfo(path)
	if err != os.ErrNotExist && err != nil {
		return err
	}*/

	// now := time.Now()

	args := []string{"clone"}
	for _, v := range c.rawOptions {
		args = append(args, v)
	}

	cmd := exec.Command("git", args...)
	/*stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get the pipe to stdout: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get pipe to stderr: %w", err)
		}*/

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to run 'git clone': %w", err)
	}

	return nil
}

func readRepositoriesInfo(path string) (*RepositoryInfoList, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var oldRepositories RepositoryInfoList
	err = json.NewDecoder(f).Decode(&oldRepositories)
	if err != nil {
		return nil, fmt.Errorf("failed to decode json: %w", err)
	}
	return &oldRepositories, nil
}
