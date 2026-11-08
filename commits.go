package main

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/google/go-github/v92/github"
)

type CommitsCmd struct {
	Repo string `arg:"" help:"format: <owner>/<repo>"`
}

func (c *CommitsCmd) Run(ctx *Context) error {
	sep := strings.Split(c.Repo, "/")
	if len(sep) != 2 {
		return fmt.Errorf("repo name should be like <owner>/<repo>")
	}
	owner := sep[0]
	repoName := sep[1]

	client, err := NewGitHubClient()
	if err != nil {
		return fmt.Errorf("failed to create github client: %w", err)
	}

	it := client.Repositories.ListCommitsIter(context.TODO(), owner, repoName, &github.CommitsListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	})
	commits := make([]*github.RepositoryCommit, 0, 10000)
	for c, err := range it {
		if err != nil {
			return fmt.Errorf("error while listing commits: %w", err)
		}
		commits = append(commits, c)
	}

	fmt.Println(len(commits))

	res := AnalyzeCommits(commits)
	sorted := res.Sorted()
	sorted.Print()

	return nil
}

type CommtitInfo struct {
	typ string
	// will be empty when no scope specified.
	scope string
}

type AnalysisResult struct {
	commitsPerType  map[string]uint
	commitsPerScope map[string]uint
}

type Entry struct {
	Key   string
	Value uint
}

type SortedAnalysisResult struct {
	commitsPerType  []Entry
	commitsPerScope []Entry
}

// NewAnalysisResult creates new [AnalysisResult].
func NewAnalysisResult() AnalysisResult {
	return AnalysisResult{
		commitsPerType:  make(map[string]uint),
		commitsPerScope: make(map[string]uint),
	}
}

func AnalyzeCommits(commits []*github.RepositoryCommit) AnalysisResult {
	res := NewAnalysisResult()
	for _, c := range commits {
		info, err := ParseCommitInfo(*c.Commit.Message)
		if err != nil {
			continue
		}
		if _, ok := res.commitsPerType[info.typ]; !ok {
			res.commitsPerType[info.typ] = 1
		} else {
			res.commitsPerType[info.typ] += 1
		}
		if _, ok := res.commitsPerScope[info.scope]; !ok {
			res.commitsPerScope[info.scope] = 1
		} else {
			res.commitsPerScope[info.scope] += 1
		}
	}
	return res
}

// Sorted sorts [AnalysisResult] and returns it.
func (sar *AnalysisResult) Sorted() SortedAnalysisResult {
	commitsPerType := make([]Entry, 0, len(sar.commitsPerType))
	for k, v := range sar.commitsPerType {
		commitsPerType = append(commitsPerType, Entry{
			Key:   k,
			Value: v,
		})
	}
	slices.SortFunc(commitsPerType, func(a, b Entry) int {
		return -cmp.Compare(a.Value, b.Value)
	})

	commitsPerScope := make([]Entry, 0, len(sar.commitsPerScope))
	for k, v := range sar.commitsPerScope {
		commitsPerScope = append(commitsPerScope, Entry{
			Key:   k,
			Value: v,
		})
	}
	slices.SortFunc(commitsPerScope, func(a, b Entry) int {
		return -cmp.Compare(a.Value, b.Value)
	})

	return SortedAnalysisResult{
		commitsPerType:  commitsPerType,
		commitsPerScope: commitsPerScope,
	}
}

// Print prints the result to stdout.
func (sar *SortedAnalysisResult) Print() {
	fmt.Println("Number of commits per scope:")
	for _, v := range sar.commitsPerScope {
		fmt.Printf("%v: %v\n", v.Key, v.Value)
	}

	fmt.Println("\n\nNumber of commits per type:")
	for _, v := range sar.commitsPerType {
		fmt.Printf("%v: %v\n", v.Key, v.Value)
	}
}

// ParseCommitInfo parses message following Conventional Commits.
func ParseCommitInfo(message string) (*CommtitInfo, error) {
	msg, _, found := strings.Cut(message, "\n")
	if !found {
		msg = message
	}

	typeScope, _, found := strings.Cut(msg, ":")
	if !found {
		return nil, errors.New("not a valid conventional commit: there's no colon")
	}

	typ, scope, found := strings.Cut(typeScope, "(")
	if found {
		if scope[len(scope)-1] != ')' {
			return nil, errors.New("not a valid conventional commit: there's no matching right parenthes after scope")
		}
		scope = scope[:len(scope)-1]
	}
	return &CommtitInfo{
		typ:   typ,
		scope: scope,
	}, nil
}
