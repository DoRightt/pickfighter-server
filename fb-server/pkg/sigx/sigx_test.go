package sigx

import (
	"os"
	"syscall"
	"testing"
	"time"

	"gopkg.in/go-playground/assert.v1"
)

func TestListen(t *testing.T) {
	tests := []struct {
		Name           string
		SignalToSend   os.Signal
		ExpectedSignal os.Signal
	}{
		{"os.Interrupt signal", os.Interrupt, os.Interrupt},
		{"syscall.SIGTERM signal", syscall.SIGTERM, syscall.SIGTERM},
		{"syscall.SIGHUP signal", syscall.SIGHUP, syscall.SIGHUP},
		{"syscall.SIGQUIT signal", syscall.SIGQUIT, syscall.SIGQUIT},
	}

	process, err := os.FindProcess(os.Getpid())
	if err != nil {
		panic(err)
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			actualSignal := make(chan os.Signal, 1)
			defer close(actualSignal)

			signalFunc := func(s os.Signal) {
				actualSignal <- s
			}

			Listen(signalFunc)

			time.Sleep(10 * time.Millisecond)

			err = process.Signal(tc.SignalToSend)
			if err != nil {
				panic(err)
			}

			assert.Equal(t, tc.ExpectedSignal, <-actualSignal)
		})
	}
}
