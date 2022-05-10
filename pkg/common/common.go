package common

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func RemoveDir(dir string) {
	if err := os.RemoveAll(dir); err != nil {
		log.Error(err)
	}
}

func CopySourceFiles(files []string, sourceDir, destDir string) error {
	for _, f := range files {
		sourcePath := filepath.Join(sourceDir, f)
		destPath := filepath.Join(destDir, f)
		log.Debugf("Copying %s to %s", sourcePath, destPath)
		if err := CopyFile(sourcePath, destPath); err != nil {
			log.Error("error copying files from source")
			return err
		}
	}
	return nil
}

func CopyFile(source, dest string) error {
	input, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}

	if baseDir, _ := filepath.Split(dest); baseDir != "" {
		if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
			return err
		}
	}

	err = ioutil.WriteFile(dest, input, 0644)
	return err
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

func RunCommandNoLog(dir, cmd string, args ...string) error {
	command := exec.Command(cmd, args...)
	command.Dir = dir

	err := command.Run()

	if err != nil {
		return fmt.Errorf("error running command")
	}
	return nil
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

func GetOwner() string {
	owner := os.Getenv("GITHUB_REPOSITORY_OWNER")
	if owner == "" {
		log.Fatal("GITHUB_REPOSITORY_OWNER env is empty")
	}
	return owner
}

func GetRepo() string {
	ownerRepo := os.Getenv("GITHUB_REPOSITORY")
	if ownerRepo == "" {
		log.Fatal("GITHUB_REPOSITORY env is empty")
	}
	parts := strings.Split(ownerRepo, "/")
	if len(parts) != 2 {
		log.Fatalf("GITHUB_REPOSITORY is in unexpected format: %s", ownerRepo)
	}
	return parts[1]
}

func GetFiles() []string {
	files := os.Getenv("INPUT_FILES")
	if files == "" {
		log.Fatal("INPUT_FILES env is empty")
	}
	return strings.Split(files, "\n")
}
