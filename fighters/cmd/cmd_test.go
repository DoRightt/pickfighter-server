package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"

	"fightbettr.com/pkg/model"
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

func TestValidateServerArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected error
	}{
		{"ValidRoute", []string{model.FightersService}, nil},
		{"EmptyArgs", []string{}, errEmptyApiRoute},
		{"InvalidRoute", []string{"invalidRoute"}, fmt.Errorf("allowed routes are: %s", model.FightersService)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			err := validateServerArgs(cmd, tc.args)

			assert.Equal(t, tc.expected, err)
		})
	}
}

func TestReadFighterData(t *testing.T) {
	fighters, err := ReadFighterData()

	assert.NoError(t, err)
	assert.True(t, len(fighters) > 0)
}

func TestWriteFighterData(t *testing.T) {

}

func TestDeleteFighterData(t *testing.T) {

}

func TestCreateFighter(t *testing.T) {

}

func TestUpdateFighter(t *testing.T) {

}
