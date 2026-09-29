package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
)

type Context struct {
	Verbose   bool
	ConfigDir string
}

type ListCmd struct {
}

func (l *ListCmd) Run(ctx *Context) error {
	return nil
}

var cli struct {
	Verbose bool `heko:"Print verbose output"`

	Clone  CloneCmd  `cmd:"" help:"Clone a remote repository."`
	List   ListCmd   `cmd:"" help:"List git repositories."`
	Latest LatestCmd `cmd:"" help:"Get the latest commit of the specified repository."`
}

func main() {
	dir, err := ConfigDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get path of config dir: %v\n", err)
	}
	err = os.MkdirAll(dir, os.ModeDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create config directory: %v\n", err)
	}

	ctx := kong.Parse(&cli)
	err = ctx.Run(&Context{Verbose: cli.Verbose, ConfigDir: dir})
	ctx.FatalIfErrorf(err)
	/*entries, err := os.ReadDir(".")
	if err != nil {
		log.Fatalf("failed to read directory: %v", err)
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		func() {
			os.Chdir(e.Name())
			defer os.Chdir("..")

			fmt.Printf("\033[32m%s\033[0m", e.Name())

			out, err := exec.Command("git", "status").Output()
			if err != nil {
				log.Fatalf("failed to execure 'git status': %v", err)
			}
			if strings.Contains(string(out), "nothing to commit, working tree clean") {
				fmt.Println("-- Clean")
			} else {
				fmt.Printf("\n%s\n", string(out))
			}
		}()
		}*/
}

// ConfigDir returns the path of config directory for the application.
func ConfigDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return base + "/yagi", nil
}

func RepositoryInfoFile(configDir string) string {
	return configDir + "/repositories.json"
}
