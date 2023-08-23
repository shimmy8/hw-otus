package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("Test dir not found error", func(t *testing.T) {
		_, err := ReadDir("./123")
		require.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("Test env values correct", func(t *testing.T) {
		env, err := ReadDir("./testdata/env")

		require.NoError(t, err)

		expectedEnv := Environment{
			"BAR":   EnvValue{Value: "bar", NeedRemove: false},
			"EMPTY": EnvValue{Value: "", NeedRemove: true},
			"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
			"UNSET": EnvValue{Value: "", NeedRemove: true},
		}
		require.Equal(t, expectedEnv, env)
	})

	t.Run("Test env file = error", func(t *testing.T) {
		_, err := ReadDir("./testdata/env2")

		require.ErrorIs(t, err, ErrFileName)
	})

	t.Run("Test env val trim tab", func(t *testing.T) {
		env, err := ReadDir("./testdata/env3")

		require.NoError(t, err)

		expectedEnv := Environment{
			"TRIMMYTAB": EnvValue{Value: "HERE IT IS", NeedRemove: false},
		}
		require.Equal(t, expectedEnv, env)
	})
}
