package cmd

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	cmd := rootCmd
	cmd.SetArgs([]string{"--version"})

	err := cmd.Execute()
	assert.NoError(t, err)
}

func TestInitZapLogger(t *testing.T) {
	initZapLogger()
	assert.NotNil(t, logger, "Logger should be initialized")
}

func TestBindViperFlag(t *testing.T) {
	cmd := &cobra.Command{}

	cmd.Flags().String("testFlag", "", "Test flag")

	bindViperFlag(cmd, "testFlag", "testFlag")

	cmd.SetArgs([]string{"--testFlag=testViperVal"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "testViperVal", viper.GetString("testFlag"), "Expected Viper flag to be bound")
}

func TestBindViperNilFlag(t *testing.T) {
	var buf bytes.Buffer

	log.SetOutput(&buf)

	cmd := &cobra.Command{}

	bindViperFlag(cmd, "testFlag", "")

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Failed to bind viper flag:", "Expected log message about failed to bind viper flag")

	log.SetOutput(os.Stderr)
}

func TestBindViperPersistentFlag(t *testing.T) {
	cmd := &cobra.Command{}

	cmd.PersistentFlags().String("testPersistentFlag", "", "Test persistent flag")

	bindViperPersistentFlag(cmd, "testPersistentFlag", "testPersistentFlag")

	cmd.SetArgs([]string{"--testPersistentFlag=testPersistentViperVal"})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "testPersistentViperVal", viper.GetString("testPersistentFlag"), "Expected Viper persistent flag to be bound")
}

func TestBindViperPersistentNilFlag(t *testing.T) {
	var buf bytes.Buffer

	log.SetOutput(&buf)

	cmd := &cobra.Command{}

	bindViperPersistentFlag(cmd, "testPersistentFlag", "")

	logOutput := buf.String()
	assert.Contains(t, logOutput, "Failed to bind viper flag:", "Expected log message about failed to bind viper flag")

	log.SetOutput(os.Stderr)
}
