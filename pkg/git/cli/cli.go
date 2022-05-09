package cli

import (
	"fmt"
	"github.com/champ-oss/file-sync/pkg/common"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

func Clone(sourceRepo string) (dir string, err error) {
	log.Debug("Creating temp directory for source repository")
	dir, _ = ioutil.TempDir("", "source")

	log.Infof("Cloning source repository %s to %s", sourceRepo, dir)
	output, err := common.RunCommand("./", "git", "clone", sourceRepo, dir)
	if err != nil {
		return dir, fmt.Errorf(output)
	}
	return dir, nil
}

func Fetch(repoDir string) error {
	output, err := common.RunCommand(repoDir, "git", "fetch")
	if err != nil {
		return fmt.Errorf(output)
	}
	return nil
}

func Branch(repoDir, branchName string) error {
	output, err := common.RunCommand(repoDir, "git", "branch", branchName)
	if err != nil {
		return fmt.Errorf(output)
	}
	return nil
}

func Checkout(repoDir, branchName string) error {
	output, err := common.RunCommand(repoDir, "git", "checkout", branchName)
	if err != nil {
		return fmt.Errorf(output)
	}
	return nil
}

func Status(repoDir, fileName string) string {
	output, err := common.RunCommand(repoDir, "git", "status", "--porcelain", fileName)
	if err != nil {
		return err.Error()
	}
	return output
}

func Add(repoDir, fileName string) error {
	output, err := common.RunCommand(repoDir, "git", "add", fileName)
	if err != nil {
		return fmt.Errorf(output)
	}
	return nil
}

func Commit(repoDir, message string) error {
	output, err := common.RunCommand(repoDir, "git", "commit", "-m", message)
	if err != nil {
		return fmt.Errorf(output)
	}
	return nil
}

func Push(repoDir, branchName string) error {
	output, err := common.RunCommand(repoDir, "git", "push", "--set-upstream", "origin", branchName)
	if err != nil {
		return fmt.Errorf(output)
	}
	return nil
}
