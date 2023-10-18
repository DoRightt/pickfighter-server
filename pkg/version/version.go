package version

import (
	"time"
)

const (
	DevVersion = "v0.1.0"
	Name       = "fb-server-api"
)

var RunDate = time.Now().Unix()
var GitVersion string
var BuildDate string
var BuildCommit string
