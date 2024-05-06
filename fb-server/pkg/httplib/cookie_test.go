package httplib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetCookieName(t *testing.T) {
	initialCookieName := CookieName
	newCookieName := "test-name"

	SetCookieName(newCookieName)

	assert.Equal(t, newCookieName, CookieName, "CookieName should be updated")

	CookieName = initialCookieName
}
