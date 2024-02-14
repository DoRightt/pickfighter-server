package version

import (
	"time"
)

const (
	DevVersion = "v0.1.0"
	Name       = "fb-server-api"
)

// RunDate contains the Unix timestamp representing the date and time of application start.
var RunDate = time.Now().Unix()

// TODO
var GitVersion string
var BuildDate string
var BuildCommit string
