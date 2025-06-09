package util

import (
	"path/filepath"
	"strings"
	"os"
	"fmt"
)

func ExpandHomeDirectory(path string) string {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if home == "" || err != nil {
			return path // If HOME is not set, return path unchanged
		}

		return strings.Replace(path, "~", home, 1)
	}
	return path
}

func isAbsolutePath(path string) bool {
	return filepath.IsAbs(path)
}

func isExplicitRelativepath(path string) bool {
	return strings.HasPrefix(path, "./") || strings.HasPrefix(path, "../")
}

func isRegularExecutableFile(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}

	// Check that it's a regular file and has at least one executable (`x`) permission bit.
	return stat.Mode().IsRegular() && (stat.Mode().Perm()&0111) != 0
}

func ResolveExecutable(command string) (string, error) {

	if isAbsolutePath(command) || isExplicitRelativepath(command) {
		if isRegularExecutableFile(command) {
			return command, nil
		}
		return "", fmt.Errorf("%s: not an executable file", command)
	}

	pathEnv := os.Getenv("PATH")
	for dir := range strings.SplitSeq(pathEnv, string(os.PathListSeparator)) {
		fullPath := filepath.Join(dir, command)
		if isRegularExecutableFile(fullPath) {
			return fullPath, nil
		}
	}

	return "", fmt.Errorf("%s: command not found", command)
}