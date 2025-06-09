package cmd 

import (
	"fmt"
	"os"
	"strings"
	"path/filepath"
)

func Pwd(args []string) error {
	var physical = false
	validOption := false

	for _, arg := range args {
		if len(args) > 0 && strings.HasPrefix(arg, "-") {
			lastChar := rune(arg[len(arg)-1])
			switch lastChar {
			case 'L':
				physical = false
				validOption = true
			case 'P':
				physical = true
				validOption = true
			default:
				fmt.Printf("pwd: invalid option -- '%c'\n", lastChar)
				return fmt.Errorf("pwd: invalid option -- '%c'", lastChar)
			}
		}
	}

	// Default to logical (`-L`) behavior if no valid option was provided
	if !validOption {
		physical = false
	}

	// Print the current directory based on the selected mode
	if physical {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Printf("pwd: error getting current working directory: %v\n", err)
		}
		realPath, err := filepath.EvalSymlinks(wd)
		if err != nil {
			fmt.Printf("pwd: error getting physical path: %v\n", err)
			return fmt.Errorf("pwd: error getting physical path: %v", err)
		}
		fmt.Println(realPath)
	} else {
		fmt.Println(currentPwd)
	}
	return nil
}