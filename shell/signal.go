package shell 

import (
	"os"
	"os/signal"
	"syscall"
	"fmt"
)

var fgPid int = 0

func SetupSignalHandling(fgPid *int) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTSTP)

	go func() {
		for sig := range signalChan {
			if *fgPid == 0 {
				if sig == syscall.SIGINT || sig == syscall.SIGTSTP {
					fmt.Print("\n" + GetPrompt())
				}
				continue
			}

			// Forward signal to the foreground process group
			if err := syscall.Kill(-*fgPid, sig.(syscall.Signal)); err != nil {
				fmt.Printf("Error killing process: %v\n", err)
			}

		}
	}()
}
