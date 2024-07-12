package sigx

import (
	"os"
	"os/signal"
	"syscall"
)

// Listen starts a goroutine to listen for OS signals such as os.Interrupt,
// syscall.SIGTERM, syscall.SIGHUP, and syscall.SIGQUIT. When any of these signals
// is received, it invokes the provided function fn with the received signal as an argument.
// If fn is nil, no action is taken when a signal is received.
func Listen(fn func(os.Signal)) {
	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

		sig := <-sigchan
		if fn != nil {
			fn(sig)
		}
	}()
}
