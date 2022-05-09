package common

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

func RemoveDir(dir string) {
	if err := os.RemoveAll(dir); err != nil {
		log.Error(err)
	}
}

func RunCommand(dir, cmd string, args ...string) (output string, err error) {
	LogCommand(cmd, args...)
	command := exec.Command(cmd, args...)

	var stdout bytes.Buffer
	command.Stdout = &stdout
	var stderr bytes.Buffer
	command.Stderr = &stderr
	command.Dir = dir

	err = command.Run()
	LogOutput(stdout)
	LogOutput(stderr)

	if err != nil {
		return stderr.String(), err
	}
	return stdout.String(), nil
}

func LogCommand(cmd string, args ...string) {
	logMessage := cmd
	for _, a := range args {
		logMessage += " " + a
	}
	log.Info(logMessage)
}

func LogOutput(output bytes.Buffer) {
	if output.String() == "" {
		return
	}
	fmt.Print(output.String())
}
