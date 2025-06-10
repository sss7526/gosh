package shell

import (
	"os"
	"strings"
)

// NewShell creates a new instance of the Shell with default state.
func NewShell() *Shell {
	return &Shell{
		CurrentPwd: getPwd(),			// Start with the current working directory
		OldPwd: "",
		Env:		buildEnvMap(),		// Build environment map from the OS
		Jobs:		make(map[int]*Job),	// No jobs initially
		fgPid:		0,					// No current foreground process
		LastStatus: 0,					// Neutral exit status
	}
}

// Clone creates a new Shell instance that shares the same environment
// but starts with a clean job list and state.
func (sh *Shell) Clone() *Shell {
	return &Shell{
		CurrentPwd:	sh.CurrentPwd,
		OldPwd:		sh.OldPwd,
		Env:		copyEnvMap(sh.Env),	// Copy the environment variables
		Jobs:		make(map[int]*Job),	// Subshell starts with no active jobs
		LastStatus:	sh.LastStatus,		// Copies the last status
	}
}

// Helper function to deeply clone the environment map.
func copyEnvMap(env map[string]string) map[string]string {
	newEnv := make(map[string]string)
	for key, value := range env {
		newEnv[key] = value
	}
	return newEnv
}

// Helper function to get the current working directory.
func getPwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "/" // Fallback to root directory
	}
	return cwd
}

// Helper function to build an environment variable map from os.Environ().
func buildEnvMap() map[string]string {
	envMap := make(map[string]string)
	for _, kv := range os.Environ() {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}
	return envMap
}