package httplib

// CookieName is the name of the cookie used to store the session ID.
var CookieName = "session"

// SetCookieName sets the name of the cookie used to store the session ID.
func SetCookieName(name string) {
	CookieName = name
}
