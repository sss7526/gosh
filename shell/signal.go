package shell 

import (
	"os"
	"os/signal"
	"syscall"
	"fmt"
)

// var fgPid int = 0

func (sh *Shell) SetupSignalHandling() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTSTP)

	go func() {
		for sig := range signalChan {
			sh.lock.Lock()
			fmt.Printf("FGPID: %d\n", sh.fgPid)
			if sh.fgPid == 0 {
				// No fg process, just print a fresh prompt
				if sig == syscall.SIGINT || sig == syscall.SIGTSTP {
					fmt.Print("\n" + GetPrompt())
				}
				sh.lock.Unlock()
				continue
			}

			// Forward signal to the foreground process group
			if err := syscall.Kill(sh.fgPid, sig.(syscall.Signal)); err != nil {
				fmt.Printf("Error killing process: %v\n", err)
			}
			sh.lock.Unlock()
		}
	}()
}
