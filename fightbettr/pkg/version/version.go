package version

import "time"

const (
	DevVersion = "v0.1.0"
	Name       = "Fightbettr-gateway-service"
)

// RunDate contains the Unix timestamp representing the date and time of service start.
var RunDate = time.Now().Unix()
