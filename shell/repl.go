package shell

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"gosh/util"
	// "gosh/shell/execution"
)

// type Shell struct {
// 	CurrentPwd string				// Tracks current working directory
// 	OldPwd		string				// Tracks previous working directory
// 	Env			map[string]string	// Stores shell environment variables
// 	Jobs		map[int]*Job		// Tracks fg/bg jobs
// 	LastStatus	int 				// Stores the exit code of the last command
// }

func (sh *Shell) Start() {
	
	sh.SetupSignalHandling()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(GetPrompt()) // Display prompt
		// Read the command from the user
		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("exit") // Print "exit" like bash/zsh does
				os.Exit(0)
			}
			fmt.Println("Error reading input:", err)
			break
		}

		args := strings.Fields(input) // Tokenize input
		if len(args) == 0 {
			continue
		}

		for i, arg := range args {
			args[i] = util.ExpandHomeDirectory(arg)
		}

		// Check if the command is built-in or external
		if err := sh.HandleBuiltInCommand(args); err == nil {
			continue
		}

		// If not a built-in command, handle execution
		sh.Execute(args)
	}
}