package logger

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInitialize(t *testing.T) {
	logLevel := "info"
	wrongLogLevel := "wrong"
	logFilePath := "test.log"

	err := Initialize(logLevel, logFilePath)

	assert.Nil(t, err, "Initialization should not return an error")

	assert.NotNil(t, logger, "Logger should be initialized")

	err = Initialize(wrongLogLevel, logFilePath)

	assert.NotNil(t, err, "Initialize should return error")
}

func TestGet(t *testing.T) {
	logger := Get()

	assert.NotNil(t, logger, "Logger should be initialized")
}

func TestNew(t *testing.T) {
	logger := New()

	assert.NotNil(t, logger, "Logger shouldn't be nil")
	viper.Reset()
}

func TestNewSugared(t *testing.T) {
	sugaredLogger := NewSugared()

	assert.NotNil(t, sugaredLogger, "SugaredLogger shouldn't be nil")
	viper.Reset()
}

func TestScrapperLogger(t *testing.T) {
	filePath := "logger/logs/test-scraper-log.json"
	logDirectory := "logger/logs/"

	if _, err := os.Stat(logDirectory); os.IsNotExist(err) {
		err = os.MkdirAll(logDirectory, 0755)
		if err != nil {
			t.Fatalf("Failed to create log directory: %v", err)
		}
	}

	defer os.RemoveAll(logDirectory)

	scrapperLogger, err := ScraperLogger()
	defer os.Remove(filePath)

	assert.NotNil(t, scrapperLogger, "ScrapperLogger should not be nil")
	assert.Nil(t, err, "Error should be nil")
}
