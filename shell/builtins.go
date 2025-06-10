package shell

import (
	"errors"
)

// Predefined error for commands that are not built-in
var ErrNotBuiltInCommand = errors.New("not a built-in command")

// Map to store built-in command with their corresponding handlers
var builtInCommands = map[string]func(*Shell, []string) error{
	"exit": (*Shell).Exit,
	"cd": (*Shell).Cd,
	"pwd": (*Shell).Pwd,
}


// HandleBuiltInCommand tries to match and execute a built-in command.
// If the command is not built-in, it returns an error.
func (sh *Shell) HandleBuiltInCommand(args []string) error {
	if handler, exists := builtInCommands[args[0]]; exists {
		return handler(sh, args[1:])
	}
	return ErrNotBuiltInCommand
}