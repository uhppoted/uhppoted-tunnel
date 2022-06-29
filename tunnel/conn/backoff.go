package conn

import (
	"fmt"
	"time"

	"github.com/uhppoted/uhppoted-tunnel/log"
)

const RETRY_MIN_DELAY = 5 * time.Second

type Backoff struct {
	retries       int
	retryDelay    time.Duration
	maxRetries    int
	maxRetryDelay time.Duration
}

func NewBackoff(maxRetries int, maxRetryDelay time.Duration) Backoff {
	return Backoff{
		retries:       0,
		retryDelay:    RETRY_MIN_DELAY,
		maxRetries:    maxRetries,
		maxRetryDelay: maxRetryDelay,
	}
}

func (b *Backoff) Reset() {
	b.retries = 0
	b.retryDelay = RETRY_MIN_DELAY
}

func (b *Backoff) Wait(tag string, closing chan struct{}) bool {
	b.retries++
	if b.maxRetries >= 0 && b.retries > b.maxRetries {
		fatalf(tag, "retry count exceeded %v", b.maxRetries)
		return false
	}

	infof(tag, "retrying in %v", b.retryDelay)

	select {
	case <-time.After(b.retryDelay):
		b.retryDelay *= 2
		if b.retryDelay > b.maxRetryDelay {
			b.retryDelay = b.maxRetryDelay
		}

	case <-closing:
		return false
	}

	return true
}

func debugf(tag, format string, args ...any) {
	f := fmt.Sprintf("%-6v %v", tag, format)

	log.Debugf(f, args...)
}

func infof(tag string, format string, args ...any) {
	f := fmt.Sprintf("%-6v %v", tag, format)

	log.Infof(f, args...)
}

func warnf(tag, format string, args ...any) {
	f := fmt.Sprintf("%-6v %v", tag, format)

	log.Warnf(f, args...)
}

func errorf(tag string, format string, args ...any) {
	f := fmt.Sprintf("%-6v %v", tag, format)

	log.Errorf(f, args...)
}

func fatalf(tag string, format string, args ...any) {
	f := fmt.Sprintf("%-6v %v", tag, format)

	log.Fatalf(f, args...)
}
