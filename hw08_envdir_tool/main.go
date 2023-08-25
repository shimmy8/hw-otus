package main

import (
	"fmt"
	"os"
)

func main() {
	envDir := os.Args[1]
	cmd := os.Args[2:]

	env, err := ReadDir(envDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	cmdExitCode := RunCmd(cmd, env)
	os.Exit(cmdExitCode)
}
