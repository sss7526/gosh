package shell

import (
	"os"
	"strings"
	"fmt"
)

func GetPrompt() string {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "?"
	}

	homeDir := os.Getenv("HOME")
	if strings.HasPrefix(cwd, homeDir) {
		// if cwd == homeDir {
		// 	cwd = "~"
		// } else {
		// 	cwd = strings.Replace(cwd, homeDir, "~", 1)
		// }
		cwd = "~" + strings.TrimPrefix(cwd, homeDir)
	}

	user := os.Getenv("USER")
	if user == "" {
		user = "?"
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "?"
	}

	promptEnd := "$"
	if os.Geteuid() == 0 {
		promptEnd = "#"
	}
	return fmt.Sprintf("%s@%s:%s%s ", user, hostname, cwd, promptEnd)
}