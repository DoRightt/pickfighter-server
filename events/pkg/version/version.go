package version

import "time"

const (
	DevVersion = "v0.1.0"
	Name       = "event-service"
)

// RunDate contains the Unix timestamp representing the date and time of service start.
var RunDate = time.Now().Unix()
