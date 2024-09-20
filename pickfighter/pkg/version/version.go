package version

import "time"

const (
	DevVersion = "v0.3.2"
	Name       = "Pickfighter-gateway-service"
)

// RunDate contains the Unix timestamp representing the date and time of service start.
var RunDate = time.Now().Unix()
