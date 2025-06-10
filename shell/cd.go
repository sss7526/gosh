package shell

import (
	"os"
	"path/filepath"
	"fmt"
	"strings"
)

func (sh *Shell) Cd(args []string) error {
	var physical = false
	var resolveErrorOnPhysical = false

	// Handle the `cd -` case FIRST
	if len(args) > 0 && args[0] == "-" {
		if sh.OldPwd == "" {
			return fmt.Errorf("cd: OLDPWD not set")
		}
		targetDir := sh.OldPwd
		absTarget, err := filepath.Abs(targetDir)
		if err != nil {
			return fmt.Errorf("cd: error resolving absolute path of OLDPWD: %v", err)
		}
		fmt.Println(targetDir) // Print the absolute path for the new directory for `cd -` which is what bash does
		if err := os.Chdir(absTarget); err != nil {
			return fmt.Errorf("cd: %v", err)
		}
		sh.OldPwd = sh.CurrentPwd
		sh.CurrentPwd = absTarget
		os.Setenv("OLDPWD", sh.OldPwd)
		os.Setenv("PWD", sh.CurrentPwd)
		return nil
	}

	// Parse options
	i := 0
	for i < len(args) {
		if !strings.HasPrefix(args[i], "-") { // Stop option parsing at the first non-flag
			break
		}
		for _, c := range args[i][1:] {
			switch c {
			case 'L':
				physical = false // Logical path resolution (default)
			case 'P':
				physical = true // Resolve real paths, avoiding symlinks
			case 'e':
				resolveErrorOnPhysical = true // Fail on physical path issues
			case '@':
				// Extended attributes placeholder (not implemented)
			default:
				return fmt.Errorf("cd: invalid option -- '%c'", c)
			}
		}
		i++
	}

	var targetDir string
	if i >= len(args) {
		// No directory provided: use `$HOME`
		targetDir = os.Getenv("HOME")
		if targetDir == "" {
			return fmt.Errorf("cd: HOME not set")
		}
	} else {
		targetDir = args[i]
	}

	// Apply CDPATH logic for relative paths
	if !filepath.IsAbs(targetDir) {
		cdpath := os.Getenv("CDPATH")
		if cdpath != "" {
			for dir := range strings.SplitSeq(cdpath, string(os.PathListSeparator)) {
				candidate := filepath.Join(dir, targetDir)
				if _, err := os.Stat(candidate); err == nil {
					targetDir = candidate
					break
				}
			}
		}
	}


	// Resolve the absolute path of targetDir
	absTarget, err := filepath.Abs(targetDir)
	if err != nil {
		return fmt.Errorf("cd: error resolving absolute path: %v", err)
	}
	targetDir = absTarget
	
	// Handle physical path resolution (`-P`) before changing directories
	if physical {
		absTarget, err := filepath.EvalSymlinks(targetDir)
		if err != nil {
			fmt.Printf("cd: %v\n", err)
			if resolveErrorOnPhysical {
				return fmt.Errorf("cd: error resolving physical path") // Fail only for `-Pe`
			}
			return fmt.Errorf("cd: %v", err)
		}
		targetDir = absTarget // Physical path for `-P`
	}

	// Change the directory to the resolved `targetDir`
	if err := os.Chdir(targetDir); err != nil {
		return fmt.Errorf("cd: %v", err)
	}

	// Update `$PWD` and `$OLDPWD`
	sh.OldPwd = sh.CurrentPwd
	sh.CurrentPwd = targetDir // Use logical (`-L`) or physical (`-P`) path as appropriate
	os.Setenv("OLDPWD", sh.OldPwd)
	os.Setenv("PWD", sh.CurrentPwd)

	return nil
}