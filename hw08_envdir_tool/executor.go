package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec

	command.Env = buildCmdEnv(env)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err := command.Run()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}
		log.Fatal(err)
	}
	return 0
}

func buildCmdEnv(env Environment) []string {
	// build map of current os.Environ variables
	cmdEnvMap := make(map[string]string)
	for _, envStr := range os.Environ() {
		envSlice := strings.Split(envStr, "=")
		cmdEnvMap[envSlice[0]] = envSlice[1]
	}
	// change map to match with input env values
	for envName, envVal := range env {
		_, inEnv := cmdEnvMap[envName]
		if envVal.NeedRemove && inEnv {
			delete(cmdEnvMap, envName)
		} else {
			cmdEnvMap[envName] = envVal.Value
		}
	}
	// reform env map to a slice of key=value strings
	cmdEnvSlice := []string{}
	for envName, envVal := range cmdEnvMap {
		cmdEnvSlice = append(cmdEnvSlice, strings.Join([]string{envName, envVal}, "="))
	}
	return cmdEnvSlice
}
