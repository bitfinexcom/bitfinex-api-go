package utils

import (
	"fmt"
	"os"
	"runtime/debug"
	"syscall"
)

var PanicMessage string

func GoWithRecover(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				PanicMessage = fmt.Sprintf("panic: %v\n%v\n", err, string(debug.Stack()))
				_, _ = fmt.Fprint(os.Stderr, PanicMessage)

				if err := os.Stderr.Sync(); err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "sync stdout err: %v\n", err.Error())
				}

				if err := syscall.Kill(os.Getpid(), syscall.SIGUSR1); err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "send signal err: %v", err)
				}
			}
		}()

		f()
	}()
}
