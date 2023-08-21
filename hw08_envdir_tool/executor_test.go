package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("Test set env var", func(t *testing.T) {
		testFilename := "./testdata/test_env.txt"
		defer os.Remove(testFilename)

		cmd := []string{"/bin/bash", "-c", "echo -n \"$TESTME\" > " + testFilename}

		exptectedValue := "test"
		env := Environment{
			"TESTME": EnvValue{Value: exptectedValue, NeedRemove: false},
		}
		RunCmd(cmd, env)

		actualValue, _ := os.ReadFile(testFilename)

		require.Equal(t, exptectedValue, string(actualValue))
	})

	t.Run("Test unset env var", func(t *testing.T) {
		testFilename := "./testdata/test_env_rm.txt"
		defer os.Remove(testFilename)

		os.Setenv("UNSET", "123")
		cmd := []string{"/bin/bash", "-c", "echo -n \"$UNSET\" > " + testFilename}

		env := Environment{
			"UNSET": EnvValue{Value: "must not be set", NeedRemove: true},
		}
		RunCmd(cmd, env)

		actualValue, _ := os.ReadFile(testFilename)

		require.Equal(t, "", string(actualValue))
	})
}
