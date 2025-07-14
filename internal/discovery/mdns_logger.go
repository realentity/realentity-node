package discovery

import (
	"log"
	"os"
	"strings"
)

// QuietMDNSLogger wraps the standard logger to suppress mDNS multicast warnings
type QuietMDNSLogger struct {
	original *log.Logger
	quiet    bool
}

// NewQuietMDNSLogger creates a logger that can suppress mDNS warnings
func NewQuietMDNSLogger(quiet bool) *QuietMDNSLogger {
	return &QuietMDNSLogger{
		original: log.New(os.Stderr, "", log.LstdFlags),
		quiet:    quiet,
	}
}

// Printf implements a custom Printf that can filter mDNS warnings
func (q *QuietMDNSLogger) Printf(format string, v ...interface{}) {
	if q.quiet {
		// Filter out common mDNS multicast interface warnings
		if strings.Contains(format, "Failed to set multicast interface") ||
			strings.Contains(format, "no such interface") ||
			strings.Contains(format, "mdns:") {
			// Suppress this warning
			return
		}
	}

	// Pass through to original logger
	q.original.Printf(format, v...)
}

// SetOutput sets the output destination for the logger
func (q *QuietMDNSLogger) SetOutput(w os.File) {
	q.original.SetOutput(&w)
}
