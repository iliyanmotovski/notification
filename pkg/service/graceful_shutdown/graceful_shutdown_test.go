package gracefulshutdown

import (
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGracefulShutdown(t *testing.T) {
	callbackChan := make(chan bool)

	callback := func() {
		callbackChan <- true
	}

	GracefulShutdown(callback)

	termChan <- syscall.SIGINT

	assert.True(t, <-callbackChan)
}
