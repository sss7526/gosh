package shell

import (
	"errors"
	// "gosh/cmd"
)

// Map to store built-in command with their corresponding handlers
var builtInCommands = map[string]func(*Shell, []string) error{
	"exit": (*Shell).Exit,
	"cd": (*Shell).Cd,
	"pwd": (*Shell).Pwd,
}

func (sh *Shell) HandleBuiltInCommand(args []string) error {
	if handler, exists := builtInCommands[args[0]]; exists {
		return handler(sh, args[1:])
	}
	return errors.New("not a built-in command")
}