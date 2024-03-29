package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var ErrFileName = errors.New("env file name error")

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	env := make(Environment)
	for _, file := range files {
		value, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return env, err
		}
		if strings.Contains(file.Name(), "=") {
			return nil, ErrFileName
		}
		cleanValue := cleanValue(value)
		env[file.Name()] = EnvValue{Value: cleanValue, NeedRemove: len(cleanValue) == 0}
	}
	return env, nil
}

func cleanValue(inp []byte) string {
	value := strings.Split(string(inp), "\n")[0]
	value = strings.ReplaceAll(value, "\x00", "\n")
	return strings.TrimRight(value, " \t")
}
