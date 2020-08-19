package internal

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/magefile/mage/sh"
	"github.com/wavesoftware/go-ensure"
	"github.com/wavesoftware/go-magetasks/config"
	"github.com/wavesoftware/go-magetasks/internal/tasks"
)

var (
	git                 = sh.OutCmd("git")
	gitVerCache *string = nil
)

// BuildDir returns project build dir
func BuildDir() string {
	return relativeToRepo(config.BuildDirPath)
}

func RepoDir() string {
	if config.RepoDir != "" {
		return config.RepoDir
	}
	repoDir, err := os.Getwd()
	ensure.NoError(err)
	return repoDir
}

func AppendLdflags(args []string, t *tasks.Task) []string {
	if isVersionVariableSet() {
		p, err := versionVariablePath()
		if err != nil {
			t.End(err)
			return args
		}
		ldflags := fmt.Sprintf("-ldflags=-X %s=%s", p, gitVersion())
		args = append(args, ldflags)
	}
	return args
}

func isVersionVariableSet() bool {
	return config.VersionVariablePath != ""
}

func versionVariablePath() (string, error) {
	if !isVersionVariableSet() {
		return "", errors.New("version variable path is unset")
	}
	return config.VersionVariablePath, nil
}

func relativeToRepo(paths []string) string {
	fullpath := make([]string, len(paths)+1)
	fullpath[0] = RepoDir()
	for ix, elem := range paths {
		fullpath[ix+1] = elem
	}
	return path.Join(fullpath...)
}
