package version

import "time"

const (
	DevVersion = "v0.4.0"
	Name       = "pickfighter-gateway-service"
)

// RunDate contains the Unix timestamp representing the date and time of service start.
var RunDate = time.Now().Unix()
